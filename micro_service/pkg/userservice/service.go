package userservice

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"userService/pkg/cache"
	"userService/pkg/common"
	insmodel "userService/pkg/model/institution"
	mchtmodel "userService/pkg/model/merchant"
	usermodel "userService/pkg/model/user"
	"userService/pkg/pb"
	"userService/pkg/rbac"
	"userService/pkg/util"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/metadata"
)

type userService struct{}

func (u *userService) RemoveRoute(ctx context.Context, in *pb.RemoveRouteRequest) (*pb.RemoveRouteReply, error) {
	reply := &pb.RemoveRouteReply{}
	if in.Route == "" {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "路由不能为空",
		}
		return reply, nil
	}
	db := common.DB

	err := usermodel.DeleteRoute(db, &usermodel.Route{Name: in.Route})
	if err != nil {
		return nil, err
	}
	return reply, err
}

type UserInfo struct {
	ID       int64
	UserName string
}

func (u *userService) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginReply, error) {
	db := common.DB.New()
	if in.GetUsername() == "" || in.GetPassword() == "" {
		return &pb.LoginReply{
			Err: &pb.Error{
				Code:        http.StatusBadRequest,
				Message:     "InvalidParamsError",
				Description: "用户或密码为空",
			},
		}, nil
	}

	// 查询用户
	user, err := usermodel.FindUserByUserName(db, in.GetUsername())
	if err != nil {
		return nil, err
	}
	if user == nil {
		return &pb.LoginReply{
			Err: &pb.Error{
				Code:        http.StatusNotFound,
				Message:     "NotFoundError",
				Description: "用户不存在",
			},
		}, nil
	}

	hash := user.PasswordHashNew
	if hash == "" {
		hash = user.PasswordHash
	}
	// 校验密码
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(in.GetPassword()))
	if err != nil {
		return &pb.LoginReply{
			Err: &pb.Error{
				Code:        http.StatusBadRequest,
				Message:     InvalidParam,
				Description: "密码错误",
			},
		}, nil
	}

	// 生成token
	expiredAt := time.Now().Add(time.Hour * 72)
	userIdStr := fmt.Sprintf("%d", user.UserID)
	tk, err := genToken(userIdStr, expiredAt)
	if err != nil {
		return nil, err
	}

	userMap := map[string]interface{}{
		"id":        user.UserID,
		"username":  user.UserName,
		"email":     user.Email,
		"leaguerNo": user.LeaguerNO,
		"token":     tk,
	}

	bs, _ := json.Marshal(userMap)
	client := common.RedisClient
	err = cache.SetUserInfo(client, userIdStr, string(bs), time.Duration(expiredAt.UnixNano()))
	if err != nil {
		return &pb.LoginReply{
			Err: &pb.Error{
				Code:        http.StatusInternalServerError,
				Message:     INTERNAL,
				Description: err.Error(),
			},
		}, nil
	}

	//添加新密码
	if user.PasswordHashNew == "" {
		logrus.Infoln("更新密码")
		newHash, _ := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.MinCost)
		logrus.Infoln("newHash", newHash)
		err = usermodel.UpdateUser(db, user.UserID, &usermodel.User{
			PasswordHashNew: string(newHash),
		})
		if err != nil {
			return nil, err
		}
	}
	return &pb.LoginReply{
		Id:         user.UserID,
		Username:   user.UserName,
		UserType:   user.UserType,
		LeaguerNo:  user.LeaguerNO,
		Email:      user.Email,
		UserStatus: user.UserStatus,
		CreatedAt:  user.CreatedAt.Unix(),
		Token:      tk,
	}, nil
}

func (u *userService) GetPermissions(ctx context.Context, in *pb.GetPermissionsRequest) (*pb.GetPermissionsReply, error) {
	reply := &pb.GetPermissionsReply{}
	db := common.DB
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "找不到用户信息",
		}
		return reply, nil
	}
	userids := md.Get("userid")
	if len(userids) == 0 {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "找不到用户信息",
		}
		return reply, nil
	}

	permissions := common.Enforcer.GetImplicitPermissionsForUser(fmt.Sprintf("user:%s", userids[0]))

	itemNames := make([]string, 0)
	wild := make([]string, 0)
	for _, permission := range permissions {
		if len(permission) > 1 {
			ps := permission[1]
			itemNames = append(itemNames, ps)
			if strings.Contains(ps, "*") {
				wild = append(wild, strings.ReplaceAll(ps, "*", "%"))
			}
		}
	}

	menus, err := usermodel.GetAuthMenu(db, itemNames, false)
	if err != nil {
		return nil, err
	}

	wildMenus, err := usermodel.GetAuthMenu(db, wild, true)
	if err != nil {
		return nil, err
	}
	menus = append(menus, wildMenus...)

	idMap := make(map[int32]bool)
	replyMenus := make([]*pb.Menu, 0, len(menus))

	for _, m := range menus {
		if !idMap[m.ID] {
			replyMenus = append(replyMenus, &pb.Menu{
				Id:     m.ID,
				Name:   m.Name,
				Parent: m.Parent,
				Route:  m.MenuRoute,
				Data:   m.MenuData,
				Order:  m.MenuOrder,
			})
			idMap[m.ID] = true
		}
	}
	logrus.Debugln("------------", len(replyMenus))
	return &pb.GetPermissionsReply{
		Menus: replyMenus,
	}, nil
}

