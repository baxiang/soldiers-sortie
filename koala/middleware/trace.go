package middleware

import (
	"context"
	"github.com/baxiang/koala/logs"
	"github.com/baxiang/koala/meta"
	"github.com/baxiang/koala/util"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc/metadata"
)

func TraceServerMiddleware(next MiddlewareFunc)MiddlewareFunc  {
	return func(ctx context.Context, req interface{}) (rsp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok{
			md =metadata.Pairs()
		}
		tracer := opentracing.GlobalTracer()
		parentSpanContext, err := tracer.Extract(opentracing.HTTPHeaders, metadataTextMap(md))
		if err != nil && err != opentracing.ErrSpanContextNotFound {
			logs.Warn(ctx, "trace extract failed, parsing trace information: %v", err)
		}

		serverMeta := meta.GetServerMeta(ctx)
		//开始追踪该方法
		serverSpan := tracer.StartSpan(
			serverMeta.Method,
			ext.RPCServerOption(parentSpanContext),
			ext.SpanKindRPCServer,
		)

		serverSpan.SetTag(util.TraceID, logs.GetTraceId(ctx))
		ctx = opentracing.ContextWithSpan(ctx, serverSpan)
		rsp, err = next(ctx, req)
		//记录错误
		if err != nil {
			ext.Error.Set(serverSpan, true)
			serverSpan.LogFields(log.String("event", "error"), log.String("message", err.Error()))
		}

		serverSpan.Finish()
		return
	}
}