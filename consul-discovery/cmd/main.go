package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	uuid "github.com/satori/go.uuid"

	"github.com/baxiang/soldiers-sortie/consul-discovery/config"
	"github.com/baxiang/soldiers-sortie/consul-discovery/discover"
	"github.com/baxiang/soldiers-sortie/consul-discovery/endpoint"
	"github.com/baxiang/soldiers-sortie/consul-discovery/service"
	"github.com/baxiang/soldiers-sortie/consul-discovery/transport"

)

func main() {
	// 从命令行中读取相关参数，没有时使用默认值
	var (
		// 服务地址和服务名
		servicePort = flag.Int("service.port", 10086, "service port")
		serviceHost = flag.String("service.host", "192.168.1.108", "service host")
		serviceName = flag.String("service.name", "HelloConsul", "service name")
		// consul 地址
		consulPort = flag.Int("consul.port", 8500, "consul port")
		consulHost = flag.String("consul.host", "127.0.0.1", "consul host")
	)

	flag.Parse()

	ctx := context.Background()
	errChan := make(chan error)

	// 声明服务发现客户端
	var discoveryClient discover.DiscoveryClient

	discoveryClient, err := discover.NewKitDiscoverClient(*consulHost, *consulPort)
	// 获取服务发现客户端失败，直接关闭服务
	if err != nil{
		config.Logger.Println("Get Consul Client failed")
		os.Exit(-1)
	}

	// 声明并初始化 Service
	var svc = service.NewDiscoveryServiceImpl(discoveryClient)

	// 创建打招呼的Endpoint
	sayHelloEndpoint := endpoint.MakeHelloEndpoint(svc)
	// 创建服务发现的Endpoint
	discoveryEndpoint := endpoint.MakeDiscoveryEndpoint(svc)
	//创建健康检查的Endpoint
	healthEndpoint := endpoint.MakeHealthCheckEndpoint(svc)

	ep := endpoint.DiscoveryEndpoints{
		HelloEndpoint:		sayHelloEndpoint,
		DiscoveryEndpoint:		discoveryEndpoint,
		HealthCheckEndpoint:	healthEndpoint,
	}

	//创建http.Handler
	r := transport.MakeHttpHandler(ctx, ep, config.KitLogger)
	// 定义服务实例ID
	instanceId := *serviceName + "-" + uuid.NewV4().String()
	// 启动 http server
	go func() {
		config.Logger.Println("Http Server start at port:" + strconv.Itoa(*servicePort))
		//启动前执行注册
		if err :=discoveryClient.Register(*serviceName, instanceId, "/health", *serviceHost,  *servicePort, nil);err!=nil{
			config.Logger.Printf("string-service for service %s failed.", serviceName)
			// 注册失败，服务启动失败
			os.Exit(-1)
		}
		handler := r
		errChan <- http.ListenAndServe(":"  + strconv.Itoa(*servicePort), handler)
	}()

	go func() {
		// 监控系统信号，等待 ctrl + c 系统信号通知服务关闭
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	error := <-errChan
	//服务退出取消注册
	discoveryClient.DeRegister(instanceId)
	config.Logger.Println(error)
}