func (u *userService) CheckPermission(ctx context.Context, in *pb.CheckPermissionRequest) (*pb.CheckPermissionReply, error) {
	reply := &pb.CheckPermissionReply{}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "找不到用户信息",
		}
		return reply, nil
	}
	ids := md.Get("userid")
	if !ok || len(ids) == 0 {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "找不到用户信息",
		}
		return reply, nil
	}

	ok = common.Enforcer.Enforce(fmt.Sprintf("user:%s", ids[0]), in.GetRoute(), "read")
	return &pb.CheckPermissionReply{
		Result: ok,
	}, nil
}

func (u *userService) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterReply, error) {
	reply := &pb.RegisterReply{}
	db := common.DB
	if in.Username == "" || in.Password == "" || in.UserType == "" || in.Email == "" || in.LeaguerNo == "" {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "用户信息不全",
		}
		return reply, nil
	}
	switch in.UserType {
	case "admin":
		break
	case "merchant":
		mcht, err := mchtmodel.FindMerchantInfoById(db, in.UserGroupNo)
		if err != nil {
			return nil, err
		}
		if mcht == nil {
			reply.Err = &pb.Error{
				Code:        http.StatusBadRequest,
				Message:     InvalidParam,
				Description: "商户号不存在",
			}
			return reply, nil
		}
	case "institution":
		ins, err := insmodel.FindInstitutionInfoById(db, in.UserGroupNo)
		if err != nil {
			return nil, err
		}
		if ins == nil {
			reply.Err = &pb.Error{
				Code:        http.StatusBadRequest,
				Message:     InvalidParam,
				Description: "机构号不存在",
			}
			return reply, nil
		}
	case "institution_group":
		id, err := strconv.ParseInt(in.UserGroupNo, 10, 64)
		if err != nil {
			reply.Err = &pb.Error{
				Code:        http.StatusBadRequest,
				Message:     InvalidParam,
				Description: "机构组号错误",
			}
			return reply, nil
		}
		ins, err := insmodel.FindInsGroupById(db, id)
		if err != nil {
			return nil, err
		}
		if ins == nil {
			reply.Err = &pb.Error{
				Code:        http.StatusBadRequest,
				Message:     InvalidParam,
				Description: "机构组号不存在",
			}
			return reply, nil
		}
	default:
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "用户类型错误",
		}
		return reply, nil
	}

	user, err := usermodel.FindUserByUserName(db, in.Username)
	if err != nil {
		return nil, err
	}
	if user != nil {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     AlreadyExists,
			Description: "用户已存在",
		}
		return reply, nil
	}

	bs, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.MinCost)
	if err != nil {
		return nil, err
	}

	user = &usermodel.User{
		UserName:        in.Username,
		UserType:        in.UserType,
		Email:           in.Email,
		LeaguerNO:       in.LeaguerNo,
		PasswordHash:    string(bs),
		PasswordHashNew: string(bs),
		UserStatus:      1,
		UserGroupNo:     in.UserGroupNo,
	}

	newUser, err := usermodel.SaveUser(db, user)
	if err != nil {
		return nil, err
	}

	return &pb.RegisterReply{
		Id:         newUser.UserID,
		LeaguerNo:  newUser.LeaguerNO,
		Username:   newUser.UserName,
		Email:      newUser.Email,
		UserType:   newUser.UserType,
		UserStatus: newUser.UserStatus,
		CreatedAt:  newUser.CreatedAt.UnixNano() / int64(time.Millisecond),
	}, nil
}

func (u *userService) AddPermissionForRole(ctx context.Context, in *pb.AddPermissionForRoleRequest) (*pb.AddPermissionForRoleReply, error) {
	reply := &pb.AddPermissionForRoleReply{}
	if in.Role == "" || in.Permission == "" {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "权限和角色不能为空",
		}
		return reply, nil
	}
	db := common.DB

	role, err := usermodel.FindRole(db, in.Role)
	if err != nil {
		return nil, err
	}

	if role == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusNotFound,
			Message:     NotFound,
			Description: "角色不存在",
		}
		return reply, nil
	}

	permission, err := usermodel.FindPermissionByName(db, in.Permission)
	if err != nil {
		return nil, err
	}
	if permission == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusNotFound,
			Message:     NotFound,
			Description: "权限不存在",
		}
		return reply, nil
	}

	if !common.Enforcer.AddRoleForUser(fmt.Sprintf("role:%d", role.ID), fmt.Sprintf("permission:%d", permission.ID)) {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     AlreadyExists,
			Description: "策略已存在",
		}
	}
	return reply, nil
}

