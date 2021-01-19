package main

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"userService/pkg/camunda"
	"userService/pkg/workflow"

	"github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/tracing/zipkin"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	consuld "github.com/hashicorp/consul/api"
	stdzipkin "github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"userService/pkg/common"
	"userService/pkg/model"
	"userService/pkg/pb"
	"userService/pkg/util"
)

var (
	mysqlHost     = "localhost"
	mysqlPort     = 3306
	mysqlUser     = "test"
	mysqlPassword = "test"
	mysqlDB       = "test"
	grpcHost      = "localhost"
	grpcPort      = 5001

	consulHost = "localhost"
	consulPort = 8500
)

func main() {
	// 初始化log
	level := os.Getenv("LOG_LEVEL")
	if level == "debug" {
		logrus.SetReportCaller(true)
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	logrus.SetFormatter(&util.LogFormatter{})

	var err error
	if err = parseConfigFile(); err != nil {
		logrus.Fatal("解析配置文件错误", err)
	}

	// 初始化mysql client
	opts := model.Options{
		User:     mysqlUser,
		Password: mysqlPassword,
		DB:       mysqlDB,
		Addr:     fmt.Sprintf("%s:%d", mysqlHost, mysqlPort),
	}
	common.DB, err = model.NewDB(&opts)
	if err != nil {
		logrus.Fatal("启动mysql错误", err)
	}
	if level == "debug" {
		common.DB = common.DB.Debug()
	}

	// 初始化consul client
	consulClient, err := newConsulClient(fmt.Sprintf("%s:%d", consulHost, consulPort))
	if err != nil {
		logrus.Fatal("连接consul失败", err)
	}

	// 启动grpc server
	go func() {
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
		var tracer grpctransport.ServerOption

		if stdTracer != nil {
			tracer = zipkin.GRPCServerTrace(stdTracer)
		}
		if err := runGRPCServer(fmt.Sprintf("%s:%d", "", grpcPort), tracer); err != nil {
			logrus.Fatal("grpc server shutdown", err)
		}
	}()

	// 注册consul service
	err = registerConsulService(consulClient, "workflow", grpcHost, grpcPort)
	if err != nil {
		logrus.Errorln("注册workflow失败", err)
	}
	logrus.Infoln("启动成功")

	log := &util.ConsulLogger{}
	camunda.Load(consulClient, log)

	runtime.Goexit()
}

func parseConfigFile() error {
	fileName := os.Getenv("CONFIG_FILE")

	viper.SetConfigType("toml")
	viper.AddConfigPath("./configs/workflow/")

	if fileName != "" {
		viper.SetConfigFile(fileName)
	}

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	settings := viper.AllSettings()
	logrus.Infoln(settings)

	mysqlHost = viper.GetString("db.mysql.host")
	mysqlPort = viper.GetInt("db.mysql.port")
	mysqlUser = viper.GetString("db.mysql.user")
	mysqlPassword = viper.GetString("db.mysql.password")
	mysqlDB = viper.GetString("db.mysql.db")

	grpcHost = viper.GetString("grpc.host")
	grpcPort = viper.GetInt("grpc.port")

	if host := os.Getenv("GRPC_HOST"); host != "" {
		grpcHost = host
	}

	consulHost = viper.GetString("consul.host")
	consulPort = viper.GetInt("consul.port")
	return nil
}

func newConsulClient(addr string) (consul.Client, error) {
	consulClient, err := consuld.NewClient(&consuld.Config{
		Address: addr,
	})
	if err != nil {
		return nil, err
	}
	return consul.NewClient(consulClient), nil
}

func registerConsulService(client consul.Client, name string, host string, port int) error {
	return client.Register(&consuld.AgentServiceRegistration{
		Address: host,
		Port:    port,
		Name:    name,
	})
}

func runGRPCServer(addr string, tracer grpctransport.ServerOption) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	svr := grpc.NewServer()
	pb.RegisterWorkflowServer(svr, workflow.New(tracer))
	return svr.Serve(l)
}
