package service

import (
	"context"
	"errors"
	"github.com/baxiang/soldiers-sortie/consul-discovery/discover"
)

type Service interface {

	// sayHelloService
	Hello() string
	// HealthCheck check service health status
	HealthCheck() bool


	//  discovery service from consul by serviceName
	DiscoveryService(ctx context.Context, serviceName string) ([]interface{}, error)

}

// 服务实例不存在
var ErrNotServiceInstances = errors.New("instances are not existed")

type DiscoveryServiceImpl struct {
	discoveryClient discover.DiscoveryClient
}

func NewDiscoveryServiceImpl(discoveryClient discover.DiscoveryClient) Service  {
	return &DiscoveryServiceImpl{
		discoveryClient:discoveryClient,
	}
}

func (*DiscoveryServiceImpl) Hello() string {
	return "Hello World!"
}

func (service *DiscoveryServiceImpl) DiscoveryService(_ context.Context, serviceName string) ([]interface{}, error)  {
    // 获取服务实例
	instances,err:= service.discoveryClient.DiscoverServices(serviceName)

	if err!= nil || len(instances) == 0 {
		return nil, ErrNotServiceInstances
	}
	return instances, nil
}


// HealthCheck implement Service method
// 用于检查服务的健康状态，这里仅仅返回true
func (*DiscoveryServiceImpl) HealthCheck() bool {
	return true
}

