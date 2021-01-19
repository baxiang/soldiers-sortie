package userservice

import (
	"context"
	"userService/pkg/pb"

	"github.com/go-kit/kit/endpoint"
)

func MakeLoginEndpoint(service pb.UserServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.LoginRequest)
		if !ok {
			return nil, ErrRequestTypeInvalid
		}
		return service.Login(ctx, req)
	}
}

func MakeGetPermissionsEndpoint(service pb.UserServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.GetPermissionsRequest)
		if !ok {
			return nil, ErrRequestTypeInvalid
		}
		return service.GetPermissions(ctx, req)
	}
}

func MakeCheckPermissionEndpoint(service pb.UserServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.CheckPermissionRequest)
		if !ok {
			return nil, ErrRequestTypeInvalid
		}
		return service.CheckPermission(ctx, req)
	}
}

func MakeRegisterEndpoint(service pb.UserServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.RegisterRequest)
		if !ok {
			return nil, ErrRequestTypeInvalid
		}
		return service.Register(ctx, req)
	}
}

func MakeAddPermissionForRoleEndpoint(service pb.UserServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.AddPermissionForRoleRequest)
		if !ok {
			return nil, ErrRequestTypeInvalid
		}
		return service.AddPermissionForRole(ctx, req)
	}
}

func MakeAddRoleForRoleEndpoint(service pb.UserServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.AddRoleForRoleRequest)
		if !ok {
			return nil, ErrRequestTypeInvalid
		}
		return service.AddRoleForRole(ctx, req)
	}
}

func MakeCreateRoleEndpoint(service pb.UserServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.CreateRoleRequest)
		if !ok {
			return nil, ErrRequestTypeInvalid
		}
		return service.CreateRole(ctx, req)
	}
}

func MakeAddRoleForUserEndpoint(service pb.UserServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.AddRoleForUserRequest)
		if !ok {
			return nil, ErrRequestTypeInvalid
		}
		return service.AddRoleForUser(ctx, req)
	}
}

func MakeAddRoutesEndpoint(service pb.UserServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.AddRoutesRequest)
		if !ok {
			return nil, ErrRequestTypeInvalid
		}
		return service.AddRoutes(ctx, req)
	}
}

func MakeListRoutesEndpoint(service pb.UserServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.ListRoutesRequest)
		if !ok {
			return nil, ErrRequestTypeInvalid
		}
		return service.ListRoutes(ctx, req)
	}
}

func MakeCreatePermissionEndpoint(service pb.UserServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.CreatePermissionRequest)
		if !ok {
			return nil, ErrRequestTypeInvalid
		}
		return service.CreatePermission(ctx, req)
	}
}

func MakeUpdatePermissionEndpoint(service pb.UserServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.UpdatePermissionRequest)
		if !ok {
			return nil, ErrRequestTypeInvalid
		}
		return service.UpdatePermission(ctx, req)
	}
}

func MakeAddRouteForPermissionEndpoint(service pb.UserServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.AddRouteForPermissionRequest)
		if !ok {
			return nil, ErrRequestTypeInvalid
		}
		return service.AddRouteForPermission(ctx, req)
	}
}

func MakeRemoveRouteForPermissionEndpoint(service pb.UserServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.RemoveRouteForPermissionRequest)
		if !ok {
			return nil, ErrRequestTypeInvalid
		}
		return service.RemoveRouteForPermission(ctx, req)
	}
}

func MakeRemovePermissionEndpoint(service pb.UserServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.RemovePermissionRequest)
		if !ok {
			return nil, ErrRequestTypeInvalid
		}
		return service.RemovePermission(ctx, req)
	}
}

func MakeListPermissionsEndpoint(service pb.UserServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.ListPermissionsRequest)
		if !ok {
			return nil, ErrRequestTypeInvalid
		}
		return service.ListPermissions(ctx, req)
	}
}

func MakeAddPermissionForPermissionEndpoint(service pb.UserServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.AddPermissionForPermissionRequest)
		if !ok {
			return nil, ErrRequestTypeInvalid
		}
		return service.AddPermissionForPermission(ctx, req)
	}
}

func MakeRemovePermissionForPermissionEndpoint(service pb.UserServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.RemovePermissionForPermissionRequest)
		if !ok {
			return nil, ErrRequestTypeInvalid
		}
		return service.RemovePermissionForPermission(ctx, req)
	}
}

func MakeListRoleEndpoint(service pb.UserServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.ListRoleRequest)
		if !ok {
			return nil, ErrRequestTypeInvalid
		}
		return service.ListRole(ctx, req)
	}
}

