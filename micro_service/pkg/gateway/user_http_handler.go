package gateway

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"
	"userService/pkg/common"
	"userService/pkg/pb"
	"userService/pkg/userservice"

	"github.com/dgrijalva/jwt-go"
	stdjwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/sd/lb"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

func RegisterUserHandler(engine *gin.Engine, endpoints *UserEndpoints) {
	userGroup := engine.Group("/user")

	userGroup.POST("/login", convertHttpHandlerToGinHandler(httptransport.NewServer(
		endpoints.LoginEndpoint,
		decodeHttpRequest(&pb.LoginRequest{}),
		encodeHttpResponse,
		httptransport.ServerErrorEncoder(errorEncoder),
	)))

	userGroup.GET("/getPermissions",
		userservice.JwtMiddleware(keyFunc, stdjwt.SigningMethodHS256, userservice.UserClaimFactory),
		convertHttpHandlerToGinHandler(httptransport.NewServer(
			endpoints.GetPermissionsEndpoint,
			decodeHttpRequest(&pb.GetPermissionsRequest{}),
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(errorEncoder),
			httptransport.ServerBefore(setUserInfoContext),
		)))

	userGroup.POST("/checkPermission",
		userservice.JwtMiddleware(keyFunc, stdjwt.SigningMethodHS256, userservice.UserClaimFactory),
		convertHttpHandlerToGinHandler(httptransport.NewServer(
			endpoints.CheckPermissionEndpoint,
			decodeHttpRequest(&pb.CheckPermissionRequest{}),
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(errorEncoder),
			httptransport.ServerBefore(setUserInfoContext),
		)))

	userGroup.POST("/register", convertHttpHandlerToGinHandler(httptransport.NewServer(
		endpoints.RegisterEndpoint,
		decodeHttpRequest(&pb.RegisterRequest{}),
		encodeHttpResponse,
		httptransport.ServerErrorEncoder(errorEncoder),
	)))

	userGroup.POST("/addPermissionForRole",
		userservice.JwtMiddleware(keyFunc, stdjwt.SigningMethodHS256, userservice.UserClaimFactory),
		convertHttpHandlerToGinHandler(httptransport.NewServer(
			endpoints.AddPermissionForRoleEndpoint,
			decodeHttpRequest(&pb.AddPermissionForRoleRequest{}),
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(errorEncoder),
		)))

	userGroup.POST("/createRole",
		userservice.JwtMiddleware(keyFunc, stdjwt.SigningMethodHS256, userservice.UserClaimFactory),
		convertHttpHandlerToGinHandler(httptransport.NewServer(
			endpoints.CreateRoleEndpoint,
			decodeHttpRequest(&pb.CreateRoleRequest{}),
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(errorEncoder),
		)))

	userGroup.POST("/addRoleForUser",
		userservice.JwtMiddleware(keyFunc, stdjwt.SigningMethodHS256, userservice.UserClaimFactory),
		convertHttpHandlerToGinHandler(httptransport.NewServer(
			endpoints.AddRoleForUserEndpoint,
			decodeHttpRequest(&pb.AddRoleForUserRequest{}),
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(errorEncoder),
		)))

	userGroup.POST("/addRoutes",
		userservice.JwtMiddleware(keyFunc, stdjwt.SigningMethodHS256, userservice.UserClaimFactory),
		convertHttpHandlerToGinHandler(httptransport.NewServer(
			endpoints.AddRoutesEndpoint,
			decodeHttpRequest(&pb.AddRoutesRequest{}),
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(errorEncoder),
		)))

	userGroup.GET("/listRoutes",
		userservice.JwtMiddleware(keyFunc, stdjwt.SigningMethodHS256, userservice.UserClaimFactory),
		convertHttpHandlerToGinHandler(httptransport.NewServer(
			endpoints.ListRoutesEndpoint,
			decodeHttpRequest(&pb.ListRoutesRequest{}),
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(errorEncoder),
		)))

	userGroup.POST("/createPermission",
		userservice.JwtMiddleware(keyFunc, stdjwt.SigningMethodHS256, userservice.UserClaimFactory),
		convertHttpHandlerToGinHandler(httptransport.NewServer(
			endpoints.CreatePermissionEndpoint,
			decodeHttpRequest(&pb.CreatePermissionRequest{}),
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(errorEncoder),
		)))

	userGroup.POST("/updatePermission",
		userservice.JwtMiddleware(keyFunc, stdjwt.SigningMethodHS256, userservice.UserClaimFactory),
		convertHttpHandlerToGinHandler(httptransport.NewServer(
			endpoints.UpdatePermissionEndpoint,
			decodeHttpRequest(&pb.UpdatePermissionRequest{}),
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(errorEncoder),
		)))

	userGroup.POST("/addRouteForPermission",
		userservice.JwtMiddleware(keyFunc, stdjwt.SigningMethodHS256, userservice.UserClaimFactory),
		convertHttpHandlerToGinHandler(httptransport.NewServer(
			endpoints.AddRouteForPermissionEndpoint,
			decodeHttpRequest(&pb.AddRouteForPermissionRequest{}),
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(errorEncoder),
		)))

	userGroup.POST("/removeRouteForPermission",
		userservice.JwtMiddleware(keyFunc, stdjwt.SigningMethodHS256, userservice.UserClaimFactory),
		convertHttpHandlerToGinHandler(httptransport.NewServer(
			endpoints.RemoveRouteForPermissionEndpoint,
			decodeHttpRequest(&pb.RemoveRouteForPermissionRequest{}),
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(errorEncoder),
		)))

	userGroup.POST("/removePermission",
		userservice.JwtMiddleware(keyFunc, stdjwt.SigningMethodHS256, userservice.UserClaimFactory),
		convertHttpHandlerToGinHandler(httptransport.NewServer(
			endpoints.RemovePermissionEndpoint,
			decodeHttpRequest(&pb.RemovePermissionRequest{}),
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(errorEncoder),
		)))

	userGroup.POST("/listPermissions",
		userservice.JwtMiddleware(keyFunc, stdjwt.SigningMethodHS256, userservice.UserClaimFactory),
		convertHttpHandlerToGinHandler(httptransport.NewServer(
			endpoints.ListPermissionsEndpoint,
			decodeHttpRequest(&pb.ListPermissionsRequest{}),
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(errorEncoder),
		)))

	userGroup.POST("/addPermissionForPermission",
		userservice.JwtMiddleware(keyFunc, stdjwt.SigningMethodHS256, userservice.UserClaimFactory),
		convertHttpHandlerToGinHandler(httptransport.NewServer(
			endpoints.AddPermissionForPermissionEndpoint,
			decodeHttpRequest(&pb.AddPermissionForPermissionRequest{}),
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(errorEncoder),
		)))

	userGroup.POST("/removePermissionForPermission",
		userservice.JwtMiddleware(keyFunc, stdjwt.SigningMethodHS256, userservice.UserClaimFactory),
		convertHttpHandlerToGinHandler(httptransport.NewServer(
			endpoints.RemovePermissionForPermissionEndpoint,
			decodeHttpRequest(&pb.RemovePermissionForPermissionRequest{}),
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(errorEncoder),
		)))

	userGroup.POST("/listRole",
		userservice.JwtMiddleware(keyFunc, stdjwt.SigningMethodHS256, userservice.UserClaimFactory),
		convertHttpHandlerToGinHandler(httptransport.NewServer(
			endpoints.ListRoleEndpoint,
			decodeHttpRequest(&pb.ListRoleRequest{}),
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(errorEncoder),
		)))

	userGroup.POST("/updateRole",
		userservice.JwtMiddleware(keyFunc, stdjwt.SigningMethodHS256, userservice.UserClaimFactory),
		convertHttpHandlerToGinHandler(httptransport.NewServer(
			endpoints.UpdateRoleEndpoint,
			decodeHttpRequest(&pb.UpdateRoleRequest{}),
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(errorEncoder),
		)))

	userGroup.POST("/removePermissionForRole",
		userservice.JwtMiddleware(keyFunc, stdjwt.SigningMethodHS256, userservice.UserClaimFactory),
		convertHttpHandlerToGinHandler(httptransport.NewServer(
			endpoints.RemovePermissionForRoleEndpoint,
			decodeHttpRequest(&pb.RemovePermissionForRoleRequest{}),
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(errorEncoder),
		)))

	userGroup.POST("/addRoleForRole",
		userservice.JwtMiddleware(keyFunc, stdjwt.SigningMethodHS256, userservice.UserClaimFactory),
		convertHttpHandlerToGinHandler(httptransport.NewServer(
			endpoints.AddRoleForRoleEndpoint,
			decodeHttpRequest(&pb.AddRoleForRoleRequest{}),
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(errorEncoder),
		)))

	userGroup.POST("/removeRoleForRole",
		userservice.JwtMiddleware(keyFunc, stdjwt.SigningMethodHS256, userservice.UserClaimFactory),
		convertHttpHandlerToGinHandler(httptransport.NewServer(
			endpoints.RemoveRoleForRoleEndpoint,
			decodeHttpRequest(&pb.RemoveRoleForRoleRequest{}),
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(errorEncoder),
		)))

	userGroup.POST("/removeRole",
		userservice.JwtMiddleware(keyFunc, stdjwt.SigningMethodHS256, userservice.UserClaimFactory),
		convertHttpHandlerToGinHandler(httptransport.NewServer(
			endpoints.RemoveRoleEndpoint,
			decodeHttpRequest(&pb.RemoveRoleRequest{}),
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(errorEncoder),
		)))

	userGroup.POST("/listUsers",
		userservice.JwtMiddleware(keyFunc, stdjwt.SigningMethodHS256, userservice.UserClaimFactory),
		convertHttpHandlerToGinHandler(httptransport.NewServer(
			endpoints.ListUsersEndpoint,
			decodeHttpRequest(&pb.ListUsersRequest{}),
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(errorEncoder),
		)))

	userGroup.POST("/updateUser",
		userservice.JwtMiddleware(keyFunc, stdjwt.SigningMethodHS256, userservice.UserClaimFactory),
		convertHttpHandlerToGinHandler(httptransport.NewServer(
			endpoints.UpdateUserEndpoint,
			decodeHttpRequest(&pb.UpdateUserRequest{}),
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(errorEncoder),
		)))

	userGroup.POST("/addPermissionForUser",
		userservice.JwtMiddleware(keyFunc, stdjwt.SigningMethodHS256, userservice.UserClaimFactory),
		convertHttpHandlerToGinHandler(httptransport.NewServer(
			endpoints.AddPermissionForUserEndpoint,
			decodeHttpRequest(&pb.AddPermissionForUserRequest{}),
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(errorEncoder),
		)))

	userGroup.POST("/removePermissionForUser",
		userservice.JwtMiddleware(keyFunc, stdjwt.SigningMethodHS256, userservice.UserClaimFactory),
		convertHttpHandlerToGinHandler(httptransport.NewServer(
			endpoints.RemovePermissionForUserEndpoint,
			decodeHttpRequest(&pb.RemovePermissionForUserRequest{}),
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(errorEncoder),
		)))

	userGroup.POST("/removeRoleForUser",
		userservice.JwtMiddleware(keyFunc, stdjwt.SigningMethodHS256, userservice.UserClaimFactory),
		convertHttpHandlerToGinHandler(httptransport.NewServer(
			endpoints.RemoveRoleForUserEndpoint,
			decodeHttpRequest(&pb.RemoveRoleForUserRequest{}),
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(errorEncoder),
		)))

	userGroup.POST("/listMenus",
		userservice.JwtMiddleware(keyFunc, stdjwt.SigningMethodHS256, userservice.UserClaimFactory),
		convertHttpHandlerToGinHandler(httptransport.NewServer(
			endpoints.ListMenusEndpoint,
			decodeHttpRequest(&pb.ListMenusRequest{}),
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(errorEncoder),
		)))

	userGroup.POST("/createMenu",
		userservice.JwtMiddleware(keyFunc, stdjwt.SigningMethodHS256, userservice.UserClaimFactory),
		convertHttpHandlerToGinHandler(httptransport.NewServer(
			endpoints.CreateMenuEndpoint,
			decodeHttpRequest(&pb.CreateMenuRequest{}),
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(errorEncoder),
		)))

	userGroup.POST("/removeMenu",
		userservice.JwtMiddleware(keyFunc, stdjwt.SigningMethodHS256, userservice.UserClaimFactory),
		convertHttpHandlerToGinHandler(httptransport.NewServer(
			endpoints.RemoveMenuEndpoint,
			decodeHttpRequest(&pb.RemoveMenuRequest{}),
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(errorEncoder),
		)))

	userGroup.GET("/getUserTypeInfo",
		userservice.JwtMiddleware(keyFunc, stdjwt.SigningMethodHS256, userservice.UserClaimFactory),
		convertHttpHandlerToGinHandler(httptransport.NewServer(
			endpoints.GetUserTypeInfoEndpoint,
			decodeHttpRequest(&pb.GetUserTypeInfoRequest{}),
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(errorEncoder),
		)))

	userGroup.GET("/getUser",
		userservice.JwtMiddleware(keyFunc, stdjwt.SigningMethodHS256, userservice.UserClaimFactory),
		convertHttpHandlerToGinHandler(httptransport.NewServer(
			endpoints.GetUserEndpoint,
			decodeGetUserRequest,
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(errorEncoder),
		)))

	userGroup.POST("/getUserPermissionsAndRoles",
		userservice.JwtMiddleware(keyFunc, stdjwt.SigningMethodHS256, userservice.UserClaimFactory),
		convertHttpHandlerToGinHandler(httptransport.NewServer(
			endpoints.GetUserPermissionsAndRolesEndpoint,
			decodeHttpRequest(&pb.GetUserPermissionsAndRolesRequest{}),
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(errorEncoder),
		)))

	userGroup.POST("/getRolePermissionsAndRoles",
		userservice.JwtMiddleware(keyFunc, stdjwt.SigningMethodHS256, userservice.UserClaimFactory),
		convertHttpHandlerToGinHandler(httptransport.NewServer(
			endpoints.GetRolePermissionsAndRolesEndpoint,
			decodeHttpRequest(&pb.GetRolePermissionsAndRolesRequest{}),
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(errorEncoder),
		)))

	userGroup.POST("/getPermissionsAndRoutes",
		userservice.JwtMiddleware(keyFunc, stdjwt.SigningMethodHS256, userservice.UserClaimFactory),
		convertHttpHandlerToGinHandler(httptransport.NewServer(
			endpoints.GetPermissionsAndRoutesEndpoint,
			decodeHttpRequest(&pb.GetPermissionsAndRoutesRequest{}),
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(errorEncoder),
		)))

	userGroup.POST("/listLeaguer",
		userservice.JwtMiddleware(keyFunc, stdjwt.SigningMethodHS256, userservice.UserClaimFactory),
		convertHttpHandlerToGinHandler(httptransport.NewServer(
			endpoints.ListLeaguerEndpoint,
			decodeHttpRequest(&pb.ListLeaguerRequest{}),
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(errorEncoder),
		)))

	userGroup.POST("/removeRoute",
		userservice.JwtMiddleware(keyFunc, stdjwt.SigningMethodHS256, userservice.UserClaimFactory),
		convertHttpHandlerToGinHandler(httptransport.NewServer(
			endpoints.RemoveRouteEndpoint,
			decodeHttpRequest(&pb.RemoveRouteRequest{}),
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(errorEncoder),
		)))
}

