package client

import (
	"github.com/baxiang/gorpc/auth"
	"github.com/baxiang/gorpc/interceptor"
	"github.com/baxiang/gorpc/transport"
	"time"
)

type Options struct {
	serviceName string // 服务名称
	method string // 方法
	target string // ip:port 127.0.0.1:8000
	timeout time.Duration // 超时
	network string // 网络形式 tcp,udp
	protocol string
	serializationType string
	transportOpts transport.ClientTransportOptions
	interceptors []interceptor.ClientInterceptor
	selectorName string // 服务发现方式
	perRPCAuth []auth.PerRPCAuth
	transportAuth auth.TransportAuth
}

type Option func(*Options)

func WithServiceName(serviceName string) Option {
	return func(o *Options) {
		o.serviceName = serviceName
	}
}

func WithMethod(method string) Option {
	return func(o *Options) {
		o.method = method
	}
}

func WithTarget(target string) Option {
	return func(o *Options) {
		o.target = target
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.timeout = timeout
	}
}

func WithNetwork(network string) Option {
	return func(o *Options) {
		o.network = network
	}
}

func WithProtocol(protocol string) Option {
	return func(o *Options) {
		o.protocol = protocol
	}
}
func WithSerializationType(serializationType string) Option {
	return func(o *Options) {
		o.serializationType = serializationType
	}
}
func WithSelectorName(selectorName string) Option {
	return func(o *Options) {
		o.selectorName = selectorName
	}
}

func WithInterceptor(interceptors ...interceptor.ClientInterceptor) Option {
	return func(o *Options) {
		o.interceptors = append(o.interceptors, interceptors...)
	}
}

func WithPerRPCAuth(rpcAuth auth.PerRPCAuth) Option {
	return func(o *Options) {
		o.perRPCAuth = append(o.perRPCAuth,rpcAuth)
	}
}

func WithTransportAuth(transportAuth auth.TransportAuth) Option {
	return func(o *Options) {
		o.transportAuth = transportAuth
	}
}