func (u *userService) CreateRole(ctx context.Context, in *pb.CreateRoleRequest) (*pb.CreateRoleReply, error) {
	reply := &pb.CreateRoleReply{}
	if in.Role == "" {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "角色名不能为空",
		}
		return reply, nil
	}
	db := common.DB

	r, err := usermodel.FindRole(db, in.Role)
	if err != nil {
		return nil, err
	}

	if r != nil {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     AlreadyExists,
			Description: "角色已存在",
		}
		return reply, nil
	}

	err = usermodel.SaveRole(db, &usermodel.Role{Role: in.Role})
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func (u *userService) AddRoleForUser(ctx context.Context, in *pb.AddRoleForUserRequest) (*pb.AddRoleForUserReply, error) {
	reply := &pb.AddRoleForUserReply{}
	if in.Username == "" || in.Role == "" {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "用户名和角色名不能为空",
		}
		return reply, nil
	}

	db := common.DB

	role, err := usermodel.FindRole(db, in.Role)
	if err != nil {
		return nil, err
	}
	if role == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusNotFound,
			Message:     NotFound,
			Description: "找不到角色",
		}
		return reply, nil
	}

	user, err := usermodel.FindUserByUserName(db, in.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusNotFound,
			Message:     NotFound,
			Description: "找不到用户",
		}
		return reply, nil
	}

	if !common.Enforcer.AddRoleForUser(fmt.Sprintf("user:%d", user.UserID), fmt.Sprintf("role:%d", role.ID)) {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     AlreadyExists,
			Description: "策略已存在",
		}
	}
	return reply, nil
}

func (u *userService) AddRoutes(ctx context.Context, in *pb.AddRoutesRequest) (*pb.AddRoutesReply, error) {
	reply := &pb.AddRoutesReply{}
	if len(in.Routes) == 0 {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "路由不能为空",
		}
		return reply, nil
	}
	db := common.DB

	rs, err := usermodel.FindRoutesByNames(db, in.Routes)
	if err != nil {
		return nil, err
	}

	if len(rs) != 0 {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     AlreadyExists,
			Description: "路由已存在",
		}
		return reply, nil
	}

	err = usermodel.SaveRoutes(db, in.Routes)
	if err != nil {
		return nil, err
	}
	return reply, err
}

func (u *userService) ListRoutes(ctx context.Context, in *pb.ListRoutesRequest) (*pb.ListRoutesReply, error) {
	db := common.DB

	routes, err := usermodel.ListRoutes(db)
	if err != nil {
		return nil, err
	}

	names := make([]*pb.Route, len(routes))
	for i := range routes {
		names[i] = &pb.Route{
			Id:   routes[i].ID,
			Name: routes[i].Name,
		}
	}

	return &pb.ListRoutesReply{
		Routes: names,
	}, nil
}

func (u *userService) CreatePermission(ctx context.Context, in *pb.CreatePermissionRequest) (*pb.CreatePermissionReply, error) {
	reply := &pb.CreatePermissionReply{}
	if in.Name == "" {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "权限信息不全",
		}
		return reply, nil
	}
	db := common.DB

	p, err := usermodel.FindPermissionByName(db, in.Name)
	if err != nil {
		return nil, err
	}

	if p != nil {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     AlreadyExists,
			Description: "权限已存在",
		}
		return reply, nil
	}
	err = usermodel.SavePermission(db, in.Name)
	if err != nil {
		return nil, err
	}
	return reply, err
}

func (u *userService) UpdatePermission(ctx context.Context, in *pb.UpdatePermissionRequest) (*pb.UpdatePermissionReply, error) {
	reply := &pb.UpdatePermissionReply{}
	if in.Id == 0 {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "id不能为空",
		}
		return reply, nil
	}
	db := common.DB

	p, err := usermodel.FindPermissionByID(db, in.Id)
	if err != nil {
		return nil, err
	}

	if p == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusNotFound,
			Message:     NotFound,
			Description: "权限不存在",
		}
		return reply, nil
	}
	err = usermodel.UpdatePermission(db, in.Id, &usermodel.Permission{Name: in.Name})
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func (u *userService) AddRouteForPermission(ctx context.Context, in *pb.AddRouteForPermissionRequest) (*pb.AddRouteForPermissionReply, error) {
	db := common.DB
	reply := &pb.AddRouteForPermissionReply{}
	if in.Permission == "" || in.Route == "" {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "路由和权限名不能为空",
		}
		return reply, nil
	}

	permission, err := usermodel.FindPermissionByName(db, in.Permission)
	if err != nil {
		return nil, err
	}

	if permission == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusNotFound,
			Message:     NotFound,
			Description: "权限不存在",
		}
		return reply, nil
	}

	route, err := usermodel.FindRouteByName(db, in.Route)
	if err != nil {
		return nil, err
	}

	if route == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusNotFound,
			Message:     NotFound,
			Description: "路由不存在",
		}
		return reply, nil
	}

	if !common.Enforcer.AddPolicy(fmt.Sprintf("permission:%d", permission.ID), in.Route, "*") {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     AlreadyExists,
			Description: "策略已存在",
		}
	}
	return reply, nil

}

func (u *userService) RemoveRouteForPermission(ctx context.Context, in *pb.RemoveRouteForPermissionRequest) (*pb.RemoveRouteForPermissionReply, error) {
	reply := &pb.RemoveRouteForPermissionReply{}
	if in.Permission == "" || in.Route == "" {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "路由和权限名不能为空",
		}
		return reply, nil
	}
	db := common.DB

	permission, err := usermodel.FindPermissionByName(db, in.Permission)
	if err != nil {
		return nil, err
	}

	if permission == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusNotFound,
			Message:     NotFound,
			Description: "权限不存在",
		}
		return reply, nil
	}

	route, err := usermodel.FindRouteByName(db, in.Route)
	if err != nil {
		return nil, err
	}

	if route == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusNotFound,
			Message:     NotFound,
			Description: "路由不存在",
		}
		return reply, nil
	}

	if !common.Enforcer.RemovePolicy(fmt.Sprintf("permission:%d", permission.ID), in.Route, "*") {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     NotFound,
			Description: "策略不存在",
		}
	}
	return reply, nil
}