func convertHttpHandlerToGinHandler(handler http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}

func decodeHttpRequest(ins interface{}) httptransport.DecodeRequestFunc {
	tp := reflect.TypeOf(ins)
	if tp.Kind() == reflect.Ptr {
		tp = tp.Elem()
	}
	return func(_ context.Context, r *http.Request) (interface{}, error) {
		request := reflect.New(tp).Interface()
		if r.Method == http.MethodPost {
			defer r.Body.Close()
			err := json.NewDecoder(r.Body).Decode(&request)
			return request, err
		} else {
			return request, nil
		}
	}
}

func encodeHttpResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	s, ok := response.(StatusError)
	if ok {
		err := s.GetErr()
		if err != nil {
			w.WriteHeader(int(err.GetCode()))
			return json.NewEncoder(w).Encode(gin.H{
				"err":  err.GetMessage(),
				"desc": err.GetDescription(),
			})
		}
	}
	pMsg, ok := response.(proto.Message)
	if ok {
		marshaler := jsonpb.Marshaler{
			EmitDefaults: true,
		}
		return marshaler.Marshal(w, pMsg)
	} else {
		return json.NewEncoder(w).Encode(response)
	}
}

func errorEncoder(ctx context.Context, err error, w http.ResponseWriter) {
	code, msg := err2codeAndMessage(err)
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(gin.H{
		"err":  msg,
		"desc": "未知错误",
	})
}

