package gateway

import (
	"io"
	"os"
	"userService/pkg/pb"
	"userService/pkg/workflow"

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

func GetWorkflowEndpoints(instancer sd.Instancer, log log.Logger) *WorkflowEndpoints {
	var endpoints WorkflowEndpoints

	hystrix.ConfigureCommand(workflowBreaker, hystrix.CommandConfig{
		MaxConcurrentRequests: 1000,
		Timeout:               10000,
		ErrorPercentThreshold: 25,
		SleepWindow:           10000,
	})
	userBreaker := circuitbreaker.Hystrix(workflowBreaker)

	{
		factory := workflowFactory(workflow.MakeListTaskEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(rpcRetryTimes, rpcTimeOut, balancer)
		retry = userBreaker(retry)
		endpoints.ListTaskEndpoint = retry
	}

	{
		factory := workflowFactory(workflow.MakeHandleTaskEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(rpcRetryTimes, rpcTimeOut, balancer)
		retry = userBreaker(retry)
		endpoints.HandleTaskEndpoint = retry
	}

	{
		factory := workflowFactory(workflow.MakeStartWorkflowEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(rpcRetryTimes, rpcTimeOut, balancer)
		retry = userBreaker(retry)
		endpoints.StartEndpoint = retry
	}

	{
		factory := workflowFactory(workflow.MakeListRemarkEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(rpcRetryTimes, rpcTimeOut, balancer)
		retry = userBreaker(retry)
		endpoints.ListRemarkEndpoint = retry
	}

	return &endpoints
}

func workflowFactory(makeEndpoint func(pb.WorkflowServer) endpoint.Endpoint) sd.Factory {
	return func(instance string) (endpoint endpoint.Endpoint, closer io.Closer, e error) {
		conn, err := grpc.Dial(instance, grpc.WithInsecure())
		if err != nil {
			return nil, nil, err
		}
		addr := os.Getenv("ZIPKIN_ADDR")
		if addr == "" {
			addr = "127.0.0.1:9411"
		}

		localEndpoint, _ := stdzipkin.NewEndpoint("workflow", "localhost:9411")
		reporter := zipkinhttp.NewReporter("http://" + addr + "/api/v2/spans")
		stdTracer, err := stdzipkin.NewTracer(
			reporter,
			stdzipkin.WithLocalEndpoint(localEndpoint),
		)
		if err != nil {
			logrus.Errorln(err)
		}
		var service *WorkflowEndpoints
		if stdTracer == nil {
			service = NewWorkflowGRPCClient(conn, nil)
		} else {
			tracer := zipkin.GRPCClientTrace(stdTracer)
			service = NewWorkflowGRPCClient(conn, tracer)
		}

		return makeEndpoint(service), conn, nil
	}
}
