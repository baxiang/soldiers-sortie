package main

import (
	"net/http"
	"os"
	"userService/pkg/gateway"
	"userService/pkg/util"

	consulsd "github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
)

func main() {
	level := os.Getenv("LOG_LEVEL")
	if level == "debug" {
		logrus.SetLevel(logrus.DebugLevel)
	}
	logrus.SetFormatter(&util.LogFormatter{})
	port := os.Getenv("GATEWAY_PORT")
	if port == "" {
		port = "8080"
	}
	consulAddr := os.Getenv("CONSUL_ADDR")

	consulConfig := api.DefaultConfig()
	if consulAddr != "" {
		consulConfig.Address = consulAddr
	}
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		logrus.Fatal(err)
	}
	client := consulsd.NewClient(consulClient)
	log := &logger{}

	var (
		tags        = []string{}
		passingOnly = true
		instancer   = consulsd.NewInstancer(client, log, "userService", tags, passingOnly)
	)
	endpoints := new(gateway.ClientEndpoints)
	endpoints.UserEndpoints = gateway.GetUserEndpoints(instancer, log)
	endpoints.StaticEndpoints = gateway.GetStaticCliEndpoints(instancer, log)
	endpoints.InstitutionEndpoints = gateway.GetInstitutionCliEndpoints(instancer, log)
	endpoints.MerchantEndpoints = gateway.GetMerchantEndpoints(instancer, log)
	endpoints.TermEndpoints = gateway.GetTermEndpoints(instancer, log)
	endpoints.WorkflowEndpoints = gateway.GetWorkflowEndpoints(instancer, log)
	endpoints.ScanEndpoints = gateway.GetScanEndpoints(instancer, log)

	userHandler := gateway.NewHttpHandler(endpoints)
	http.ListenAndServe(":"+port, userHandler)
}

type logger struct{}

func (l *logger) Log(args ...interface{}) error {
	logrus.Infoln(args...)
	return nil
}