func MakeUpdateRoleEndpoint(service pb.UserServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.UpdateRoleRequest)
		if !ok {
			return nil, ErrRequestTypeInvalid
		}
		return service.UpdateRole(ctx, req)
	}
}

func MakeRemovePermissionForRoleEndpoint(service pb.UserServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.RemovePermissionForRoleRequest)
		if !ok {
			return nil, ErrRequestTypeInvalid
		}
		return service.RemovePermissionForRole(ctx, req)
	}
}

func MakeRemoveRoleForRoleEndpoint(service pb.UserServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.RemoveRoleForRoleRequest)
		if !ok {
			return nil, ErrRequestTypeInvalid
		}
		return service.RemoveRoleForRole(ctx, req)
	}
}

func MakeRemoveRoleEndpoint(service pb.UserServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.RemoveRoleRequest)
		if !ok {
			return nil, ErrRequestTypeInvalid
		}
		return service.RemoveRole(ctx, req)
	}
}

func MakeListUsersEndpoint(service pb.UserServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.ListUsersRequest)
		if !ok {
			return nil, ErrRequestTypeInvalid
		}
		return service.ListUsers(ctx, req)
	}
}

func MakeUpdateUserEndpoint(service pb.UserServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.UpdateUserRequest)
		if !ok {
			return nil, ErrRequestTypeInvalid
		}
		return service.UpdateUser(ctx, req)
	}
}

func MakeAddPermissionForUserEndpoint(service pb.UserServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.AddPermissionForUserRequest)
		if !ok {
			return nil, ErrRequestTypeInvalid
		}
		return service.AddPermissionForUser(ctx, req)
	}
}

func MakeRemovePermissionForUserEndpoint(service pb.UserServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.RemovePermissionForUserRequest)
		if !ok {
			return nil, ErrRequestTypeInvalid
		}
		return service.RemovePermissionForUser(ctx, req)
	}
}

func MakeRemoveRoleForUserEndpoint(service pb.UserServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.RemoveRoleForUserRequest)
		if !ok {
			return nil, ErrRequestTypeInvalid
		}
		return service.RemoveRoleForUser(ctx, req)
	}
}

func MakeListMenusEndpoint(service pb.UserServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.ListMenusRequest)
		if !ok {
			return nil, ErrRequestTypeInvalid
		}
		return service.ListMenus(ctx, req)
	}
}

func MakeCreateMenuEndpoint(service pb.UserServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.CreateMenuRequest)
		if !ok {
			return nil, ErrRequestTypeInvalid
		}
		return service.CreateMenu(ctx, req)
	}
}

func MakeRemoveMenuEndpoint(service pb.UserServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.RemoveMenuRequest)
		if !ok {
			return nil, ErrRequestTypeInvalid
		}
		return service.RemoveMenu(ctx, req)
	}
}

func MakeGetUserTypeInfoEndpoint(service pb.UserServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.GetUserTypeInfoRequest)
		if !ok {
			return nil, ErrRequestTypeInvalid
		}
		return service.GetUserTypeInfo(ctx, req)
	}
}

func MakeGetUserEndpoint(service pb.UserServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.GetUserRequest)
		if !ok {
			return nil, ErrRequestTypeInvalid
		}
		return service.GetUser(ctx, req)
	}
}

func MakeGetUserPermissionsAndRolesEndpoint(service pb.UserServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.GetUserPermissionsAndRolesRequest)
		if !ok {
			return nil, ErrRequestTypeInvalid
		}
		return service.GetUserPermissionsAndRoles(ctx, req)
	}
}

func MakeGetRolePermissionsAndRolesEndpoint(service pb.UserServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.GetRolePermissionsAndRolesRequest)
		if !ok {
			return nil, ErrRequestTypeInvalid
		}
		return service.GetRolePermissionsAndRoles(ctx, req)
	}
}

func MakeGetPermissionsAndRoutesEndpoint(service pb.UserServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.GetPermissionsAndRoutesRequest)
		if !ok {
			return nil, ErrRequestTypeInvalid
		}
		return service.GetPermissionsAndRoutes(ctx, req)
	}
}

func MakeListLeaguerEndpoint(service pb.UserServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.ListLeaguerRequest)
		if !ok {
			return nil, ErrRequestTypeInvalid
		}
		return service.ListLeaguer(ctx, req)
	}
}

func MakeRemoveRouteEndpoint(service pb.UserServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return service.RemoveRoute(ctx, request.(*pb.RemoveRouteRequest))
	}
}

func decodeRequest(ctx context.Context, request interface{}) (interface{}, error) {
	return request, nil
}

func encodeResponse(ctx context.Context, response interface{}) (interface{}, error) {
	return response, nil
}
