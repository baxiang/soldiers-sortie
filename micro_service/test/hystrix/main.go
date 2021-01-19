package main

import (
	"database/sql"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

func main() {
	_, err := sql.Open("mysql", "micro:micro@tcp(172.16.7.120:3306)/apsmgm?charset=utf8&parseTime=true")
	logrus.Fatal("连接结果", err)

	hystrix.ConfigureCommand("my", hystrix.CommandConfig{
		Timeout:               1000,
		MaxConcurrentRequests: 100,
		ErrorPercentThreshold: 25,
	})
	err = hystrix.Do("my", func() error {
		time.Sleep(10 * time.Second)
		return nil
	}, nil)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Infoln("my success")
}