func (u *userService) RemovePermission(ctx context.Context, in *pb.RemovePermissionRequest) (*pb.RemovePermissionReply, error) {
	reply := &pb.RemovePermissionReply{}
	if in.Permission == "" {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "权限名不能为空",
		}
		return reply, nil
	}
	db := common.DB

	permission, err := usermodel.FindPermissionByName(db, in.Permission)
	if err != nil {
		return nil, err
	}

	if permission == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusNotFound,
			Message:     NotFound,
			Description: "权限不存在",
		}
		return reply, nil
	}

	db = db.Begin()
	defer db.Rollback()
	err = usermodel.DeletePermission(db, permission)
	if err != nil {
		return nil, err
	}
	common.Enforcer.DeleteRole(fmt.Sprintf("permission:%d", permission.ID))
	common.Enforcer.DeleteUser(fmt.Sprintf("permission:%d", permission.ID))
	err = db.Commit().Error
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func (u *userService) ListPermissions(ctx context.Context, in *pb.ListPermissionsRequest) (*pb.ListPermissionsReply, error) {
	db := common.DB

	if in.Size == 0 {
		in.Size = 10
	}
	if in.Page == 0 {
		in.Page = 1
	}

	query := &usermodel.Permission{}
	if in.Permission != nil {
		query.ID = in.Permission.Id
		query.Name = in.Permission.Permission
	}
	ps, count, err := usermodel.ListPermissions(db, query, in.Page, in.Size)
	if err != nil {
		return nil, err
	}

	names := make([]*pb.PermissionField, len(ps))
	for i := range ps {
		names[i] = &pb.PermissionField{
			Id:         ps[i].ID,
			Permission: ps[i].Name,
		}
	}

	return &pb.ListPermissionsReply{
		Permissions: names,
		Page:        in.Page,
		Size:        in.Size,
		Count:       count,
	}, nil
}

func (u *userService) AddPermissionForPermission(ctx context.Context, in *pb.AddPermissionForPermissionRequest) (*pb.AddPermissionForPermissionReply, error) {
	reply := &pb.AddPermissionForPermissionReply{}
	if in.From == "" || in.Child == "" {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "权限名不能为空",
		}
		return reply, nil
	}
	db := common.DB

	from, err := usermodel.FindPermissionByName(db, in.From)
	if err != nil {
		return nil, err
	}
	if from == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusNotFound,
			Message:     NotFound,
			Description: "权限不存在",
		}
		return reply, nil
	}
	child, err := usermodel.FindPermissionByName(db, in.Child)
	if err != nil {
		return nil, err
	}
	if child == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusNotFound,
			Message:     NotFound,
			Description: "权限不存在",
		}
		return reply, nil
	}

	if !common.Enforcer.AddRoleForUser(fmt.Sprintf("permission:%d", from.ID), fmt.Sprintf("permission:%d", child.ID)) {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     AlreadyExists,
			Description: "策略已存在",
		}
	}
	return reply, nil
}

func (u *userService) RemovePermissionForPermission(ctx context.Context, in *pb.RemovePermissionForPermissionRequest) (*pb.RemovePermissionForPermissionReply, error) {
	reply := &pb.RemovePermissionForPermissionReply{}

	if in.From == "" || in.Child == "" {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "权限名不能为空",
		}
		return reply, nil
	}
	db := common.DB

	from, err := usermodel.FindPermissionByName(db, in.From)
	if err != nil {
		return nil, err
	}
	if from == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusNotFound,
			Message:     NotFound,
			Description: "权限不存在",
		}
		return reply, nil
	}
	child, err := usermodel.FindPermissionByName(db, in.Child)
	if err != nil {
		return nil, err
	}
	if child == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusNotFound,
			Message:     NotFound,
			Description: "权限不存在",
		}
		return reply, nil
	}

	if !common.Enforcer.DeleteRoleForUser(fmt.Sprintf("permission:%d", from.ID), fmt.Sprintf("permission:%d", child.ID)) {
		reply.Err = &pb.Error{
			Code:        http.StatusNotFound,
			Message:     NotFound,
			Description: "策略不存在",
		}
	}
	return reply, nil
}

func (u *userService) ListRole(ctx context.Context, in *pb.ListRoleRequest) (*pb.ListRoleReply, error) {
	db := common.DB

	if in.Size == 0 {
		in.Size = 10
	}
	if in.Page == 0 {
		in.Page = 1
	}

	query := &usermodel.Role{}

	if in.Role != nil {
		query.ID = in.Role.Id
		query.Role = in.Role.Role
	}
	roles, count, err := usermodel.ListRole(db, query, in.Page, in.Size)
	if err != nil {
		return nil, err
	}

	names := make([]*pb.RoleField, len(roles))
	for i := range roles {
		names[i] = &pb.RoleField{
			Id:   roles[i].ID,
			Role: roles[i].Role,
		}
	}
	return &pb.ListRoleReply{
		Roles: names,
		Page:  in.Page,
		Size:  in.Size,
		Count: count,
	}, nil
}

