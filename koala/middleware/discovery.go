package middleware

import (
	"github.com/baxiang/koala/logs"
	"github.com/baxiang/koala/meta"
	"github.com/baxiang/koala/registry"
	"context"
)

func NewDiscoveryMiddleware(discovery registry.Registry)Middleware{
	return func(next MiddlewareFunc) MiddlewareFunc {
		return func(ctx context.Context, req interface{}) (rsp interface{}, err error) {
			rpcMeta := meta.GetRpcMeta(ctx)
			if len(rpcMeta.AllNodes) > 0 {
				return next(ctx, req)
			}

			service, err := discovery.GetService(ctx, rpcMeta.ServiceName)
			if err != nil {
				logs.Error(ctx, "discovery service:%s failed, err:%v", rpcMeta.ServiceName, err)
				return
			}

			rpcMeta.AllNodes = service.Nodes
			rsp, err = next(ctx, req)
			return
		}
	}
}