package jaeger

import (
	"context"
	"github.com/baxiang/gorpc/interceptor"
	"github.com/baxiang/gorpc/plugin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"strings"
)

type Jaeger struct {
	opts *plugin.Options
}

const Name  = "jaeger"
const JaegerClientName = "gorpc-client-jaeger"
const JaegerServerName = "gorpc-server-jaeger"

func init(){
	plugin.Register(Name,JaegerSvr)
}

var JaegerSvr = &Jaeger{opts:&plugin.Options{}}

type jaegerCarrier map[string][]byte

func (m jaegerCarrier) Set(key, val string) {
	key = strings.ToLower(key)
	m[key] = []byte(val)
}

func (m jaegerCarrier) ForeachKey(handler func(key, val string) error) error {
	for k, v := range m {
		handler(k, string(v))
	}
	return nil
}

// OpenTracingClientInterceptor packaging jaeger tracer as a client interceptor
func OpenTracingClientInterceptor(tracer opentracing.Tracer, spanName string) interceptor.ClientInterceptor {
	return func (ctx context.Context, req, rsp interface{}, ivk interceptor.Invoker) error {

		//clientSpan := tracer.StartSpan(spanName, ext.SpanKindRPCClient, opentracing.ChildOf(parentCtx))
		clientSpan := tracer.StartSpan(spanName, ext.SpanKindRPCClient)
		defer clientSpan.Finish()

		mdCarrier := &jaegerCarrier{}

		if err := tracer.Inject(clientSpan.Context(), opentracing.HTTPHeaders, mdCarrier); err != nil {
			clientSpan.LogFields(log.String("event", "Tracer.Inject() failed"), log.Error(err))
		}

		clientSpan.LogFields(log.String("spanName", spanName))

		return ivk(ctx, req, rsp)

	}
}