func (u *userService) UpdateRole(ctx context.Context, in *pb.UpdateRoleRequest) (*pb.UpdateRoleReply, error) {
	db := common.DB
	reply := &pb.UpdateRoleReply{}
	if in.Id == 0 {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "id不能为空",
		}
		return reply, nil
	}

	p, err := usermodel.FindRoleByID(db, in.Id)
	if err != nil {
		return nil, err
	}

	if p == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusNotFound,
			Message:     NotFound,
			Description: "角色不存在",
		}
		return reply, nil
	}

	err = usermodel.UpdateRole(db, in.Id, &usermodel.Role{Role: in.Name})
	if err != nil {
		return nil, err
	}
	return reply, err
}

func (u *userService) RemovePermissionForRole(ctx context.Context, in *pb.RemovePermissionForRoleRequest) (*pb.RemovePermissionForRoleReply, error) {
	reply := &pb.RemovePermissionForRoleReply{}
	if in.Role == "" || len(in.Permission) == 0 {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "角色名和权限不能为空",
		}
		return reply, nil
	}
	db := common.DB

	role, err := usermodel.FindRole(db, in.Role)
	if err != nil {
		return nil, err
	}

	if role == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusNotFound,
			Message:     NotFound,
			Description: "角色不存在",
		}
		return reply, nil
	}

	permission, err := usermodel.FindPermissionByName(db, in.Permission)
	if err != nil {
		return nil, err
	}
	if permission == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusNotFound,
			Message:     NotFound,
			Description: "权限不存在",
		}
		return reply, nil
	}

	if !common.Enforcer.DeleteRoleForUser(fmt.Sprintf("role:%d", role.ID), fmt.Sprintf("permission:%d", permission.ID)) {
		reply.Err = &pb.Error{
			Code:        http.StatusNotFound,
			Message:     NotFound,
			Description: "策略不存在",
		}
	}
	return reply, nil
}

func (u *userService) AddRoleForRole(ctx context.Context, in *pb.AddRoleForRoleRequest) (*pb.AddRoleForRoleReply, error) {
	reply := &pb.AddRoleForRoleReply{}
	if in.From == "" || in.Child == "" {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "角色名不能为空",
		}
		return reply, nil
	}
	db := common.DB

	from, err := usermodel.FindRole(db, in.From)
	if err != nil {
		return nil, err
	}
	if from == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusNotFound,
			Message:     NotFound,
			Description: "角色不存在",
		}
		return reply, nil
	}
	child, err := usermodel.FindRole(db, in.Child)
	if err != nil {
		return nil, err
	}
	if child == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusNotFound,
			Message:     NotFound,
			Description: "角色不存在",
		}
		return reply, nil
	}

	if !common.Enforcer.AddRoleForUser(fmt.Sprintf("role:%d", from.ID), fmt.Sprintf("role:%d", child.ID)) {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     AlreadyExists,
			Description: "策略已存在",
		}
	}
	return reply, nil
}

func (u *userService) RemoveRoleForRole(ctx context.Context, in *pb.RemoveRoleForRoleRequest) (*pb.RemoveRoleForRoleReply, error) {
	reply := &pb.RemoveRoleForRoleReply{}
	if in.From == "" || in.Child == "" {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "角色名不能为空",
		}
		return reply, nil
	}
	db := common.DB

	from, err := usermodel.FindRole(db, in.From)
	if err != nil {
		return nil, err
	}
	if from == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusNotFound,
			Message:     NotFound,
			Description: "角色不存在",
		}
		return reply, nil
	}
	child, err := usermodel.FindRole(db, in.Child)
	if err != nil {
		return nil, err
	}
	if child == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusNotFound,
			Message:     NotFound,
			Description: "角色不存在",
		}
		return reply, nil
	}

	if !common.Enforcer.DeleteRoleForUser(fmt.Sprintf("role:%d", from.ID), fmt.Sprintf("role:%d", child.ID)) {
		reply.Err = &pb.Error{
			Code:        http.StatusNotFound,
			Message:     NotFound,
			Description: "策略不存在",
		}
	}
	return reply, nil
}

func (u *userService) RemoveRole(ctx context.Context, in *pb.RemoveRoleRequest) (*pb.RemoveRoleReply, error) {
	reply := &pb.RemoveRoleReply{}
	if in.Role == "" {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "角色名不能为空",
		}
		return reply, nil
	}
	db := common.DB

	role, err := usermodel.FindRole(db, in.Role)
	if err != nil {
		return nil, err
	}

	if role == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusNotFound,
			Message:     NotFound,
			Description: "角色不存在",
		}
		return reply, nil
	}

	db = db.Begin()
	defer db.Rollback()
	err = usermodel.DeleteRole(db, role)
	if err != nil {
		return nil, err
	}
	common.Enforcer.DeleteRole(fmt.Sprintf("role:%d", role.ID))
	common.Enforcer.DeleteUser(fmt.Sprintf("role:%d", role.ID))
	err = db.Commit().Error
	if err != nil {
		return nil, err
	}
	return reply, err
}

