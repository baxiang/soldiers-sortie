package camunda

import (
	"io"
	"time"
	"userService/pkg/camunda/pb"
	"userService/pkg/transport"

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

const (
	processDefinitionBreaker = "process_definition_breaker"
)

func GetProcessDefinitionConsulEndpoints(instancer sd.Instancer, log log.Logger) *ProcessDefinitionEndpoints {
	var endpoints ProcessDefinitionEndpoints

	hystrix.ConfigureCommand(processDefinitionBreaker, hystrix.CommandConfig{
		MaxConcurrentRequests: 1000,
		Timeout:               10000,
		ErrorPercentThreshold: 25,
		SleepWindow:           10000,
	})
	institutionBreaker := circuitbreaker.Hystrix(processDefinitionBreaker)

	{
		factory := processDefinitionFactory(transport.MakeGetProcessDefinitionEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(3, 5000*time.Millisecond, balancer)
		retry = institutionBreaker(retry)
		endpoints.GetEndpoint = retry
	}

	{
		factory := processDefinitionFactory(transport.MakeStartProcessDefinitionEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(3, 5000*time.Millisecond, balancer)
		retry = institutionBreaker(retry)
		endpoints.StartEndpoint = retry
	}

	return &endpoints
}

func processDefinitionFactory(makeEndpoint func(pb.ProcessDefinitionServer) endpoint.Endpoint) sd.Factory {
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

		var service *ProcessDefinitionEndpoints
		if stdTracer == nil {
			service = NewProcessDefinitionClient(conn, nil)
		} else {
			tracer := zipkin.GRPCClientTrace(stdTracer)
			service = NewProcessDefinitionClient(conn, tracer)
		}

		return makeEndpoint(service), conn, nil
	}
}
