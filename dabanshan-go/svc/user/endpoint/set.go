package endpoint

import (
	"context"
	"github.com/baxiang/soldiers-sortie/dabanshan-go/svc/user/model"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/sony/gobreaker"

	//"github.com/go-kit/kit/ratelimit"
	//rl "github.com/juju/ratelimit"

	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/baxiang/soldiers-sortie/dabanshan-go/svc/user/service"
)

type Set struct {
	GetUserEndpoint  endpoint.Endpoint
	RegisterEndpoint endpoint.Endpoint
	LoginEndpoint    endpoint.Endpoint
}

func New(svc service.Service, logger log.Logger, duration metrics.Histogram,
	trace stdopentracing.Tracer)Set{
	var (
		getUserEndpoint  endpoint.Endpoint
		registerEndpoint endpoint.Endpoint
		loginEndpoint    endpoint.Endpoint
	)

	{
		getUserEndpoint = MakeGetUserEndpoint(svc)
		//getUserEndpoint = ratelimit.NewTokenBucketLimiter(rl.NewBucketWithRate(1, 1))(getUserEndpoint)
		getUserEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(getUserEndpoint)
		getUserEndpoint = opentracing.TraceServer(trace, "GetUser")(getUserEndpoint)
		getUserEndpoint = LoggingMiddleware(log.With(logger, "method", "GetUser"))(getUserEndpoint)
		getUserEndpoint = InstrumentingMiddleware(duration.With("method", "GetUser"))(getUserEndpoint)
	}

	{
		registerEndpoint = MakeRegisterEndpoint(svc)
		//registerEndpoint = ratelimit.NewTokenBucketLimiter(rl.NewBucketWithRate(1, 1))(registerEndpoint)
		registerEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(registerEndpoint)
		registerEndpoint = opentracing.TraceServer(trace, "Register")(registerEndpoint)
		registerEndpoint = LoggingMiddleware(log.With(logger, "method", "Register"))(registerEndpoint)
		registerEndpoint = InstrumentingMiddleware(duration.With("method", "Register"))(registerEndpoint)
	}

	{
		loginEndpoint = MakeLoginEndpoint(svc)
		//loginEndpoint = ratelimit.NewTokenBucketLimiter(rl.NewBucketWithRate(1, 1))(loginEndpoint)
		loginEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(loginEndpoint)
		loginEndpoint = opentracing.TraceServer(trace, "Login")(loginEndpoint)
		loginEndpoint = LoggingMiddleware(log.With(logger, "method", "Login"))(loginEndpoint)
		loginEndpoint = InstrumentingMiddleware(duration.With("method", "Login"))(loginEndpoint)
	}
	return Set{
		GetUserEndpoint:  getUserEndpoint,
		RegisterEndpoint: registerEndpoint,
		LoginEndpoint:    loginEndpoint,
	}
}

func MakeGetUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(model.GetUserRequest)
		v, err := s.GetUser(ctx, req.A)
		return v, err
	}
}

// MakeRegisterEndpoint constructs a Register endpoint wrapping the service.
func MakeRegisterEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(model.RegisterRequest)
		v, err := s.Register(ctx, req)
		return v, err
	}
}

// MakeLoginEndpoint constructs a Login endpoint wrapping the service.
func MakeLoginEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(model.LoginRequest)
		v, err := s.Login(ctx, req)
		return v, err
	}
}


// GetUser implements the service interface, so Set may be used as a service.
func (s Set) GetUser(ctx context.Context, a string) (model.GetUserResponse, error) {
	resp, err := s.GetUserEndpoint(ctx, model.GetUserRequest{A: a})
	if err != nil {
		return model.GetUserResponse{}, err
	}
	response := resp.(model.GetUserResponse)
	return response, response.Err
}

// Register implements the service interface,
func (s Set) Register(ctx context.Context, us model.RegisterRequest) (r model.RegisterUserResponse, err error) {
	resp, err := s.RegisterEndpoint(ctx, us)
	if err != nil {
		return model.RegisterUserResponse{ID: ""}, err
	}
	response := resp.(model.RegisterUserResponse)
	return response, err
}

// Login implements the service interface.
func (s Set) Login(ctx context.Context, login model.LoginRequest) (model.LoginResponse, error) {
	resp, err := s.LoginEndpoint(ctx, login)
	if err != nil {
		return model.LoginResponse{}, err
	}
	response := resp.(model.LoginResponse)
	return response, err
}