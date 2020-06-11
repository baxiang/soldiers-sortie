package registry

import "time"

type Options struct {
	Addrs []string
	Timeout time.Duration
	RegistryPath string
	HeartBeat int64
}


type Option func (options *Options)

func WithAddrs(addrs []string)Option{
	return func(options *Options) {
		options.Addrs = addrs
	}
}

func WithTimeout(timeout time.Duration)Option{
	return func(options *Options) {
		options.Timeout = timeout
	}
}

func WithRegistryPath(registryPath string)Option{
	return func(options *Options) {
		options.RegistryPath = registryPath
	}
}

func WithHeartBeat(heatBeat int64)Option{
	return func(options *Options) {
		options.HeartBeat = heatBeat
	}
}


