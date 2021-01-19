package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"path"
	"runtime"
	"time"
	"userService/pkg/camunda"
	"userService/pkg/common"
	"userService/pkg/institutionservice"
	"userService/pkg/model"
	"userService/pkg/pb"
	"userService/pkg/util"

	"github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var (
	logPath = "./log"
	logFile = "log"
)

func main() {
	//初始化log
	level := os.Getenv("LOG_LEVEL")
	if level == "debug" {
		logrus.SetReportCaller(true)
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	logrus.SetFormatter(&logFormatter{})

	conf, err := ParseConfigFile()
	if err != nil {
		logrus.Fatal("解析配置文件错误", err)
	}

	if logFile != "" {
		os.MkdirAll(logPath, os.ModePerm)
		logFilePath := path.Join(logPath, logFile)
		writer, err := rotatelogs.New(
			logFilePath+".%Y%m%d%H%M",
			rotatelogs.WithLinkName(logFilePath),
			rotatelogs.WithMaxAge(time.Duration(24)*time.Hour),
			rotatelogs.WithRotationTime(time.Duration(30*24)*time.Hour),
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
			&logFormatter{},
		))

	}

	logrus.Info("正在链接mysql...")
	common.DB, err = model.NewDB(&model.Options{
		User:     conf.MysqlUser,
		Password: conf.MysqlPassword,
		DB:       conf.MysqlDB,
		Addr:     fmt.Sprintf("%s:%d", conf.MysqlHost, conf.MysqlPort),
	})

	if err != nil {
		logrus.Fatal("启动mysql错误", err)
	}
	if level == "debug" {
		common.DB = common.DB.Debug()
	}

	go func() {
		if err := runGRPCServer(fmt.Sprintf("%s:%d", conf.GrpcHost, conf.GrpcPort)); err != nil {
			logrus.Fatal("grpc server shutdown", err)
		}
	}()

	// register service.
	logrus.Info("正在链接consul...")
	consulClient, err := newConsulClient(fmt.Sprintf("%s:%d", conf.ConsulHost, conf.ConsulPort))
	if err != nil {
		logrus.Fatal("consul链接失败： ", err)
	}

	err = registerService(consulClient, conf.ServiceName, conf.GrpcRegistHost, conf.GrpcRegistPort)
	if err != nil {
		logrus.Fatal("服务注册失败:", err)
	}

	logrus.Info("启动成功")

	// 加载 camunda grpc 服务
	log := &util.ConsulLogger{}
	camunda.Load(consulClient, log)

	runtime.Goexit()
}

func runGRPCServer(addr string) error {

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	insHandle := institutionservice.NewGRPCServer()

	svr := grpc.NewServer()
	pb.RegisterInstitutionServer(svr, insHandle)
	return svr.Serve(l)
}

func newConsulClient(addr string) (consul.Client, error) {
	client, err := api.NewClient(&api.Config{
		Address: addr,
	})
	if err != nil {
		return nil, err
	}

	return consul.NewClient(client), nil

}

func registerService(client consul.Client, name, host string, port int) error {
	return client.Register(&api.AgentServiceRegistration{
		Address: host,
		Port:    port,
		Name:    name,
	})
}

//Conf .
type Conf struct {
	MysqlHost      string
	MysqlPort      int
	MysqlUser      string
	MysqlPassword  string
	MysqlDB        string
	RedisHost      string
	RedisPort      int
	GrpcHost       string
	GrpcPort       int
	GrpcRegistHost string
	GrpcRegistPort int
	ConsulHost     string
	ConsulPort     int
	ServiceName    string
}

type logFormatter struct{}

func (l logFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var buffer *bytes.Buffer
	if entry.Buffer != nil {
		buffer = entry.Buffer
	} else {
		buffer = &bytes.Buffer{}
	}
	buffer.Write([]byte("["))
	buffer.Write([]byte(entry.Time.Format("2006-01-02 15:04:05.000")))
	buffer.Write([]byte("] "))
	buffer.Write([]byte("["))
	buffer.Write([]byte(entry.Level.String()))
	buffer.Write([]byte("] "))
	if entry.HasCaller() {
		buffer.Write([]byte("["))
		buffer.Write([]byte(fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line)))
		buffer.Write([]byte("] "))
	}
	buffer.Write([]byte(entry.Message))
	buffer.Write([]byte("\n"))
	return buffer.Bytes(), nil
}

//ParseConfigFile .
func ParseConfigFile() (*Conf, error) {
	fileName := os.Getenv("CONFIG_FILE")

	viper.SetConfigType("toml")
	viper.AddConfigPath("./configs/")

	if fileName != "" {
		viper.SetConfigFile(fileName)
	}

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	settings := viper.AllSettings()
	logrus.Infoln(settings)

	conf := Conf{}
	conf.MysqlHost = viper.GetString("db.mysql.host")
	conf.MysqlPort = viper.GetInt("db.mysql.port")
	conf.MysqlUser = viper.GetString("db.mysql.user")
	conf.MysqlPassword = viper.GetString("db.mysql.password")
	conf.MysqlDB = viper.GetString("db.mysql.db")

	conf.RedisHost = viper.GetString("db.redis.host")
	conf.RedisPort = viper.GetInt("db.redis.port")

	conf.GrpcHost = viper.GetString("grpc.host")
	conf.GrpcPort = viper.GetInt("grpc.port")
	conf.GrpcRegistHost = viper.GetString("grpc.registHost")
	conf.GrpcRegistPort = viper.GetInt("grpc.registPost")

	conf.ServiceName = viper.GetString("info.serviceName")

	if host := os.Getenv("GRPC_HOST"); host != "" {
		conf.GrpcHost = host
	}

	conf.ConsulHost = viper.GetString("consul.host")
	conf.ConsulPort = viper.GetInt("consul.port")

	return &conf, err
}
