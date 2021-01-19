package userservice

import (
	"context"
	"time"
	"userService/pkg/pb"

	"github.com/go-kit/kit/ratelimit"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"golang.org/x/time/rate"
)

type userServer struct {
	loginHandler                         grpctransport.Handler
	getPermissionsHandler                grpctransport.Handler
	checkPermissionHandler               grpctransport.Handler
	registerHandler                      grpctransport.Handler
	addPermissionForRoleHandler          grpctransport.Handler
	addRoleForRoleHandler                grpctransport.Handler
	createRoleHandler                    grpctransport.Handler
	addRoleForUserHandler                grpctransport.Handler
	addRoutesHandler                     grpctransport.Handler
	listRouteHandler                     grpctransport.Handler
	createPermissionHandler              grpctransport.Handler
	updatePermissionHandler              grpctransport.Handler
	addRouteForPermissionHandler         grpctransport.Handler
	removeRouteForPermissionHandler      grpctransport.Handler
	removePermissionHandler              grpctransport.Handler
	listPermissionsHandler               grpctransport.Handler
	addPermissionForPermissionHandler    grpctransport.Handler
	removePermissionForPermissionHandler grpctransport.Handler
	listRoleHandler                      grpctransport.Handler
	updateRoleHandler                    grpctransport.Handler
	removePermissionForRoleHandler       grpctransport.Handler
	removeRoleForRoleHandler             grpctransport.Handler
	removeRoleHandler                    grpctransport.Handler
	listUsersHandler                     grpctransport.Handler
	updateUserHandler                    grpctransport.Handler
	addPermissionForUserHandler          grpctransport.Handler
	removePermissionForUserHandler       grpctransport.Handler
	removeRoleForUserHandler             grpctransport.Handler
	listMenusHandler                     grpctransport.Handler
	createMenuHandler                    grpctransport.Handler
	removeMenuHandler                    grpctransport.Handler
	getUserTypeInfoHandler               grpctransport.Handler
	getUserHandler                       grpctransport.Handler
	getUserPermissionsAndRolesHandler    grpctransport.Handler
	getRolePermissionsAndRolesHandler    grpctransport.Handler
	getPermissionsAndRoutesHandler       grpctransport.Handler
	listLeaguerHandler                   grpctransport.Handler
	RemoveRouteHandler                   grpctransport.Handler
}

func (u *userServer) RemoveRoute(ctx context.Context, in *pb.RemoveRouteRequest) (*pb.RemoveRouteReply, error) {
	_, res, err := u.RemoveRouteHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.RemoveRouteReply), nil
}