func (u *userService) ListUsers(ctx context.Context, in *pb.ListUsersRequest) (*pb.ListUsersReply, error) {
	db := common.DB
	if in.Size == 0 {
		in.Size = 10
	}
	if in.Page == 0 {
		in.Page = 1
	}

	query := &usermodel.User{}
	if in.User != nil {
		query.UserID, _ = strconv.ParseInt(in.User.Id, 10, 64)
		query.UserName = in.User.Username
		query.LeaguerNO = in.User.LeaguerNo
		query.Email = in.User.Email
		query.UserType = in.User.UserType
		query.UserStatus, _ = strconv.ParseInt(in.User.UserStatus, 10, 64)
	}

	us, count, err := usermodel.ListUsers(db, query, in.Page, in.Size)
	if err != nil {
		return nil, err
	}

	users := make([]*pb.UserField, 0)

	for _, u := range us {
		users = append(users, &pb.UserField{
			Id:          fmt.Sprintf("%d", u.UserID),
			Username:    u.UserName,
			UserType:    u.UserType,
			LeaguerNo:   u.LeaguerNO,
			Email:       u.Email,
			UserStatus:  fmt.Sprintf("%d", u.UserStatus),
			UserGroupNo: u.UserGroupNo,
			CreatedAt:   u.CreatedAt.Format(util.TimePattern),
		})
	}

	return &pb.ListUsersReply{
		Users: users,
		Page:  in.Page,
		Size:  in.Size,
		Count: count,
	}, nil
}

func (u *userService) UpdateUser(ctx context.Context, in *pb.UpdateUserRequest) (*pb.UpdateUserReply, error) {
	db := common.DB
	reply := &pb.UpdateUserReply{}
	if in.Id == "" {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "id不能为空",
		}
		return reply, nil
	}
	userId, _ := strconv.ParseInt(in.Id, 10, 64)

	user, err := usermodel.FindUserByID(db, userId)
	if err != nil {
		return nil, err
	}

	if user == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusNotFound,
			Message:     NotFound,
			Description: "用户不存在",
		}
		return reply, nil
	}

	err = usermodel.UpdateUser(db, userId, &usermodel.User{
		UserName:    in.Username,
		Email:       in.Email,
		UserType:    in.UserType,
		UserGroupNo: in.UserGroupNo,
	})
	if err != nil {
		return nil, err
	}
	_ = cache.DelUserInfo(common.RedisClient, fmt.Sprintf("%d", user.UserID))

	return reply, nil
}

func (u *userService) AddPermissionForUser(ctx context.Context, in *pb.AddPermissionForUserRequest) (*pb.AddPermissionForUserReply, error) {
	reply := &pb.AddPermissionForUserReply{}
	if in.Username == "" || in.Permission == "" {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "用户名和权限名不能为空",
		}
		return reply, nil
	}
	db := common.DB

	user, err := usermodel.FindUserByUserName(db, in.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusNotFound,
			Message:     NotFound,
			Description: "用户不存在",
		}
		return reply, nil
	}

	permission, err := usermodel.FindPermissionByName(db, in.Permission)
	if err != nil {
		return nil, err
	}
	if permission == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusNotFound,
			Message:     NotFound,
			Description: "权限不存在",
		}
		return reply, nil
	}

	if !common.Enforcer.AddRoleForUser(fmt.Sprintf("user:%d", user.UserID), fmt.Sprintf("permission:%d", permission.ID)) {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     AlreadyExists,
			Description: "策略已存在",
		}
	}
	return reply, nil
}

func (u *userService) RemovePermissionForUser(ctx context.Context, in *pb.RemovePermissionForUserRequest) (*pb.RemovePermissionForUserReply, error) {
	reply := &pb.RemovePermissionForUserReply{}
	if in.Username == "" || in.Permission == "" {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "用户名和权限名不能为空",
		}
		return reply, nil
	}
	db := common.DB

	user, err := usermodel.FindUserByUserName(db, in.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusNotFound,
			Message:     NotFound,
			Description: "用户不存在",
		}
		return reply, nil
	}

	permission, err := usermodel.FindPermissionByName(db, in.Permission)
	if err != nil {
		return nil, err
	}
	if permission == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusNotFound,
			Message:     NotFound,
			Description: "权限不存在",
		}
		return reply, nil
	}

	if !common.Enforcer.DeleteRoleForUser(fmt.Sprintf("user:%d", user.UserID), fmt.Sprintf("permission:%d", permission.ID)) {
		reply.Err = &pb.Error{
			Code:        http.StatusNotFound,
			Message:     NotFound,
			Description: "策略不存在",
		}
		return reply, nil
	}
	return reply, nil
}

func (u *userService) RemoveRoleForUser(ctx context.Context, in *pb.RemoveRoleForUserRequest) (*pb.RemoveRoleForUserReply, error) {
	reply := &pb.RemoveRoleForUserReply{}
	if in.Username == "" || in.Role == "" {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "用户名和角色名不能为空",
		}
		return reply, nil
	}

	db := common.DB

	role, err := usermodel.FindRole(db, in.Role)
	if err != nil {
		return nil, err
	}
	if role == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusNotFound,
			Message:     NotFound,
			Description: "角色不存在",
		}
		return reply, nil
	}

	user, err := usermodel.FindUserByUserName(db, in.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusNotFound,
			Message:     NotFound,
			Description: "用户不存在",
		}
		return reply, nil
	}

	if !common.Enforcer.DeleteRoleForUser(fmt.Sprintf("user:%d", user.UserID), fmt.Sprintf("role:%d", role.ID)) {
		reply.Err = &pb.Error{
			Code:        http.StatusNotFound,
			Message:     NotFound,
			Description: "策略不存在",
		}
		return reply, nil
	}
	return reply, nil

}

