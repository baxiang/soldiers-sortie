package gateway

import (
	"io"
	"os"
	"time"
	"userService/pkg/pb"
	"userService/pkg/userservice"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	"github.com/go-kit/kit/tracing/zipkin"
	stdzipkin "github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func GetUserEndpoints(instancer sd.Instancer, log log.Logger) *UserEndpoints {
	var endpoints UserEndpoints

	hystrix.ConfigureCommand(userBreaker, hystrix.CommandConfig{
		MaxConcurrentRequests: 1000,
		Timeout:               10000,
		ErrorPercentThreshold: 25,
		SleepWindow:           10000,
	})
	userBreaker := circuitbreaker.Hystrix(userBreaker)

	{
		factory := userserviceFactory(userservice.MakeLoginEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(rpcRetryTimes, rpcTimeOut, balancer)
		retry = userBreaker(retry)
		endpoints.LoginEndpoint = retry
	}

	{
		factory := userserviceFactory(userservice.MakeGetPermissionsEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(3, 500*time.Millisecond, balancer)
		retry = userBreaker(retry)
		endpoints.GetPermissionsEndpoint = retry
	}

	{
		factory := userserviceFactory(userservice.MakeCheckPermissionEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(3, 500*time.Millisecond, balancer)
		retry = userBreaker(retry)
		endpoints.CheckPermissionEndpoint = retry
	}

	{
		factory := userserviceFactory(userservice.MakeRegisterEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(3, 500*time.Millisecond, balancer)
		retry = userBreaker(retry)
		endpoints.RegisterEndpoint = retry
	}

	{
		factory := userserviceFactory(userservice.MakeAddPermissionForRoleEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(3, 500*time.Millisecond, balancer)
		retry = userBreaker(retry)
		endpoints.AddPermissionForRoleEndpoint = retry
	}

	{
		factory := userserviceFactory(userservice.MakeCreateRoleEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(3, 500*time.Millisecond, balancer)
		retry = userBreaker(retry)
		endpoints.CreateRoleEndpoint = retry
	}

	{
		factory := userserviceFactory(userservice.MakeAddRoleForUserEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(3, 500*time.Millisecond, balancer)
		retry = userBreaker(retry)
		endpoints.AddRoleForUserEndpoint = retry
	}

	{
		factory := userserviceFactory(userservice.MakeAddRoutesEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(3, 500*time.Millisecond, balancer)
		retry = userBreaker(retry)
		endpoints.AddRoutesEndpoint = retry
	}

	{
		factory := userserviceFactory(userservice.MakeListRoutesEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(3, 500*time.Millisecond, balancer)
		retry = userBreaker(retry)
		endpoints.ListRoutesEndpoint = retry
	}

	{
		factory := userserviceFactory(userservice.MakeCreatePermissionEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(3, 500*time.Millisecond, balancer)
		retry = userBreaker(retry)
		endpoints.CreatePermissionEndpoint = retry
	}

	{
		factory := userserviceFactory(userservice.MakeUpdatePermissionEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(3, 500*time.Millisecond, balancer)
		retry = userBreaker(retry)
		endpoints.UpdatePermissionEndpoint = retry
	}

	{
		factory := userserviceFactory(userservice.MakeAddRouteForPermissionEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(3, 500*time.Millisecond, balancer)
		retry = userBreaker(retry)
		endpoints.AddRouteForPermissionEndpoint = retry
	}

	{
		factory := userserviceFactory(userservice.MakeRemoveRouteForPermissionEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(3, 500*time.Millisecond, balancer)
		retry = userBreaker(retry)
		endpoints.RemoveRouteForPermissionEndpoint = retry
	}

	{
		factory := userserviceFactory(userservice.MakeRemovePermissionEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(3, 500*time.Millisecond, balancer)
		retry = userBreaker(retry)
		endpoints.RemovePermissionEndpoint = retry
	}

	{
		factory := userserviceFactory(userservice.MakeListPermissionsEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(3, 500*time.Millisecond, balancer)
		retry = userBreaker(retry)
		endpoints.ListPermissionsEndpoint = retry
	}

	{
		factory := userserviceFactory(userservice.MakeAddPermissionForPermissionEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(3, 500*time.Millisecond, balancer)
		retry = userBreaker(retry)
		endpoints.AddPermissionForPermissionEndpoint = retry
	}

	{
		factory := userserviceFactory(userservice.MakeRemovePermissionForPermissionEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(3, 500*time.Millisecond, balancer)
		retry = userBreaker(retry)
		endpoints.RemovePermissionForPermissionEndpoint = retry
	}

	{
		factory := userserviceFactory(userservice.MakeListRoleEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(3, 500*time.Millisecond, balancer)
		retry = userBreaker(retry)
		endpoints.ListRoleEndpoint = retry
	}

	{
		factory := userserviceFactory(userservice.MakeUpdateRoleEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(3, 500*time.Millisecond, balancer)
		retry = userBreaker(retry)
		endpoints.UpdateRoleEndpoint = retry
	}

	{
		factory := userserviceFactory(userservice.MakeRemovePermissionForRoleEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(3, 500*time.Millisecond, balancer)
		retry = userBreaker(retry)
		endpoints.RemovePermissionForRoleEndpoint = retry
	}

	{
		factory := userserviceFactory(userservice.MakeAddRoleForRoleEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(3, 500*time.Millisecond, balancer)
		retry = userBreaker(retry)
		endpoints.AddRoleForRoleEndpoint = retry
	}

	{
		factory := userserviceFactory(userservice.MakeRemoveRoleForRoleEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(3, 500*time.Millisecond, balancer)
		retry = userBreaker(retry)
		endpoints.RemoveRoleForRoleEndpoint = retry
	}

	{
		factory := userserviceFactory(userservice.MakeRemoveRoleEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(3, 500*time.Millisecond, balancer)
		retry = userBreaker(retry)
		endpoints.RemoveRoleEndpoint = retry
	}

	{
		factory := userserviceFactory(userservice.MakeListUsersEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(3, 500*time.Millisecond, balancer)
		retry = userBreaker(retry)
		endpoints.ListUsersEndpoint = retry
	}

	{
		factory := userserviceFactory(userservice.MakeUpdateUserEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(3, 500*time.Millisecond, balancer)
		retry = userBreaker(retry)
		endpoints.UpdateUserEndpoint = retry
	}

	{
		factory := userserviceFactory(userservice.MakeAddPermissionForUserEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(3, 500*time.Millisecond, balancer)
		retry = userBreaker(retry)
		endpoints.AddPermissionForUserEndpoint = retry
	}

	{
		factory := userserviceFactory(userservice.MakeRemovePermissionForUserEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(3, 500*time.Millisecond, balancer)
		retry = userBreaker(retry)
		endpoints.RemovePermissionForUserEndpoint = retry
	}

	{
		factory := userserviceFactory(userservice.MakeRemoveRoleForUserEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(3, 500*time.Millisecond, balancer)
		retry = userBreaker(retry)
		endpoints.RemoveRoleForUserEndpoint = retry
	}

	{
		factory := userserviceFactory(userservice.MakeListMenusEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(3, 500*time.Millisecond, balancer)
		retry = userBreaker(retry)
		endpoints.ListMenusEndpoint = retry
	}

	{
		factory := userserviceFactory(userservice.MakeCreateMenuEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(3, 500*time.Millisecond, balancer)
		retry = userBreaker(retry)
		endpoints.CreateMenuEndpoint = retry
	}

	{
		factory := userserviceFactory(userservice.MakeRemoveMenuEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(3, 500*time.Millisecond, balancer)
		retry = userBreaker(retry)
		endpoints.RemoveMenuEndpoint = retry
	}

	{
		factory := userserviceFactory(userservice.MakeGetUserTypeInfoEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(3, 500*time.Millisecond, balancer)
		retry = userBreaker(retry)
		endpoints.GetUserTypeInfoEndpoint = retry
	}

	{
		factory := userserviceFactory(userservice.MakeGetUserEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(3, 500*time.Millisecond, balancer)
		retry = userBreaker(retry)
		endpoints.GetUserEndpoint = retry
	}

	{
		factory := userserviceFactory(userservice.MakeGetUserPermissionsAndRolesEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(3, 500*time.Millisecond, balancer)
		retry = userBreaker(retry)
		endpoints.GetUserPermissionsAndRolesEndpoint = retry
	}

	{
		factory := userserviceFactory(userservice.MakeGetRolePermissionsAndRolesEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(3, 500*time.Millisecond, balancer)
		retry = userBreaker(retry)
		endpoints.GetRolePermissionsAndRolesEndpoint = retry
	}

	{
		factory := userserviceFactory(userservice.MakeGetPermissionsAndRoutesEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(3, 500*time.Millisecond, balancer)
		retry = userBreaker(retry)
		endpoints.GetPermissionsAndRoutesEndpoint = retry
	}

	{
		factory := userserviceFactory(userservice.MakeListLeaguerEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(3, 500*time.Millisecond, balancer)
		retry = userBreaker(retry)
		endpoints.ListLeaguerEndpoint = retry
	}

	{
		factory := userserviceFactory(userservice.MakeRemoveRouteEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(3, 500*time.Millisecond, balancer)
		retry = userBreaker(retry)
		endpoints.RemoveRouteEndpoint = retry
	}

	return &endpoints
}

func userserviceFactory(makeEndpoint func(pb.UserServer) endpoint.Endpoint) sd.Factory {
	return func(instance string) (endpoint endpoint.Endpoint, closer io.Closer, e error) {
		conn, err := grpc.Dial(instance, grpc.WithInsecure())
		if err != nil {
			return nil, nil, err
		}
		addr := os.Getenv("ZIPKIN_ADDR")
		if addr == "" {
			addr = "127.0.0.1:9411"
		}

		localEndpoint, _ := stdzipkin.NewEndpoint("user", "localhost:9411")
		reporter := zipkinhttp.NewReporter("http://" + addr + "/api/v2/spans")
		stdTracer, err := stdzipkin.NewTracer(
			reporter,
			stdzipkin.WithLocalEndpoint(localEndpoint),
		)
		if err != nil {
			logrus.Errorln(err)
		}
		var service *UserEndpoints
		if stdTracer == nil {
			service = NewUserServiceGRPCClient(conn, nil)
		} else {
			tracer := zipkin.GRPCClientTrace(stdTracer)
			service = NewUserServiceGRPCClient(conn, tracer)
		}

		return makeEndpoint(service), conn, nil
	}
}
