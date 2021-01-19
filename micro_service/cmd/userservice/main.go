package main

import (
	"fmt"
	"net"
	"os"
	"path"
	"runtime"
	"time"
	"userService/pkg/apstfr/apstfrpb"
	"userService/pkg/apstfr/scan"
	"userService/pkg/camunda"
	"userService/pkg/institutionservice"
	"userService/pkg/merchantservice"
	"userService/pkg/staticservice"
	"userService/pkg/task"
	"userService/pkg/termservice"
	"userService/pkg/workflow"

	logrustash "github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/tracing/zipkin"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/go-redis/redis"
	consuld "github.com/hashicorp/consul/api"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	stdzipkin "github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"userService/pkg/common"
	"userService/pkg/model"
	"userService/pkg/pb"
	"userService/pkg/rbac"
	"userService/pkg/userservice"
	"userService/pkg/util"
)

var (
	mysqlHost     = "localhost"
	mysqlPort     = 3306
	mysqlUser     = "test"
	mysqlPassword = "test"
	mysqlDB       = "test"

	redisHost = "localhost"
	redisPort = 6379

	grpcHost = "localhost"
	grpcPort = 5001

	consulHost = "localhost"
	consulPort = 8500

	rbacFileName = ""

	logPath = ""
	logFile = ""

	watcherAddr = ""
	taskCron    = ""
)

func main() {
	hook, err := logrustash.NewHook("tcp", "127.0.0.1:5100", "userService")
	if err != nil {
		logrus.Errorln(err)
	} else {
		logrus.AddHook(hook)
	}

	// 初始化log
	level := os.Getenv("LOG_LEVEL")
	if level == "debug" {
		logrus.SetReportCaller(true)
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	logrus.SetFormatter(&util.LogFormatter{})

	if err = parseConfigFile(); err != nil {
		logrus.Fatal("解析配置文件错误", err)
	}

	if logFile != "" {
		os.MkdirAll(logPath, os.ModePerm)
		logFilePath := path.Join(logPath, logFile)
		writer, err := rotatelogs.New(
			logFilePath+".%Y%m%d%H%M",
			rotatelogs.WithLinkName(logFilePath),
			rotatelogs.WithMaxAge(time.Duration(30*24)*time.Hour),
			rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
		)
		if err != nil {
			logrus.Errorln(err)
		}
		logrus.AddHook(lfshook.NewHook(
			lfshook.WriterMap{
				logrus.InfoLevel:  writer,
				logrus.DebugLevel: writer,
				logrus.FatalLevel: writer,
				logrus.PanicLevel: writer,
				logrus.ErrorLevel: writer,
				logrus.WarnLevel:  writer,
				logrus.TraceLevel: writer,
			},
			&util.LogFormatter{},
		))

	}

	// 初始化redis client
	common.RedisClient = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", redisHost, redisPort),
	})

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

	common.Enforcer = rbac.NewCasbin(rbacFileName, &opts)

	// 初始化consul client
	consulClient, err := newConsulClient(fmt.Sprintf("%s:%d", consulHost, consulPort))
	if err != nil {
		logrus.Fatal("连接consul失败", err)
	}

	// 启动 http watcher
	go staticservice.StartServer(watcherAddr)

	// 启动 grpc server
	go func() {
		addr := os.Getenv("ZIPKIN_ADDR")
		if addr == "" {
			addr = "127.0.0.1:9411"
		}

		localEndpoint, _ := stdzipkin.NewEndpoint("user", "localhost:9411")
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

	log := &util.ConsulLogger{}
	camunda.Load(consulClient, log)
	// 注册consul service
	err = registerConsulService(consulClient, "userService", grpcHost, grpcPort)
	if err != nil {
		logrus.Errorln("注册userService失败", err)
	}
	logrus.Infoln("启动成功")

	// 启动定时任务
	task.RunServiceTask(taskCron, 4)

	runtime.Goexit()
}

func parseConfigFile() error {
	fileName := os.Getenv("CONFIG_FILE")

	viper.SetConfigType("toml")
	viper.AddConfigPath("./configs/")

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

	redisHost = viper.GetString("db.redis.host")
	redisPort = viper.GetInt("db.redis.port")

	grpcHost = viper.GetString("grpc.host")
	grpcPort = viper.GetInt("grpc.port")

	if host := os.Getenv("GRPC_HOST"); host != "" {
		grpcHost = host
	}

	consulHost = viper.GetString("consul.host")
	consulPort = viper.GetInt("consul.port")

	rbacFileName = viper.GetString("rbacFile")

	logPath = viper.GetString("log.path")
	logFile = viper.GetString("log.file")

	watcherAddr = viper.GetString("watcher.addr")

	taskCron = viper.GetString("task.cron")

	signKey := os.Getenv("SIGN_KEY")
	if signKey != "" {
		common.SignKey = signKey
	}
	return nil
}

func newConsulClient(addr string) (consul.Client, error) {
	consulClient, err := consuld.NewClient(&consuld.Config{
		Address: addr,
	})
	if err != nil {
		return nil, err
	}
	common.ConsulClient = consulClient
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
	pb.RegisterUserServer(svr, userservice.New(tracer))
	pb.RegisterWorkflowServer(svr, workflow.New(tracer))
	pb.RegisterInstitutionServer(svr, institutionservice.New(tracer))
	pb.RegisterMerchantServer(svr, merchantservice.New(tracer))
	pb.RegisterTermServer(svr, termservice.New(tracer))
	pb.RegisterStaticServer(svr, staticservice.New(tracer))

	apstfrpb.RegisterScanServer(svr, scan.New(tracer))
	return svr.Serve(l)
}
