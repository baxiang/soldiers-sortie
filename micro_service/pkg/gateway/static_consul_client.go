package gateway

import (
	"io"
	"userService/pkg/pb"
	"userService/pkg/staticservice"

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

func GetStaticCliEndpoints(instancer sd.Instancer, log log.Logger) *StaticEndpoints {

	var endpoints StaticEndpoints

	hystrix.ConfigureCommand(staticBreaker, hystrix.CommandConfig{
		MaxConcurrentRequests: 1000,
		Timeout:               10000,
		ErrorPercentThreshold: 25,
		SleepWindow:           10000,
	})
	breaker := circuitbreaker.Hystrix(staticBreaker)

	{
		factory := staticserviceFactory(staticservice.MakeSyncDataEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(rpcRetryTimes, rpcTimeOut, balancer)
		retry = breaker(retry)
		endpoints.SyncDataEndpoint = retry
	}

	{
		factory := staticserviceFactory(staticservice.MakeGetDictionaryItemEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(rpcRetryTimes, rpcTimeOut, balancer)
		retry = breaker(retry)
		endpoints.GetDictionaryItemEndpoint = retry
	}

	{
		factory := staticserviceFactory(staticservice.MakeGetDicByProdAndBizEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(rpcRetryTimes, rpcTimeOut, balancer)
		retry = breaker(retry)
		endpoints.GetDicByProdAndBizEndpoint = retry
	}

	{
		factory := staticserviceFactory(staticservice.MakeCheckValuesEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(rpcRetryTimes, rpcTimeOut, balancer)
		retry = breaker(retry)
		endpoints.CheckValuesEndpoint = retry
	}

	{
		factory := staticserviceFactory(staticservice.MakeGetDictionaryLayerItemEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(rpcRetryTimes, rpcTimeOut, balancer)
		retry = breaker(retry)
		endpoints.GetDictionaryLayerItemEndpoint = retry
	}

	{
		factory := staticserviceFactory(staticservice.MakeGetDictionaryItemByPkEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(rpcRetryTimes, rpcTimeOut, balancer)
		retry = breaker(retry)
		endpoints.GetDictionaryItemByPkEndpoint = retry
	}

	{
		factory := staticserviceFactory(staticservice.MakeGetUnionPayBankListEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(rpcRetryTimes, rpcTimeOut, balancer)
		retry = breaker(retry)
		endpoints.GetUnionPayBankListEndpoint = retry
	}

	{
		factory := staticserviceFactory(staticservice.MakeFindUnionPayMccListEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(rpcRetryTimes, rpcTimeOut, balancer)
		retry = breaker(retry)
		endpoints.FindUnionPayMccListEndpoint = retry
	}

	{
		factory := staticserviceFactory(staticservice.MakeGetInsProdBizFeeMapInfoEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(rpcRetryTimes, rpcTimeOut, balancer)
		retry = breaker(retry)
		endpoints.GetInsProdBizFeeMapInfoEndpoint = retry
	}

	{
		factory := staticserviceFactory(staticservice.MakeListTransMapEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(rpcRetryTimes, rpcTimeOut, balancer)
		retry = breaker(retry)
		endpoints.ListTransMapEndpoint = retry
	}

	{
		factory := staticserviceFactory(staticservice.MakeListFeeMapEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(rpcRetryTimes, rpcTimeOut, balancer)
		retry = breaker(retry)
		endpoints.ListFeeMapEndpoint = retry
	}

	{
		factory := staticserviceFactory(staticservice.MakeFindAreaEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(rpcRetryTimes, rpcTimeOut, balancer)
		retry = breaker(retry)
		endpoints.FindAreaEndpoint = retry
	}

	{
		factory := staticserviceFactory(staticservice.MakeFindMerchantFirstThreeCodeEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(rpcRetryTimes, rpcTimeOut, balancer)
		retry = breaker(retry)
		endpoints.FindMerchantFirstThreeCodeEndpoint = retry
	}

	{
		factory := staticserviceFactory(staticservice.MakeSaveOrgDictionaryItemEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(rpcRetryTimes, rpcTimeOut, balancer)
		retry = breaker(retry)
		endpoints.SaveOrgDictionaryItemEndpoint = retry
	}

	{
		factory := staticserviceFactory(staticservice.MakeListOrgDictionaryItemEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(rpcRetryTimes, rpcTimeOut, balancer)
		retry = breaker(retry)
		endpoints.ListOrgDictionaryItemEndpoint = retry
	}

	{
		factory := staticserviceFactory(staticservice.MakeSaveInsProdBizFeeMapInfoEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, log)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(rpcRetryTimes, rpcTimeOut, balancer)
		retry = breaker(retry)
		endpoints.SaveInsProdBizFeeMapInfoEndpoint = retry
	}

	return &endpoints
}

func staticserviceFactory(makeEndpoint func(pb.StaticServer) endpoint.Endpoint) sd.Factory {
	return func(instance string) (endpoint endpoint.Endpoint, closer io.Closer, e error) {
		conn, err := grpc.Dial(instance, grpc.WithInsecure())
		if err != nil {
			return nil, nil, err
		}

		localEndpoint, _ := stdzipkin.NewEndpoint("static", "localhost:9411")
		reporter := zipkinhttp.NewReporter("http://localhost:9411/api/v2/spans")
		stdTracer, err := stdzipkin.NewTracer(
			reporter,
			stdzipkin.WithLocalEndpoint(localEndpoint),
		)
		if err != nil {
			logrus.Errorln(err)
		}

		var service *StaticEndpoints
		if stdTracer == nil {
			service = NewStaticServiceGRPCClient(conn, nil)
		} else {
			tracer := zipkin.GRPCClientTrace(stdTracer)
			service = NewStaticServiceGRPCClient(conn, tracer)
		}

		return makeEndpoint(service), conn, nil
	}
}
