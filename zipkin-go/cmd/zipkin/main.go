package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/baxiang/soldiers-sortie/zipkin-go/middleware"
	"github.com/baxiang/soldiers-sortie/zipkin-go/pb"
	"github.com/baxiang/soldiers-sortie/zipkin-go/service"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
	"github.com/openzipkin/zipkin-go"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"google.golang.org/grpc"


	kitzipkin "github.com/go-kit/kit/tracing/zipkin"
	"github.com/baxiang/soldiers-sortie/zipkin-go/transports"
	edpts "github.com/baxiang/soldiers-sortie/zipkin-go/endpoint"
	"github.com/pborman/uuid"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
)

func main() {
	var (
		consulHost  = flag.String("consul.host", "114.67.98.210", "consul ip address")
		consulPort  = flag.String("consul.port", "8500", "consul port")
		serviceHost = flag.String("service.host", "localhost", "service ip address")
		servicePort = flag.String("service.port", "9009", "service port")
		zipkinURL   = flag.String("zipkin.url", "http://114.67.98.210:9411/api/v2/spans", "Zipkin server url")
		grpcAddr    = flag.String("grpc", ":9008", "gRPC listen address.")
	)

	flag.Parse()

	ctx := context.Background()
	errChan := make(chan error)

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var zipkinTracer *zipkin.Tracer
	{
		var (
			err           error
			hostPort      = *serviceHost + ":" + *servicePort
			serviceName   = "string-service"
			useNoopTracer = (*zipkinURL == "")
			reporter      = zipkinhttp.NewReporter(*zipkinURL)
		)
		defer reporter.Close()
		zEP, _ := zipkin.NewEndpoint(serviceName, hostPort)
		zipkinTracer, err = zipkin.NewTracer(
			reporter, zipkin.WithLocalEndpoint(zEP), zipkin.WithNoopTracer(useNoopTracer),
		)
		if err != nil {
			logger.Log("err", err)
			os.Exit(1)
		}
		if !useNoopTracer {
			logger.Log("tracer", "Zipkin", "type", "Native", "URL", *zipkinURL)
		}
	}

	var svc service.Service
	svc = service.StringService{}

	// add logging middleware to service
	svc = middleware.LoggingMiddleware(logger)(svc)

	endpoint := edpts.MakeStringEndpoint(ctx, svc)
	endpoint = kitzipkin.TraceEndpoint(zipkinTracer, "string-endpoint")(endpoint)

	//创建健康检查的Endpoint
	healthEndpoint := edpts.MakeHealthCheckEndpoint(svc)
	healthEndpoint = kitzipkin.TraceEndpoint(zipkinTracer, "health-endpoint")(healthEndpoint)

	//把算术运算Endpoint和健康检查Endpoint封装至StringEndpoints
	endpts := edpts.StringEndpoints{
		StringEndpoint:      endpoint,
		HealthCheckEndpoint: healthEndpoint,
	}

	//创建http.Handler
	r := transports.MakeHttpHandler(ctx, endpts, zipkinTracer, logger)

	//创建注册对象
	registar := Register(*consulHost, *consulPort, *serviceHost, *servicePort, logger)

	go func() {
		fmt.Println("Http Server start at port:" + *servicePort)
		//启动前执行注册
		registar.Register()
		handler := r
		errChan <- http.ListenAndServe(":"+*servicePort, handler)
	}()
	//grpc server
	go func() {
		fmt.Println("grpc Server start at port" + *grpcAddr)
		listener, err := net.Listen("tcp", *grpcAddr)
		if err != nil {
			errChan <- err
			return
		}
		serverTracer := kitzipkin.GRPCServerTrace(zipkinTracer, kitzipkin.Name("string-grpc-transport"))

		handler := transports.NewGRPCServer(ctx, endpts, serverTracer)
		gRPCServer := grpc.NewServer()
		pb.RegisterStringServiceServer(gRPCServer, handler)
		errChan <- gRPCServer.Serve(listener)
	}()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	error := <-errChan
	//服务退出取消注册
	registar.Deregister()
	fmt.Println(error)
}

func Register(consulHost, consulPort, svcHost, svcPort string, logger log.Logger) (registar sd.Registrar) {

	// 创建Consul客户端连接
	var client consul.Client
	{
		consulCfg := api.DefaultConfig()
		consulCfg.Address = consulHost + ":" + consulPort
		consulClient, err := api.NewClient(consulCfg)
		if err != nil {
			logger.Log("create consul client error:", err)
			os.Exit(1)
		}

		client = consul.NewClient(consulClient)
	}

	// 设置Consul对服务健康检查的参数
	check := api.AgentServiceCheck{
		HTTP:     "http://" + svcHost + ":" + svcPort + "/health",
		Interval: "10s",
		Timeout:  "1s",
		Notes:    "Consul check service health status.",
	}

	port, _ := strconv.Atoi(svcPort)

	//设置微服务想Consul的注册信息
	reg := api.AgentServiceRegistration{
		ID:      "string-service" + uuid.New(),
		Name:    "string-service",
		Address: svcHost,
		Port:    port,
		Tags:    []string{"string-service"},
		Check:   &check,
	}

	// 执行注册
	registar = consul.NewRegistrar(client, &reg, logger)
	return
}