func (u *userService) ListMenus(ctx context.Context, in *pb.ListMenusRequest) (*pb.ListMenusReply, error) {
	db := common.DB
	if in.Size == 0 {
		in.Size = 10
	}
	if in.Page == 0 {
		in.Page = 1
	}

	query := &usermodel.Menu{}
	if in.Menu != nil {
		query.ID = in.Menu.Id
		query.Name = in.Menu.Name
		query.Parent = in.Menu.Parent
		query.MenuData = in.Menu.Data
		query.MenuRoute = in.Menu.Route
		query.MenuOrder = in.Menu.Order
	}
	menus, count, err := usermodel.ListMenus(db, query, in.Page, in.Size)
	if err != nil {
		return nil, err
	}

	ms := make([]*pb.Menu, 0)
	for _, m := range menus {
		ms = append(ms, &pb.Menu{
			Id:     m.ID,
			Name:   m.Name,
			Parent: m.Parent,
			Route:  m.MenuRoute,
			Data:   m.MenuData,
			Order:  m.MenuOrder,
		})
	}

	return &pb.ListMenusReply{
		Menus: ms,
		Page:  in.Page,
		Size:  in.Size,
		Count: count,
	}, nil
}

func (u *userService) CreateMenu(ctx context.Context, in *pb.CreateMenuRequest) (*pb.CreateMenuReply, error) {
	reply := &pb.CreateMenuReply{}
	if in.Name == "" || in.Order == 0 || in.Parent == "" || in.Route == "" {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "名称，映射，父级名称，路由不能为空",
		}
		return reply, nil
	}

	db := common.DB

	parent, err := usermodel.FindMenuByName(db, in.Parent)
	if err != nil {
		return nil, err
	}

	if parent == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusNotFound,
			Message:     NotFound,
			Description: "父级不存在",
		}
		return reply, nil
	}

	menu, err := usermodel.FindMenuByName(db, in.Name)
	if err != nil {
		return nil, err
	}

	if menu != nil {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     AlreadyExists,
			Description: "菜单已存在",
		}
		return reply, nil
	}

	err = usermodel.SaveMenu(db, &usermodel.Menu{
		Name:      in.Name,
		Parent:    parent.ID,
		MenuOrder: in.Order,
		MenuRoute: in.Route,
		MenuData:  in.Data,
	})
	return reply, err
}
func (u *userService) RemoveMenu(ctx context.Context, in *pb.RemoveMenuRequest) (*pb.RemoveMenuReply, error) {
	reply := &pb.RemoveMenuReply{}
	if in.Id == 0 {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "id不能为空",
		}
		return reply, nil
	}

	db := common.DB

	menu, err := usermodel.FindMenuByID(db, in.Id)
	if err != nil {
		return nil, err
	}

	if menu == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusNotFound,
			Message:     NotFound,
			Description: "菜单不存在",
		}
		return reply, nil
	}
	err = usermodel.DeleteMenu(db, in.Id)
	return reply, err
}

func (u *userService) GetUserTypeInfo(ctx context.Context, in *pb.GetUserTypeInfoRequest) (*pb.GetUserTypeInfoReply, error) {
	reply := &pb.GetUserTypeInfoReply{}

	db := common.DB
	infos, err := usermodel.ListUserInfo(db)
	if err != nil {
		return nil, err
	}

	userTypeInfos := make([]*pb.UserTypeInfo, len(infos))
	for i := range infos {
		userTypeInfos[i] = &pb.UserTypeInfo{
			UserType: infos[i].UserType,
			UserDesc: infos[i].UserDesc,
			UserInfo: infos[i].UserInfo,
		}
	}

	reply.Infos = userTypeInfos
	return reply, nil
}

func (u *userService) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.GetUserReply, error) {
	reply := &pb.GetUserReply{}
	if in.Id == 0 {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "id不能为空",
		}
		return reply, nil
	}
	db := common.DB

	user, err := usermodel.FindUserByID(db, in.Id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		reply.Err = &pb.Error{
			Code:        http.StatusNotFound,
			Message:     NotFound,
			Description: "用户不存在",
		}
		return reply, nil
	}
	reply = &pb.GetUserReply{
		Id:          user.UserID,
		Username:    user.UserName,
		UserType:    user.UserType,
		UserGroupNo: user.UserGroupNo,
		Email:       user.Email,
		UserStatus:  user.UserStatus,
		CreatedAt:   user.CreatedAt.Unix(),
	}
	return reply, nil
}

func (u *userService) GetUserPermissionsAndRoles(ctx context.Context, in *pb.GetUserPermissionsAndRolesRequest) (*pb.GetUserPermissionsAndRolesReply, error) {
	reply := &pb.GetUserPermissionsAndRolesReply{}
	if in.Id == 0 {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "id不能为空",
		}
		return reply, nil
	}
	db := common.DB

	values := common.Enforcer.GetImplicitRolesForUser(fmt.Sprintf("user:%d", in.Id))

	roleIds := make([]int64, 0)
	permissionIds := make([]int64, 0)

	for _, value := range values {
		info, id := rbac.Split(value)
		if info == "role" {
			roleIds = append(roleIds, id)
		} else if info == "permission" {
			permissionIds = append(permissionIds, id)
		}
	}

	roles, err := usermodel.FindRolesByIds(db, roleIds)
	if err != nil {
		return nil, err
	}
	permissions, err := usermodel.FindPermissionsByIds(db, permissionIds)
	if err != nil {
		return nil, err
	}

	repRoles := make([]*pb.RoleField, 0, len(roles))
	repPermissions := make([]*pb.PermissionField, 0, len(permissions))

	for _, role := range roles {
		repRoles = append(repRoles, &pb.RoleField{
			Id:   role.ID,
			Role: role.Role,
		})
	}

	for _, permission := range permissions {
		repPermissions = append(repPermissions, &pb.PermissionField{
			Id:         permission.ID,
			Permission: permission.Name,
		})
	}
	reply.Roles = repRoles
	reply.Permissions = repPermissions
	return reply, nil
}

