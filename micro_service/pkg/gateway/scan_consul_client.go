package gateway

import (
	"io"
	"userService/pkg/apstfr/apstfrpb"
	"userService/pkg/apstfr/scan"

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

func GetScanEndpoints(instancer sd.Instancer, log log.Logger) *ScanEndpoints {
	var endpoints ScanEndpoints

	hystrix.ConfigureCommand(scanBreaker, hystrix.CommandConfig{
		MaxConcurrentRequests: 1000,
		Timeout:               10000,
		ErrorPercentThreshold: 25,
		SleepWindow:           10000,
	})

	breaker := circuitbreaker.Hystrix(scanBreaker)

	{
		factory := scanFactory(scan.MakePayEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(rpcRetryTimes, rpcTimeOut, balancer)
		retry = breaker(retry)
		endpoints.PayEndpoint = retry
	}

	return &endpoints
}

func scanFactory(makeEndpoint func(apstfrpb.ScanServer) endpoint.Endpoint) sd.Factory {
	return func(instance string) (endpoint endpoint.Endpoint, closer io.Closer, e error) {
		conn, err := grpc.Dial(instance, grpc.WithInsecure())
		if err != nil {
			return nil, nil, err
		}

		localEndpoint, _ := stdzipkin.NewEndpoint("scan", "localhost:9411")
		reporter := zipkinhttp.NewReporter("http://localhost:9411/api/v2/spans")
		stdTracer, err := stdzipkin.NewTracer(
			reporter,
			stdzipkin.WithLocalEndpoint(localEndpoint),
		)
		if err != nil {
			logrus.Errorln(err)
		}

		var service *ScanEndpoints
		if stdTracer == nil {
			service = NewScanGRPCClient(conn, nil)
		} else {
			tracer := zipkin.GRPCClientTrace(stdTracer)
			service = NewScanGRPCClient(conn, tracer)
		}

		return makeEndpoint(service), conn, nil
	}
}
