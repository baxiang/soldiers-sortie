package transport

import (
	"github.com/baxiang/soldiers-sortie/go-mircosvc/servers/article/endpoints"
	"github.com/go-kit/kit/log"
	"github.com/baxiang/soldiers-sortie/go-mircosvc/pkg/common"
	kitOpentracing "github.com/go-kit/kit/tracing/opentracing"
	kitTransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func MakeHTTPHandler(eps *endpoints.Endpoints,
	otTracer opentracing.Tracer,
	logger log.Logger,
	opts []kitTransport.ServerOption,
) http.Handler {
	m := mux.NewRouter()
	m.Handle("/metrics", promhttp.Handler())

	{
		handler := kitTransport.NewServer(
			eps.GetCategoriesEP,
			common.DecodeEmptyHttpRequest,
			kitTransport.EncodeJSONResponse,
			append(opts, kitTransport.ServerBefore(kitOpentracing.HTTPToContext(otTracer, "GetCategories", logger)))...,
		)
		m.Handle("/category", handler).Methods("GET")
	}

	return m
}