func (u *userService) GetRolePermissionsAndRoles(ctx context.Context, in *pb.GetRolePermissionsAndRolesRequest) (*pb.GetRolePermissionsAndRolesReply, error) {
	reply := &pb.GetRolePermissionsAndRolesReply{}
	if in.Id == 0 {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "id不能为空",
		}
		return reply, nil
	}
	db := common.DB

	values := common.Enforcer.GetImplicitRolesForUser(fmt.Sprintf("role:%d", in.Id))

	roleIds := make([]int64, 0)
	permissionIds := make([]int64, 0)

	for _, value := range values {
		info, id := rbac.Split(value)
		if info == "role" {
			roleIds = append(roleIds, id)
		} else if info == "permission" {
			permissionIds = append(permissionIds, id)
		}
	}

	roles, err := usermodel.FindRolesByIds(db, roleIds)
	if err != nil {
		return nil, err
	}
	permissions, err := usermodel.FindPermissionsByIds(db, permissionIds)
	if err != nil {
		return nil, err
	}

	repRoles := make([]*pb.RoleField, 0, len(roles))
	repPermissions := make([]*pb.PermissionField, 0, len(permissions))

	for _, role := range roles {
		repRoles = append(repRoles, &pb.RoleField{
			Id:   role.ID,
			Role: role.Role,
		})
	}

	for _, permission := range permissions {
		repPermissions = append(repPermissions, &pb.PermissionField{
			Id:         permission.ID,
			Permission: permission.Name,
		})
	}
	reply.Roles = repRoles
	reply.Permissions = repPermissions
	return reply, nil
}

func (u *userService) GetPermissionsAndRoutes(ctx context.Context, in *pb.GetPermissionsAndRoutesRequest) (*pb.GetPermissionsAndRoutesReply, error) {
	reply := &pb.GetPermissionsAndRoutesReply{}
	if in.Id == 0 {
		reply.Err = &pb.Error{
			Code:        http.StatusBadRequest,
			Message:     InvalidParam,
			Description: "id不能为空",
		}
		return reply, nil
	}
	db := common.DB

	values := common.Enforcer.GetImplicitRolesForUser(fmt.Sprintf("permission:%d", in.Id))
	policy := common.Enforcer.GetImplicitPermissionsForUser(fmt.Sprintf("permission:%d", in.Id))
	routes := make([]string, 0)
	for _, p := range policy {
		if len(p) >= 2 {
			routes = append(routes, p[1])
		}
	}

	permissionIds := make([]int64, 0)

	for _, value := range values {
		info, id := rbac.Split(value)
		if info == "permission" {
			permissionIds = append(permissionIds, id)
		}
	}

	permissions, err := usermodel.FindPermissionsByIds(db, permissionIds)
	if err != nil {
		return nil, err
	}

	repPermissions := make([]*pb.PermissionField, 0, len(permissions))

	for _, permission := range permissions {
		repPermissions = append(repPermissions, &pb.PermissionField{
			Id:         permission.ID,
			Permission: permission.Name,
		})
	}
	reply.Permissions = repPermissions
	reply.Routes = routes
	return reply, nil
}

func (u *userService) ListLeaguer(ctx context.Context, in *pb.ListLeaguerRequest) (*pb.ListLeaguerReply, error) {
	db := common.DB
	if in.Size == 0 {
		in.Size = 10
	}
	if in.Page == 0 {
		in.Page = 1
	}

	query := &usermodel.Leaguer{}
	if in.Leaguer != nil {
		query.LeaguerNo = in.Leaguer.LeaguerNo
		query.LeaguerName = in.Leaguer.LeaguerName
		query.LeaguerType = in.Leaguer.LeaguerType
		query.LeaguerInfo = in.Leaguer.LeaguerInfo
		query.LeaguerStatus = in.Leaguer.LeaguerStatus
	}
	leaguers, count, err := usermodel.ListLeaguers(db, query, in.Page, in.Size)
	if err != nil {
		return nil, err
	}

	ms := make([]*pb.LeaguerField, 0)
	for _, m := range leaguers {
		ms = append(ms, &pb.LeaguerField{
			LeaguerNo:     m.LeaguerNo,
			LeaguerName:   m.LeaguerName,
			LeaguerType:   m.LeaguerType,
			LeaguerInfo:   m.LeaguerInfo,
			LeaguerStatus: m.LeaguerStatus,
			Created:       int64(m.Created.Unix()),
			Updated:       int64(m.Updated.Unix()),
		})
	}

	return &pb.ListLeaguerReply{
		Leaguers: ms,
		Page:     in.Page,
		Size:     in.Size,
		Count:    count,
	}, nil
}
