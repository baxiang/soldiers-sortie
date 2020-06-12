package middleware

import (
	"context"
	"github.com/baxiang/koala/errno"
	"github.com/baxiang/koala/logs"
	"github.com/baxiang/koala/meta"
	"google.golang.org/grpc/balancer"
)

func NewLoadBalanceMiddleware(balancer balancer.Balancer) Middleware{
	return func(next MiddlewareFunc) MiddlewareFunc {
		return func(ctx context.Context, req interface{}) (rsp interface{}, err error) {
			rpcMeta := meta.GetRpcMeta(ctx)
			if len(rpcMeta.AllNodes) == 0 {
				err = errno.NotHaveInstance
				logs.Error(ctx, "not have instance")
				return
			}
			ctx = loadbalance.WithBalanceContext(ctx)
		}
	}
}