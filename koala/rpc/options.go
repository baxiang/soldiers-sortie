package rpc

import "time"

type RpcOptions struct {
	ConnTimeout time.Duration
	WriteTimeout time.Duration
	ReadTimeout time.Duration
	ServiceName string
	ClientServiceName string

	RegisterName string //注册中心名称
	RegisterAddr string //注册中心地址
	RegisterPath string

	MaxLimitQPS int

	TraceReportAddr string
	TraceSampleType string
	//trace sample rate
	TraceSampleRate float64
}

type RpcOption func(opts *RpcOptions)


func WithLimitQPS(qps int) RpcOption {
	return func(opts *RpcOptions) {
		opts.MaxLimitQPS = qps
	}
}

func WithConnTimeout(timeout time.Duration) RpcOption {
	return func(opts *RpcOptions) {
		opts.ConnTimeout = timeout
	}
}

func WithWriteTimeout(timeout time.Duration) RpcOption {
	return func(opts *RpcOptions) {
		opts.WriteTimeout = timeout
	}
}

func WithReadTimeout(timeout time.Duration) RpcOption {
	return func(opts *RpcOptions) {
		opts.ReadTimeout = timeout
	}
}

func WithServiceName(serviceName string) RpcOption {
	return func(opts *RpcOptions) {
		opts.ServiceName = serviceName
	}
}

func WithRegisterName(name string) RpcOption {
	return func(opts *RpcOptions) {
		opts.RegisterName = name
	}
}

func WithRegisterAddr(addr string) RpcOption {
	return func(opts *RpcOptions) {
		opts.RegisterAddr = addr
	}
}

func WithRegisterPath(path string) RpcOption {
	return func(opts *RpcOptions) {
		opts.RegisterPath = path
	}
}

func WithTraceReportAddr(addr string) RpcOption {
	return func(opts *RpcOptions) {
		opts.TraceReportAddr = addr
	}
}

func WithTraceSampleType(stype string) RpcOption {
	return func(opts *RpcOptions) {
		opts.TraceSampleType = stype
	}
}

func WithTraceSampleRate(rate float64) RpcOption {
	return func(opts *RpcOptions) {
		opts.TraceSampleRate = rate
	}
}

func WithClientServiceName(name string) RpcOption {
	return func(opts *RpcOptions) {
		opts.ClientServiceName = name
	}
}