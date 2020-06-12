package middleware

import (
	"context"
	"github.com/baxiang/koala/meta"
	"github.com/afex/hystrix-go/hystrix"
)

func HystrixMiddleware(next MiddlewareFunc) MiddlewareFunc {
	return func(ctx context.Context, req interface{}) (rsp interface{}, err error) {
		rpcMeta := meta.GetRpcMeta(ctx)
		var resp interface{}

		hystrixErr := hystrix.Do(rpcMeta.ServiceName, func() (err error) {
			resp, err = next(ctx, req)
			return err
		}, nil)

		if hystrixErr != nil {
			return nil, hystrixErr
		}

		return resp, hystrixErr
	}
}