func New(tracer grpctransport.ServerOption) pb.UserServer {
	limiter := rate.NewLimiter(rate.Every(time.Second), 10000)
	limitMiddleware := ratelimit.NewErroringLimiter(limiter)

	svr := &userServer{}
	userService := &userService{}
	options := make([]grpctransport.ServerOption, 0)
	if tracer != nil {
		options = append(options, tracer)
	}

	{
		endpoint := MakeRemoveRouteEndpoint(userService)
		endpoint = limitMiddleware(endpoint)
		endpoint = logginMiddleware(endpoint)
		svr.RemoveRouteHandler = grpctransport.NewServer(
			endpoint,
			decodeRequest,
			encodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeLoginEndpoint(userService)
		endpoint = limitMiddleware(endpoint)
		endpoint = logginMiddleware(endpoint)
		svr.loginHandler = grpctransport.NewServer(
			endpoint,
			decodeRequest,
			encodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeGetPermissionsEndpoint(userService)
		endpoint = limitMiddleware(endpoint)
		endpoint = logginMiddleware(endpoint)
		svr.getPermissionsHandler = grpctransport.NewServer(
			endpoint,
			decodeRequest,
			encodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeCheckPermissionEndpoint(userService)
		endpoint = limitMiddleware(endpoint)
		endpoint = logginMiddleware(endpoint)
		svr.checkPermissionHandler = grpctransport.NewServer(
			endpoint,
			decodeRequest,
			encodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeRegisterEndpoint(userService)
		endpoint = limitMiddleware(endpoint)
		endpoint = logginMiddleware(endpoint)
		svr.registerHandler = grpctransport.NewServer(
			endpoint,
			decodeRequest,
			encodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeAddPermissionForRoleEndpoint(userService)
		endpoint = limitMiddleware(endpoint)
		endpoint = logginMiddleware(endpoint)
		svr.addPermissionForRoleHandler = grpctransport.NewServer(
			endpoint,
			decodeRequest,
			encodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeAddRoleForRoleEndpoint(userService)
		endpoint = limitMiddleware(endpoint)
		endpoint = logginMiddleware(endpoint)
		svr.addRoleForRoleHandler = grpctransport.NewServer(
			endpoint,
			decodeRequest,
			encodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeCreateRoleEndpoint(userService)
		endpoint = limitMiddleware(endpoint)
		endpoint = logginMiddleware(endpoint)
		svr.createRoleHandler = grpctransport.NewServer(
			endpoint,
			decodeRequest,
			encodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeAddRoleForUserEndpoint(userService)
		endpoint = limitMiddleware(endpoint)
		endpoint = logginMiddleware(endpoint)
		svr.addRoleForUserHandler = grpctransport.NewServer(
			endpoint,
			decodeRequest,
			encodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeAddRoutesEndpoint(userService)
		endpoint = limitMiddleware(endpoint)
		endpoint = logginMiddleware(endpoint)
		svr.addRoutesHandler = grpctransport.NewServer(
			endpoint,
			decodeRequest,
			encodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeListRoutesEndpoint(userService)
		endpoint = limitMiddleware(endpoint)
		endpoint = logginMiddleware(endpoint)
		svr.listRouteHandler = grpctransport.NewServer(
			endpoint,
			decodeRequest,
			encodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeCreatePermissionEndpoint(userService)
		endpoint = limitMiddleware(endpoint)
		endpoint = logginMiddleware(endpoint)
		svr.createPermissionHandler = grpctransport.NewServer(
			endpoint,
			decodeRequest,
			encodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeUpdatePermissionEndpoint(userService)
		endpoint = limitMiddleware(endpoint)
		endpoint = logginMiddleware(endpoint)
		svr.updatePermissionHandler = grpctransport.NewServer(
			endpoint,
			decodeRequest,
			encodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeAddRouteForPermissionEndpoint(userService)
		endpoint = limitMiddleware(endpoint)
		endpoint = logginMiddleware(endpoint)
		svr.addRouteForPermissionHandler = grpctransport.NewServer(
			endpoint,
			decodeRequest,
			encodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeRemoveRouteForPermissionEndpoint(userService)
		endpoint = limitMiddleware(endpoint)
		endpoint = logginMiddleware(endpoint)
		svr.removeRouteForPermissionHandler = grpctransport.NewServer(
			endpoint,
			decodeRequest,
			encodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeRemovePermissionEndpoint(userService)
		endpoint = limitMiddleware(endpoint)
		endpoint = logginMiddleware(endpoint)
		svr.removePermissionHandler = grpctransport.NewServer(
			endpoint,
			decodeRequest,
			encodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeListPermissionsEndpoint(userService)
		endpoint = limitMiddleware(endpoint)
		endpoint = logginMiddleware(endpoint)
		svr.listPermissionsHandler = grpctransport.NewServer(
			endpoint,
			decodeRequest,
			encodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeAddPermissionForPermissionEndpoint(userService)
		endpoint = limitMiddleware(endpoint)
		endpoint = logginMiddleware(endpoint)
		svr.addPermissionForPermissionHandler = grpctransport.NewServer(
			endpoint,
			decodeRequest,
			encodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeRemovePermissionForPermissionEndpoint(userService)
		endpoint = limitMiddleware(endpoint)
		endpoint = logginMiddleware(endpoint)
		svr.removePermissionForPermissionHandler = grpctransport.NewServer(
			endpoint,
			decodeRequest,
			encodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeListRoleEndpoint(userService)
		endpoint = limitMiddleware(endpoint)
		endpoint = logginMiddleware(endpoint)
		svr.listRoleHandler = grpctransport.NewServer(
			endpoint,
			decodeRequest,
			encodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeUpdateRoleEndpoint(userService)
		endpoint = limitMiddleware(endpoint)
		endpoint = logginMiddleware(endpoint)
		svr.updateRoleHandler = grpctransport.NewServer(
			endpoint,
			decodeRequest,
			encodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeRemovePermissionForRoleEndpoint(userService)
		endpoint = limitMiddleware(endpoint)
		endpoint = logginMiddleware(endpoint)
		svr.removePermissionForRoleHandler = grpctransport.NewServer(
			endpoint,
			decodeRequest,
			encodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeRemoveRoleForRoleEndpoint(userService)
		endpoint = limitMiddleware(endpoint)
		endpoint = logginMiddleware(endpoint)
		svr.removeRoleForRoleHandler = grpctransport.NewServer(
			endpoint,
			decodeRequest,
			encodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeRemoveRoleEndpoint(userService)
		endpoint = limitMiddleware(endpoint)
		endpoint = logginMiddleware(endpoint)
		svr.removeRoleHandler = grpctransport.NewServer(
			endpoint,
			decodeRequest,
			encodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeListUsersEndpoint(userService)
		endpoint = limitMiddleware(endpoint)
		endpoint = logginMiddleware(endpoint)
		svr.listUsersHandler = grpctransport.NewServer(
			endpoint,
			decodeRequest,
			encodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeUpdateUserEndpoint(userService)
		endpoint = limitMiddleware(endpoint)
		endpoint = logginMiddleware(endpoint)
		svr.updateUserHandler = grpctransport.NewServer(
			endpoint,
			decodeRequest,
			encodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeAddPermissionForUserEndpoint(userService)
		endpoint = limitMiddleware(endpoint)
		endpoint = logginMiddleware(endpoint)
		svr.addPermissionForUserHandler = grpctransport.NewServer(
			endpoint,
			decodeRequest,
			encodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeRemovePermissionForUserEndpoint(userService)
		endpoint = limitMiddleware(endpoint)
		endpoint = logginMiddleware(endpoint)
		svr.removePermissionForUserHandler = grpctransport.NewServer(
			endpoint,
			decodeRequest,
			encodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeRemoveRoleForUserEndpoint(userService)
		endpoint = limitMiddleware(endpoint)
		endpoint = logginMiddleware(endpoint)
		svr.removeRoleForUserHandler = grpctransport.NewServer(
			endpoint,
			decodeRequest,
			encodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeListMenusEndpoint(userService)
		endpoint = limitMiddleware(endpoint)
		endpoint = logginMiddleware(endpoint)
		svr.listMenusHandler = grpctransport.NewServer(
			endpoint,
			decodeRequest,
			encodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeCreateMenuEndpoint(userService)
		endpoint = limitMiddleware(endpoint)
		endpoint = logginMiddleware(endpoint)
		svr.createMenuHandler = grpctransport.NewServer(
			endpoint,
			decodeRequest,
			encodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeRemoveMenuEndpoint(userService)
		endpoint = limitMiddleware(endpoint)
		endpoint = logginMiddleware(endpoint)
		svr.removeMenuHandler = grpctransport.NewServer(
			endpoint,
			decodeRequest,
			encodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeGetUserTypeInfoEndpoint(userService)
		endpoint = limitMiddleware(endpoint)
		endpoint = logginMiddleware(endpoint)
		svr.getUserTypeInfoHandler = grpctransport.NewServer(
			endpoint,
			decodeRequest,
			encodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeGetUserEndpoint(userService)
		endpoint = limitMiddleware(endpoint)
		endpoint = logginMiddleware(endpoint)
		svr.getUserHandler = grpctransport.NewServer(
			endpoint,
			decodeRequest,
			encodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeGetUserPermissionsAndRolesEndpoint(userService)
		endpoint = limitMiddleware(endpoint)
		endpoint = logginMiddleware(endpoint)
		svr.getUserPermissionsAndRolesHandler = grpctransport.NewServer(
			endpoint,
			decodeRequest,
			encodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeGetRolePermissionsAndRolesEndpoint(userService)
		endpoint = limitMiddleware(endpoint)
		endpoint = logginMiddleware(endpoint)
		svr.getRolePermissionsAndRolesHandler = grpctransport.NewServer(
			endpoint,
			decodeRequest,
			encodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeGetPermissionsAndRoutesEndpoint(userService)
		endpoint = limitMiddleware(endpoint)
		endpoint = logginMiddleware(endpoint)
		svr.getPermissionsAndRoutesHandler = grpctransport.NewServer(
			endpoint,
			decodeRequest,
			encodeResponse,
			options...,
		)
	}

	{
		endpoint := MakeListLeaguerEndpoint(userService)
		endpoint = limitMiddleware(endpoint)
		endpoint = logginMiddleware(endpoint)
		svr.listLeaguerHandler = grpctransport.NewServer(
			endpoint,
			decodeRequest,
			encodeResponse,
			options...,
		)
	}

	return svr
}

func (u *userServer) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginReply, error) {
	_, res, err := u.loginHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	reply, ok := res.(*pb.LoginReply)
	if !ok {
		return nil, ErrReplyTypeInvalid
	}
	return reply, nil
}

func (u *userServer) GetPermissions(ctx context.Context, in *pb.GetPermissionsRequest) (*pb.GetPermissionsReply, error) {
	_, res, err := u.getPermissionsHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	reply, ok := res.(*pb.GetPermissionsReply)
	if !ok {
		return nil, ErrReplyTypeInvalid
	}
	return reply, nil
}

func (u *userServer) CheckPermission(ctx context.Context, in *pb.CheckPermissionRequest) (*pb.CheckPermissionReply, error) {
	_, res, err := u.checkPermissionHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	reply, ok := res.(*pb.CheckPermissionReply)
	if !ok {
		return nil, ErrReplyTypeInvalid
	}
	return reply, nil
}

func (u *userServer) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterReply, error) {
	_, res, err := u.registerHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	reply, ok := res.(*pb.RegisterReply)
	if !ok {
		return nil, ErrReplyTypeInvalid
	}
	return reply, nil
}

func (u *userServer) AddPermissionForRole(ctx context.Context, in *pb.AddPermissionForRoleRequest) (*pb.AddPermissionForRoleReply, error) {
	_, res, err := u.addPermissionForRoleHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	reply, ok := res.(*pb.AddPermissionForRoleReply)
	if !ok {
		return nil, ErrReplyTypeInvalid
	}
	return reply, nil
}

func (u *userServer) CreateRole(ctx context.Context, in *pb.CreateRoleRequest) (*pb.CreateRoleReply, error) {
	_, res, err := u.createRoleHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	reply, ok := res.(*pb.CreateRoleReply)
	if !ok {
		return nil, ErrReplyTypeInvalid
	}
	return reply, nil
}

func (u *userServer) AddRoleForUser(ctx context.Context, in *pb.AddRoleForUserRequest) (*pb.AddRoleForUserReply, error) {
	_, res, err := u.addRoleForUserHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	reply, ok := res.(*pb.AddRoleForUserReply)
	if !ok {
		return nil, ErrReplyTypeInvalid
	}
	return reply, nil
}

func (u *userServer) AddRoutes(ctx context.Context, in *pb.AddRoutesRequest) (*pb.AddRoutesReply, error) {
	_, res, err := u.addRoutesHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	reply, ok := res.(*pb.AddRoutesReply)
	if !ok {
		return nil, ErrReplyTypeInvalid
	}
	return reply, nil
}

func (u *userServer) ListRoutes(ctx context.Context, in *pb.ListRoutesRequest) (*pb.ListRoutesReply, error) {
	_, res, err := u.listRouteHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	reply, ok := res.(*pb.ListRoutesReply)
	if !ok {
		return nil, ErrReplyTypeInvalid
	}
	return reply, nil
}

func (u *userServer) CreatePermission(ctx context.Context, in *pb.CreatePermissionRequest) (*pb.CreatePermissionReply, error) {
	_, res, err := u.createPermissionHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	reply, ok := res.(*pb.CreatePermissionReply)
	if !ok {
		return nil, ErrReplyTypeInvalid
	}
	return reply, nil
}

func (u *userServer) UpdatePermission(ctx context.Context, in *pb.UpdatePermissionRequest) (*pb.UpdatePermissionReply, error) {
	_, res, err := u.updatePermissionHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	reply, ok := res.(*pb.UpdatePermissionReply)
	if !ok {
		return nil, ErrReplyTypeInvalid
	}
	return reply, nil
}

func (u *userServer) AddRouteForPermission(ctx context.Context, in *pb.AddRouteForPermissionRequest) (*pb.AddRouteForPermissionReply, error) {
	_, res, err := u.addRouteForPermissionHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	reply, ok := res.(*pb.AddRouteForPermissionReply)
	if !ok {
		return nil, ErrReplyTypeInvalid
	}
	return reply, nil
}

func (u *userServer) RemoveRouteForPermission(ctx context.Context, in *pb.RemoveRouteForPermissionRequest) (*pb.RemoveRouteForPermissionReply, error) {
	_, res, err := u.removeRouteForPermissionHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	reply, ok := res.(*pb.RemoveRouteForPermissionReply)
	if !ok {
		return nil, ErrReplyTypeInvalid
	}
	return reply, nil
}

func (u *userServer) RemovePermission(ctx context.Context, in *pb.RemovePermissionRequest) (*pb.RemovePermissionReply, error) {
	_, res, err := u.removePermissionHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	reply, ok := res.(*pb.RemovePermissionReply)
	if !ok {
		return nil, ErrReplyTypeInvalid
	}
	return reply, nil
}

func (u *userServer) ListPermissions(ctx context.Context, in *pb.ListPermissionsRequest) (*pb.ListPermissionsReply, error) {
	_, res, err := u.listPermissionsHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	reply, ok := res.(*pb.ListPermissionsReply)
	if !ok {
		return nil, ErrReplyTypeInvalid
	}
	return reply, nil
}

func (u *userServer) AddPermissionForPermission(ctx context.Context, in *pb.AddPermissionForPermissionRequest) (*pb.AddPermissionForPermissionReply, error) {
	_, res, err := u.addPermissionForPermissionHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	reply, ok := res.(*pb.AddPermissionForPermissionReply)
	if !ok {
		return nil, ErrReplyTypeInvalid
	}
	return reply, nil
}

func (u *userServer) RemovePermissionForPermission(ctx context.Context, in *pb.RemovePermissionForPermissionRequest) (*pb.RemovePermissionForPermissionReply, error) {
	_, res, err := u.removePermissionForPermissionHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	reply, ok := res.(*pb.RemovePermissionForPermissionReply)
	if !ok {
		return nil, ErrReplyTypeInvalid
	}
	return reply, nil
}

func (u *userServer) ListRole(ctx context.Context, in *pb.ListRoleRequest) (*pb.ListRoleReply, error) {
	_, res, err := u.listRoleHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	reply, ok := res.(*pb.ListRoleReply)
	if !ok {
		return nil, ErrReplyTypeInvalid
	}
	return reply, nil
}

func (u *userServer) UpdateRole(ctx context.Context, in *pb.UpdateRoleRequest) (*pb.UpdateRoleReply, error) {
	_, res, err := u.updateRoleHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	reply, ok := res.(*pb.UpdateRoleReply)
	if !ok {
		return nil, ErrReplyTypeInvalid
	}
	return reply, nil
}

func (u *userServer) RemovePermissionForRole(ctx context.Context, in *pb.RemovePermissionForRoleRequest) (*pb.RemovePermissionForRoleReply, error) {
	_, res, err := u.removePermissionForRoleHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	reply, ok := res.(*pb.RemovePermissionForRoleReply)
	if !ok {
		return nil, ErrReplyTypeInvalid
	}
	return reply, nil
}

func (u *userServer) AddRoleForRole(ctx context.Context, in *pb.AddRoleForRoleRequest) (*pb.AddRoleForRoleReply, error) {
	_, res, err := u.addRoleForRoleHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	reply, ok := res.(*pb.AddRoleForRoleReply)
	if !ok {
		return nil, ErrReplyTypeInvalid
	}
	return reply, nil
}

func (u *userServer) RemoveRoleForRole(ctx context.Context, in *pb.RemoveRoleForRoleRequest) (*pb.RemoveRoleForRoleReply, error) {
	_, res, err := u.removeRoleForRoleHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	reply, ok := res.(*pb.RemoveRoleForRoleReply)
	if !ok {
		return nil, ErrReplyTypeInvalid
	}
	return reply, nil
}

func (u *userServer) RemoveRole(ctx context.Context, in *pb.RemoveRoleRequest) (*pb.RemoveRoleReply, error) {
	_, res, err := u.removeRoleHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	reply, ok := res.(*pb.RemoveRoleReply)
	if !ok {
		return nil, ErrReplyTypeInvalid
	}
	return reply, nil
}

func (u *userServer) ListUsers(ctx context.Context, in *pb.ListUsersRequest) (*pb.ListUsersReply, error) {
	_, res, err := u.listUsersHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	reply, ok := res.(*pb.ListUsersReply)
	if !ok {
		return nil, ErrReplyTypeInvalid
	}
	return reply, nil
}

func (u *userServer) UpdateUser(ctx context.Context, in *pb.UpdateUserRequest) (*pb.UpdateUserReply, error) {
	_, res, err := u.updateUserHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	reply, ok := res.(*pb.UpdateUserReply)
	if !ok {
		return nil, ErrReplyTypeInvalid
	}
	return reply, nil
}

func (u *userServer) AddPermissionForUser(ctx context.Context, in *pb.AddPermissionForUserRequest) (*pb.AddPermissionForUserReply, error) {
	_, res, err := u.addPermissionForUserHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	reply, ok := res.(*pb.AddPermissionForUserReply)
	if !ok {
		return nil, ErrReplyTypeInvalid
	}
	return reply, nil
}

func (u *userServer) RemovePermissionForUser(ctx context.Context, in *pb.RemovePermissionForUserRequest) (*pb.RemovePermissionForUserReply, error) {
	_, res, err := u.removePermissionForUserHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	reply, ok := res.(*pb.RemovePermissionForUserReply)
	if !ok {
		return nil, ErrReplyTypeInvalid
	}
	return reply, nil
}

func (u *userServer) RemoveRoleForUser(ctx context.Context, in *pb.RemoveRoleForUserRequest) (*pb.RemoveRoleForUserReply, error) {
	_, res, err := u.removeRoleForUserHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	reply, ok := res.(*pb.RemoveRoleForUserReply)
	if !ok {
		return nil, ErrReplyTypeInvalid
	}
	return reply, nil
}

func (u *userServer) ListMenus(ctx context.Context, in *pb.ListMenusRequest) (*pb.ListMenusReply, error) {
	_, res, err := u.listMenusHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	reply, ok := res.(*pb.ListMenusReply)
	if !ok {
		return nil, ErrReplyTypeInvalid
	}
	return reply, nil
}

func (u *userServer) CreateMenu(ctx context.Context, in *pb.CreateMenuRequest) (*pb.CreateMenuReply, error) {
	_, res, err := u.createMenuHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	reply, ok := res.(*pb.CreateMenuReply)
	if !ok {
		return nil, ErrReplyTypeInvalid
	}
	return reply, nil
}

func (u *userServer) RemoveMenu(ctx context.Context, in *pb.RemoveMenuRequest) (*pb.RemoveMenuReply, error) {
	_, res, err := u.removeMenuHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	reply, ok := res.(*pb.RemoveMenuReply)
	if !ok {
		return nil, ErrReplyTypeInvalid
	}
	return reply, nil
}

func (u *userServer) GetUserTypeInfo(ctx context.Context, in *pb.GetUserTypeInfoRequest) (*pb.GetUserTypeInfoReply, error) {
	_, res, err := u.getUserTypeInfoHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	reply, ok := res.(*pb.GetUserTypeInfoReply)
	if !ok {
		return nil, ErrReplyTypeInvalid
	}
	return reply, nil
}
func (u *userServer) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.GetUserReply, error) {
	_, res, err := u.getUserHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	reply, ok := res.(*pb.GetUserReply)
	if !ok {
		return nil, ErrReplyTypeInvalid
	}
	return reply, nil

}
func (u *userServer) GetUserPermissionsAndRoles(ctx context.Context, in *pb.GetUserPermissionsAndRolesRequest) (*pb.GetUserPermissionsAndRolesReply, error) {
	_, res, err := u.getUserPermissionsAndRolesHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	reply, ok := res.(*pb.GetUserPermissionsAndRolesReply)
	if !ok {
		return nil, ErrReplyTypeInvalid
	}
	return reply, nil
}

func (u *userServer) GetRolePermissionsAndRoles(ctx context.Context, in *pb.GetRolePermissionsAndRolesRequest) (*pb.GetRolePermissionsAndRolesReply, error) {
	_, res, err := u.getRolePermissionsAndRolesHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	reply, ok := res.(*pb.GetRolePermissionsAndRolesReply)
	if !ok {
		return nil, ErrReplyTypeInvalid
	}
	return reply, nil
}

func (u *userServer) GetPermissionsAndRoutes(ctx context.Context, in *pb.GetPermissionsAndRoutesRequest) (*pb.GetPermissionsAndRoutesReply, error) {
	_, res, err := u.getPermissionsAndRoutesHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	reply, ok := res.(*pb.GetPermissionsAndRoutesReply)
	if !ok {
		return nil, ErrReplyTypeInvalid
	}
	return reply, nil
}

func (u *userServer) ListLeaguer(ctx context.Context, in *pb.ListLeaguerRequest) (*pb.ListLeaguerReply, error) {
	_, res, err := u.listLeaguerHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	reply, ok := res.(*pb.ListLeaguerReply)
	if !ok {
		return nil, ErrReplyTypeInvalid
	}
	return reply, nil
}
