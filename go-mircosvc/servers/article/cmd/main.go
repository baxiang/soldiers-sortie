package main

import (
	"fmt"
	"github.com/go-kit/kit/log"
	articlePb "github.com/baxiang/soldiers-sortie/go-mircosvc/pb"
	"github.com/baxiang/soldiers-sortie/go-mircosvc/pkg/db"
	sharedEtcd "github.com/baxiang/soldiers-sortie/go-mircosvc/pkg/etcd"
	"github.com/baxiang/soldiers-sortie/go-mircosvc/pkg/logger"
	sharedZipkin "github.com/baxiang/soldiers-sortie/go-mircosvc/pkg/zipkin"
	"github.com/go-kit/kit/log/level"
	kitGrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	zipkinGrpc "github.com/openzipkin/zipkin-go/middleware/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/baxiang/soldiers-sortie/go-mircosvc/servers/article/config"
	"github.com/baxiang/soldiers-sortie/go-mircosvc/servers/article/services"
	"github.com/baxiang/soldiers-sortie/go-mircosvc/servers/article/endpoints"
	"github.com/baxiang/soldiers-sortie/go-mircosvc/servers/article/transport"




)

func main() {
	conf := config.GetConfig()
	log, f := logger.NewLogger(conf.LogPath)
	defer f.Close()

	zipkinTracer, reporter := sharedZipkin.NewZipkin(log, conf.ZipkinAddr, "localhost:"+conf.GrpcPort,
		conf.ServiceName)
	defer reporter.Close()

	opentracing.SetGlobalTracer(zipkinot.Wrap(zipkinTracer))
	tracer := opentracing.GlobalTracer()
	{
		etcdClient := sharedEtcd.NewEtcd(conf.EtcdAddr)
		register := sharedEtcd.Register("/articlesvc", "localhost:"+conf.GrpcPort, etcdClient, log)
		defer register.Register()
	}

	var svc services.ArticleServicer
	{
		mdb := db.NewMysql(conf.MysqlUsername, conf.MysqlPassword, conf.MysqlAddr, conf.MysqlAuthsource)
		svc = services.NewArticleService(mdb)
	}
	eps := endpoints.NewEndpoints(svc, log, tracer, zipkinTracer)

	hs := health.NewServer()
	hs.SetServingStatus(conf.ServiceName, healthgrpc.HealthCheckResponse_SERVING)

	errs := make(chan error, 1)
	go grpcServer(transport.MakeGRPCServer(eps, tracer, zipkinTracer, log), conf.GrpcPort, zipkinTracer, hs, log, errs)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	level.Info(log).Log("serviceName", conf.ServiceName, "terminated", <-errs)
}

func grpcServer(grpcsvc articlePb.ArticlesvcServer, port string, zipkinTracer *zipkin.Tracer, hs *health.Server,
	logger log.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", port)
	listener, err := net.Listen("tcp", p)
	if err != nil {
		level.Error(logger).Log("protocol", "GRPC", "listen", port, "err", err)
		os.Exit(1)
	}
	level.Info(logger).Log("protocol", "GRPC", "protocol", "GRPC", "exposed", port)

	server := grpc.NewServer(grpc.UnaryInterceptor(kitGrpc.Interceptor),
		grpc.StatsHandler(zipkinGrpc.NewServerHandler(zipkinTracer)),
	)
	articlePb.RegisterArticlesvcServer(server, grpcsvc)
	healthgrpc.RegisterHealthServer(server, hs)
	reflection.Register(server)
	errs <- server.Serve(listener)
}