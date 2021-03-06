package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/baxiang/soldiers-sortie/blog/pkg/logging"
	"github.com/baxiang/soldiers-sortie/blog/pkg/redis"
	"github.com/baxiang/soldiers-sortie/blog/pkg/setting"
	"github.com/baxiang/soldiers-sortie/blog/pkg/util"
	"github.com/baxiang/soldiers-sortie/blog/models"
	"github.com/baxiang/soldiers-sortie/blog/routers"
)

func init() {
	setting.Setup()
	models.Setup()
	logging.Setup()
	redis.Setup()
	util.Setup()
}

// @title Golang Gin API
// @version 1.0
// @description An example of gin
// @license.name MIT
func main() {
	gin.SetMode(setting.ServerSetting.RunMode)

	routersInit := routers.InitRouter()
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)

	server := &http.Server{
		Addr:           endPoint,
		Handler:        routersInit,
		ReadTimeout:    setting.ServerSetting.ReadTimeout,
		WriteTimeout:   setting.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("[info] start http server listening %s", endPoint)

	log.Fatal(server.ListenAndServe())

}
