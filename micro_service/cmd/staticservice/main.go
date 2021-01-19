package main

import (
	"bytes"
	"fmt"
	"os"
	"userService/pkg/staticservice"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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

	logrus.Info("正在链接mysql...")

	logrus.Info("启动consul watcher ...")
	go staticservice.StartServer(conf.WatcherAddr)

	//register service.
	logrus.Info("启动成功")
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
	WatcherAddr    string
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
	viper.AddConfigPath("./cmd/staticservice")

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

	conf.WatcherAddr = viper.GetString("watcher.addr")

	conf.ServiceName = viper.GetString("info.serviceName")

	if host := os.Getenv("GRPC_HOST"); host != "" {
		conf.GrpcHost = host
	}

	conf.ConsulHost = viper.GetString("consul.host")
	conf.ConsulPort = viper.GetInt("consul.port")

	return &conf, err
}
