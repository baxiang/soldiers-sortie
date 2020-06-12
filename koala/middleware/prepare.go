package middleware

import (
	"context"
	"github.com/baxiang/koala/logs"
	"google.golang.org/grpc/metadata"
)

func PrepareMiddleware(next MiddlewareFunc)MiddlewareFunc{
	return func(ctx context.Context, req interface{}) (rsp interface{}, err error) {
		//处理traceId
		var traceId string
		//从ctx获取grpc的metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			vals, ok := md[util.TraceID]
			if ok && len(vals) > 0 {
				traceId = vals[0]
			}
		}

		if len(traceId) == 0 {
			traceId = logs.GenTraceId()
		}

		ctx = logs.WithFieldContext(ctx)
		ctx = logs.WithTraceId(ctx, traceId)
		rsp, err = next(ctx, req)
		return
	}
}