package gateway

import (
	"io"
	"userService/pkg/institutionservice"
	"userService/pkg/pb"

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

func GetInstitutionCliEndpoints(instancer sd.Instancer, log log.Logger) *InstitutionEndpoints {
	var endpoints InstitutionEndpoints

	hystrix.ConfigureCommand(institutionBreaker, hystrix.CommandConfig{
		MaxConcurrentRequests: 1000,
		Timeout:               10000,
		ErrorPercentThreshold: 25,
		SleepWindow:           10000,
	})
	institutionBreaker := circuitbreaker.Hystrix(institutionBreaker)

	{
		factory := institutionserviceFactory(institutionservice.MakeTnxHisDownloadEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(rpcRetryTimes, rpcTimeOut, balancer)
		retry = institutionBreaker(retry)
		endpoints.TnxHisDownloadEndpoint = retry
	}

	{
		factory := institutionserviceFactory(institutionservice.MakeGetTfrTrnLogsEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(rpcRetryTimes, rpcTimeOut, balancer)
		retry = institutionBreaker(retry)
		endpoints.GetTfrTrnLogsEndpoint = retry
	}

	{
		factory := institutionserviceFactory(institutionservice.MakeGetTfrTrnLogEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(rpcRetryTimes, rpcTimeOut, balancer)
		retry = institutionBreaker(retry)
		endpoints.GetTfrTrnLogEndpoint = retry
	}

	{
		factory := institutionserviceFactory(institutionservice.MakeDownloadTfrTrnLogsEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(rpcRetryTimes, rpcTimeOut, balancer)
		retry = institutionBreaker(retry)
		endpoints.DownloadTfrTrnLogsEndpoint = retry
	}

	{
		factory := institutionserviceFactory(institutionservice.MakeListGroupsEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(rpcRetryTimes, rpcTimeOut, balancer)
		retry = institutionBreaker(retry)
		endpoints.ListGroupsEndpoint = retry
	}

	{
		factory := institutionserviceFactory(institutionservice.MakeListInstitutionsEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(rpcRetryTimes, rpcTimeOut, balancer)
		retry = institutionBreaker(retry)
		endpoints.ListInstitutionsEndpoint = retry
	}

	{
		factory := institutionserviceFactory(institutionservice.MakeSaveInstitutionEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(rpcRetryTimes, rpcTimeOut, balancer)
		retry = institutionBreaker(retry)
		endpoints.SaveInstitutionEndpoint = retry
	}

	{
		factory := institutionserviceFactory(institutionservice.MakeGetInstitutionByIdEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(rpcRetryTimes, rpcTimeOut, balancer)
		retry = institutionBreaker(retry)
		endpoints.GetInstitutionByIdEndpoint = retry
	}

	{
		factory := institutionserviceFactory(institutionservice.MakeSaveInstitutionFeeControlCashEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(rpcRetryTimes, rpcTimeOut, balancer)
		retry = institutionBreaker(retry)
		endpoints.SaveInstitutionFeeControlCashEndpoint = retry
	}

	{
		factory := institutionserviceFactory(institutionservice.MakeGetInstitutionControlEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(rpcRetryTimes, rpcTimeOut, balancer)
		retry = institutionBreaker(retry)
		endpoints.GetInstitutionControlEndpoint = retry
	}

	{
		factory := institutionserviceFactory(institutionservice.MakeGetInstitutionCashEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(rpcRetryTimes, rpcTimeOut, balancer)
		retry = institutionBreaker(retry)
		endpoints.GetInstitutionCashEndpoint = retry
	}

	{
		factory := institutionserviceFactory(institutionservice.MakeGetInstitutionFeeEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(rpcRetryTimes, rpcTimeOut, balancer)
		retry = institutionBreaker(retry)
		endpoints.GetInstitutionFeeEndpoint = retry
	}

	{
		factory := institutionserviceFactory(institutionservice.MakeSaveGroupEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(rpcRetryTimes, rpcTimeOut, balancer)
		retry = institutionBreaker(retry)
		endpoints.SaveGroupEndpoint = retry
	}

	{
		factory := institutionserviceFactory(institutionservice.MakeBindGroupEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(rpcRetryTimes, rpcTimeOut, balancer)
		retry = institutionBreaker(retry)
		endpoints.BindGroupEndpoint = retry
	}

	{
		factory := institutionserviceFactory(institutionservice.MakeListBindGroupEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(rpcRetryTimes, rpcTimeOut, balancer)
		retry = institutionBreaker(retry)
		endpoints.ListBindGroupEndpoint = retry
	}

	{
		factory := institutionserviceFactory(institutionservice.MakeRemoveBindGroupEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(rpcRetryTimes, rpcTimeOut, balancer)
		retry = institutionBreaker(retry)
		endpoints.RemoveBindGroupEndpoint = retry
	}

	return &endpoints
}

func institutionserviceFactory(makeEndpoint func(pb.InstitutionServer) endpoint.Endpoint) sd.Factory {
	return func(instance string) (endpoint endpoint.Endpoint, closer io.Closer, e error) {
		conn, err := grpc.Dial(instance, grpc.WithInsecure())
		if err != nil {
			return nil, nil, err
		}

		localEndpoint, _ := stdzipkin.NewEndpoint("institution", "localhost:9411")
		reporter := zipkinhttp.NewReporter("http://localhost:9411/api/v2/spans")
		stdTracer, err := stdzipkin.NewTracer(
			reporter,
			stdzipkin.WithLocalEndpoint(localEndpoint),
		)
		if err != nil {
			logrus.Errorln(err)
		}

		var service *InstitutionEndpoints
		if stdTracer == nil {
			service = NewInstitutionServiceGRPCClient(conn, nil)
		} else {
			tracer := zipkin.GRPCClientTrace(stdTracer)
			service = NewInstitutionServiceGRPCClient(conn, tracer)
		}

		return makeEndpoint(service), conn, nil
	}
}
