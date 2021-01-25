package middleware

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"

	"github.com/baxiang/soldiers-sortie/go-mircosvc/pkg/common"

)

func PermissionMiddleware(level common.RoleLevel) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			m, ok := ctx.Value(common.CookieName).(map[string]interface{})
			if ok {
				if l, ok := m[common.RoleIDKey]; ok {
					if ll, ok := l.(common.RoleLevel); ok && ll <= level {
						return next(ctx, request)
					}
				}

				return nil, common.NewError(http.StatusForbidden, "权限不足.")
			}

			return nil, common.NewError(http.StatusUnauthorized, "请登录.")
		}
	}
}

