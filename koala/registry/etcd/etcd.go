package etcd

import (
	"go.etcd.io/etcd/clientv3"
	"github.com/baxiang/koala/registry"
	"sync"
	"sync/atomic"
	"time"
)

const (
	MaxServiceNum = 8
	MaxSyncServiceInterval = time.Second *10
)

type EtcdRegistry struct {
	options   *registry.Options
	client    *clientv3.Client
	serviceCh chan *registry.Service

	value              atomic.Value
	lock               sync.Mutex
	registryServiceMap map[string]*RegisterService
}

type AllServiceInfo struct {
	serviceMap map[string]*registry.Service
}