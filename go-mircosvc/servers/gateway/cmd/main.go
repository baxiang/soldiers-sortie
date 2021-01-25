package main

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/baxiang/soldiers-sortie/go-mircosvc/pkg/db"
	sharedEtcd "github.com/baxiang/soldiers-sortie/go-mircosvc/pkg/etcd"
	"github.com/baxiang/soldiers-sortie/go-mircosvc/pkg/logger"
	"github.com/baxiang/soldiers-sortie/go-mircosvc/pkg/session"
	sharedZipkin "github.com/baxiang/soldiers-sortie/go-mircosvc/pkg/zipkin"
	"github.com/baxiang/soldiers-sortie/go-mircosvc/servers/gateway/transport"
	"github.com/go-kit/kit/log/level"
	"github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/baxiang/soldiers-sortie/go-mircosvc/servers/gateway/config"

)

func main() {
	conf := config.GetConfig()

	log, f := logger.NewLogger(conf.LogPath)
	defer f.Close()

	zipkinTracer, reporter := sharedZipkin.NewZipkin(log, conf.ZipkinAddr, "localhost:"+conf.HttpPort,
		conf.ServiceName)
	defer reporter.Close()

	opentracing.SetGlobalTracer(zipkinot.Wrap(zipkinTracer))
	tracer := opentracing.GlobalTracer()
	etcdClient := sharedEtcd.NewEtcd(conf.EtcdAddr)

	var handler http.Handler
	{
		redis := db.NewRedis(conf.RedisAddr, conf.RedisPassword, conf.RedisMaxIdle, conf.RedisMaxActive)
		sessionStorage := session.NewStorage(redis)
		handler = transport.MakeHandler(etcdClient, tracer, zipkinTracer, log, conf.RetryMax, conf.RetryTimeout, sessionStorage)
	}

	errs := make(chan error, 1)
	go httpServer(log, conf.HttpPort, accessControl(handler, conf.Origin), errs)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	level.Info(log).Log("serviceName", conf.ServiceName, "terminated", <-errs)
}

func httpServer(lg log.Logger, port string, handler http.Handler, errs chan error) {
	svr := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: handler,
	}
	err := svr.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		lg.Log("listen: %s\n", err)
	}
	errs <- err
}

func accessControl(h http.Handler, origin string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin == "*" {
			origin = r.Header.Get("Origin")
		}

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
