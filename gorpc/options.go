package gorpc

import "time"

type ServerOptions struct {
	address string
	network string
	protocol string
	timeout time.Duration
	serializationType string
	selectorSvrAddr string
	tracingSvrAddr string
	pluginNames []string
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
