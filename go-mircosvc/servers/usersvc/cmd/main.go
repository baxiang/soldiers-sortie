package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/reflection"

	"github.com/baxiang/soldiers-sortie/go-mircosvc/pkg/db"
	"github.com/baxiang/soldiers-sortie/go-mircosvc/pkg/email"
	"github.com/baxiang/soldiers-sortie/go-mircosvc/pkg/logger"
	"github.com/baxiang/soldiers-sortie/go-mircosvc/servers/usersvc/config"
	"github.com/baxiang/soldiers-sortie/go-mircosvc/servers/usersvc/endpoints"
	"github.com/baxiang/soldiers-sortie/go-mircosvc/servers/usersvc/middleware"
	"github.com/baxiang/soldiers-sortie/go-mircosvc/servers/usersvc/transport"

	kitGrpc "github.com/go-kit/kit/transport/grpc"
	userPb "github.com/baxiang/soldiers-sortie/go-mircosvc/pb"
	sharedZipkin "github.com/baxiang/soldiers-sortie/go-mircosvc/pkg/zipkin"
	sharedEtcd "github.com/baxiang/soldiers-sortie/go-mircosvc/pkg/etcd"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	zipkinGrpc "github.com/openzipkin/zipkin-go/middleware/grpc"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"

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
		register := sharedEtcd.Register("/usersvc", "localhost:"+conf.GrpcPort, etcdClient, log)
		defer register.Register()
	}

	var svc endpoints.UserSerivcer
	{
		mdb := db.NewMysql(conf.MysqlUsername, conf.MysqlPassword, conf.MysqlAddr, conf.MysqlAuthsource)
		rd := db.NewRedis(conf.RedisAddr, conf.RedisPassword, conf.RedisMaxIdle, conf.RedisMaxActive)
		email := email.NewEmail(conf.EmailFrom, conf.EmailAuthCode, conf.EmailHost, conf.EmailSender, conf.EmailPort)
		svc = endpoints.NewUserService(mdb, rd, email)
		svc = middleware.MakeServiceMiddleware(svc)
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

func grpcServer(grpcsvc userPb.UsersvcServer, port string, zipkinTracer *zipkin.Tracer, hs *health.Server,
	logger log.Logger,
	errs chan error) {
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
	userPb.RegisterUsersvcServer(server, grpcsvc)
	healthgrpc.RegisterHealthServer(server, hs)
	reflection.Register(server)
	errs <- server.Serve(listener)
}

