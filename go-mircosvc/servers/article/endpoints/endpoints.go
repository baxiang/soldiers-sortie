package endpoints

import (
	"context"
	"github.com/baxiang/soldiers-sortie/go-mircosvc/servers/article/services"
	"time"


	"github.com/go-kit/kit/log"
	articlePb "github.com/baxiang/soldiers-sortie/go-mircosvc/pb"
	"github.com/baxiang/soldiers-sortie/go-mircosvc/pkg/common"
	"github.com/baxiang/soldiers-sortie/go-mircosvc/pkg/middleware"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	kitOpentracing "github.com/go-kit/kit/tracing/opentracing"
	kitZipkin "github.com/go-kit/kit/tracing/zipkin"
	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"

)

type Endpoints struct {
	GetCategoriesEP endpoint.Endpoint
}

func (e *Endpoints) GetCategories(ctx context.Context) (res *articlePb.GetCategoriesResponse, err error) {
	r, err := e.GetCategoriesEP(ctx, nil)

	if r != nil {
		res = r.(*articlePb.GetCategoriesResponse)
	}

	return
}

func NewEndpoints(svc services.ArticleServicer, logger log.Logger, otTracer opentracing.Tracer, zipkinTracer *zipkin.Tracer) *Endpoints {

	return &Endpoints{
		GetCategoriesEP: makeEndpoint(MakeGetCategoriesEndpoint(svc), "GetCategories", logger, otTracer, zipkinTracer),
	}
}

// GetCategories
func MakeGetCategoriesEndpoint(svc services.ArticleServicer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := svc.GetCategories(ctx)

		return common.Response{Data: res}, err
	}
}

func makeEndpoint(ep endpoint.Endpoint, method string,
	logger log.Logger, otTracer opentracing.Tracer,
	zipkinTracer *zipkin.Tracer,
	middlewares ...endpoint.Middleware,
) endpoint.Endpoint {
	limiter := rate.NewLimiter(rate.Every(time.Second*1), 10)

	middlewares = append(
		middlewares,
		middleware.RateLimitterMiddleware(limiter),
		circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{})),
		kitOpentracing.TraceServer(otTracer, method),
		kitZipkin.TraceEndpoint(zipkinTracer, method),
		middleware.LoggingMiddleware(log.With(logger, "method", method)),
	)

	for _, m := range middlewares {
		ep = m(ep)
	}

	return ep
}