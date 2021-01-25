package transport

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/baxiang/soldiers-sortie/go-mircosvc/pkg/common"
	sharedEtcd "github.com/baxiang/soldiers-sortie/go-mircosvc/pkg/etcd"
	"github.com/baxiang/soldiers-sortie/go-mircosvc/pkg/middleware"
	"github.com/baxiang/soldiers-sortie/go-mircosvc/pkg/session"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/etcdv3"
	"github.com/go-kit/kit/sd/lb"
	kitZipkin "github.com/go-kit/kit/tracing/zipkin"
	kitTransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"



	usersvcEndpoints "github.com/baxiang/soldiers-sortie/go-mircosvc/servers/usersvc/endpoints"
	usersvcTransport "github.com/baxiang/soldiers-sortie/go-mircosvc/servers/usersvc/transport"

	articlesvcEndpoints "github.com/baxiang/soldiers-sortie/go-mircosvc/servers/article/endpoints"
	articlesvcTransport "github.com/baxiang/soldiers-sortie/go-mircosvc/servers/article/transport"
	articlesvcServices "github.com/baxiang/soldiers-sortie/go-mircosvc/servers/article/services"
)

func MakeHandler(
	etcdClient etcdv3.Client,
	tracer opentracing.Tracer,
	zipkinTracer *zipkin.Tracer,
	logger log.Logger,
	retryMax int,
	retryTimeout int,
	sessionStorage session.Storager,
) http.Handler {

	opts := []kitTransport.ServerOption{
		kitTransport.ServerBefore(cookieToContext()),
		kitZipkin.HTTPServerTrace(zipkinTracer),
		kitTransport.ServerErrorEncoder(common.EncodeError),
	}

	r := mux.NewRouter()
	// user endpoint
	{
		endpoints := new(usersvcEndpoints.Endponits)
		ins := sharedEtcd.NewInstancer("/usersvc", etcdClient, logger)

		{
			factory := usersvcFactory(usersvcEndpoints.MakeGetUserEndpoint, tracer, zipkinTracer, logger)
			endpoints.GetUserEP = makeEndpoint(factory, ins, logger, retryMax, retryTimeout)
		}
		{
			factory := usersvcFactory(usersvcEndpoints.MakeLoginEndpoint, tracer, zipkinTracer, logger)
			endpoints.LoginEP = makeEndpoint(factory, ins, logger, retryMax, retryTimeout)
		}
		{
			factory := usersvcFactory(usersvcEndpoints.MakeRegisterEndpoint, tracer, zipkinTracer, logger)
			endpoints.RegisterEP = makeEndpoint(factory, ins, logger, retryMax, retryTimeout)
		}
		{
			factory := usersvcFactory(usersvcEndpoints.MakeSendCodeEndpoint, tracer, zipkinTracer, logger)
			endpoints.SendCodeEP = makeEndpoint(factory, ins, logger, retryMax, retryTimeout)
		}
		{
			factory := usersvcFactory(usersvcEndpoints.MakeUserListEndpoint, tracer, zipkinTracer, logger)
			endpoints.UserListEP = makeEndpoint(factory, ins, logger, retryMax, retryTimeout,
				middleware.CookieMiddleware(sessionStorage),
				middleware.PermissionMiddleware(common.RootUser),
			)
		}
		{
			factory := usersvcFactory(usersvcEndpoints.MakeAuthEndpoint, tracer, zipkinTracer, logger)
			endpoints.AuthEP = makeEndpoint(factory, ins, logger, retryMax, retryTimeout,
				middleware.CookieMiddleware(sessionStorage),
			)
		}
		{
			factory := usersvcFactory(usersvcEndpoints.MakeLogoutEndpoint, tracer, zipkinTracer, logger)
			endpoints.LogoutEP = makeEndpoint(factory, ins, logger, retryMax, retryTimeout,
				middleware.CookieMiddleware(sessionStorage),
			)
		}

		r.PathPrefix("/usersvc").Handler(http.StripPrefix("/usersvc", usersvcTransport.MakeHTTPHandler(endpoints, tracer,
			logger, opts)))
	}

	// article endpoint
	{
		endpoints := new(articlesvcEndpoints.Endpoints)
		ins := sharedEtcd.NewInstancer("/articlesvc", etcdClient, logger)

		{
			factory := articlesvcFactory(articlesvcEndpoints.MakeGetCategoriesEndpoint, tracer, zipkinTracer, logger)
			endpoints.GetCategoriesEP = makeEndpoint(factory, ins, logger, retryMax, retryTimeout)

			r.PathPrefix("/articlesvc").Handler(http.StripPrefix("/articlesvc", articlesvcTransport.MakeHTTPHandler(endpoints, tracer,
				logger, opts)))
		}
	}

	return r
}

func usersvcFactory(
	makeEndpoint func(service usersvcEndpoints.UserSerivcer) endpoint.Endpoint,
	tracer opentracing.Tracer,
	zipkinTracer *zipkin.Tracer,
	logger log.Logger,
) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		conn, err := grpc.Dial(instance, grpc.WithInsecure())
		if err != nil {
			return nil, nil, err
		}
		service := usersvcTransport.MakeGRPCClient(conn, tracer, zipkinTracer, logger)
		return makeEndpoint(service), conn, nil
	}
}

func articlesvcFactory(
	makeEndpoint func(service articlesvcServices.ArticleServicer) endpoint.Endpoint,
	tracer opentracing.Tracer,
	zipkinTracer *zipkin.Tracer,
	logger log.Logger,
) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		conn, err := grpc.Dial(instance, grpc.WithInsecure())
		if err != nil {
			return nil, nil, err
		}
		service := articlesvcTransport.MakeGRPCClient(conn, tracer, zipkinTracer, logger)
		return makeEndpoint(service), conn, nil
	}
}

func makeEndpoint(
	factory sd.Factory,
	ins *etcdv3.Instancer,
	logger log.Logger,
	retryMax int,
	retryTimeout int,
	middlewares ...endpoint.Middleware,
) endpoint.Endpoint {
	endpointer := sd.NewEndpointer(ins, factory, logger)
	balancer := lb.NewRoundRobin(endpointer)

	ep := lb.RetryWithCallback(time.Duration(retryTimeout)*time.Second, balancer, func(n int, received error) (bool,
		error) {
		if st, ok := status.FromError(received); ok {
			return false, errors.New(st.Message())
		}
		return n < retryMax, nil
	})

	for _, m := range middlewares {
		ep = m(ep)
	}

	return ep
}

func cookieToContext() kitTransport.RequestFunc {
	return func(ctx context.Context, r *http.Request) context.Context {
		c, err := r.Cookie(common.CookieName)
		if err != nil {
			return ctx
		}

		return context.WithValue(ctx, common.SessionKey, c.Value)
	}
}