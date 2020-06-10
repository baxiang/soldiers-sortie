package gorpc

import (
	"github.com/baxiang/gorpc/interceptor"
	"time"
)

// 服务配置模块
type ServerOptions struct {
	address string // listening address, e.g. :( ip://127.0.0.1:8080、 dns://www.google.com)
	network string  // network type, e.g. : tcp、udp 传输协议
	protocol string // protocol typpe, e.g. : proto、json 文件传输格式
	timeout time.Duration

	serializationType string
	selectorSvrAddr string
	tracingSvrAddr string
	tracingSpanName string
	pluginNames []string
	interceptors []interceptor.ServerInterceptor
}


type ServerOption func(*ServerOptions)

func WithAddress(address string)ServerOption{
	return func(o *ServerOptions) {
		 o.address = address
	}
}

func WithNetwork(network string)ServerOption{
	return func(o *ServerOptions) {
		o.network = network
	}
}

func WithProtocol(protocol string) ServerOption {
	return func(o *ServerOptions) {
		o.protocol = protocol
	}
}

func WithTimeout(timeout time.Duration) ServerOption {
	return func(o *ServerOptions) {
		o.timeout = timeout
	}
}

func WithSerializationType(serializationType string) ServerOption {
	return func(o *ServerOptions) {
		o.serializationType = serializationType
	}
}

func WithSelectorSvrAddr(addr string) ServerOption {
	return func(o *ServerOptions) {
		o.selectorSvrAddr = addr
	}
}

func WithPlugin(pluginName ... string) ServerOption {
	return func(o *ServerOptions) {
		o.pluginNames = append(o.pluginNames, pluginName ...)
	}
}


func WithInterceptor(interceptors ...interceptor.ServerInterceptor) ServerOption {
	return func(o *ServerOptions) {
		o.interceptors = append(o.interceptors, interceptors...)
	}
}

func WithTracingSvrAddr(addr string) ServerOption {
	return func(o *ServerOptions) {
		o.tracingSvrAddr = addr
	}
}

func WithTracingSpanName(name string) ServerOption {
	return func(o *ServerOptions) {
		o.tracingSpanName = name
	}
}