func err2codeAndMessage(err error) (int, string) {
	if err == nil {
		return http.StatusOK, ""
	}
	switch e := err.(type) {
	case lb.RetryError:
		err = e.Final
	}
	return http.StatusInternalServerError, err.Error()
}

func keyFunc(token *jwt.Token) (interface{}, error) {
	return []byte(common.SignKey), nil
}

func setUserInfoContext(ctx context.Context, r *http.Request) context.Context {
	id := r.Form.Get("userid")
	return context.WithValue(ctx, "userid", id)
}

func decodeGetUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)
	return &pb.GetUserRequest{
		Id: int64(id),
	}, nil
}

func decodeListUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	query := r.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	size, _ := strconv.Atoi(query.Get("size"))
	id := query.Get("id")
	userStatus := query.Get("userStatus")
	createdAt := query.Get("createdAt")
	leaguerNo := query.Get("leaguerNo")
	username := query.Get("username")
	email := query.Get("email")
	userType := query.Get("userType")
	return &pb.ListUsersRequest{
		Page: int32(page),
		Size: int32(size),
		User: &pb.UserField{
			Id:         id,
			LeaguerNo:  leaguerNo,
			Username:   username,
			Email:      email,
			UserType:   userType,
			UserStatus: userStatus,
			CreatedAt:  createdAt,
		},
	}, nil
}
