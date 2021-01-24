package main

import (
	"github.com/baxiang/soldiers-sortie/micro_service/pkg/common"
	"github.com/baxiang/soldiers-sortie/micro_service/pkg/util"
	logrustash "github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/sirupsen/logrus"
	"os"

	"github.com/spf13/viper"
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

func main(){
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

}