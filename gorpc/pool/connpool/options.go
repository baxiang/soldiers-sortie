package connpool

import "time"

type Options struct {
	initialCap int // 初始化连接数
	maxCap int // 最大链接数
	maxIdle int // 最大空闲链接数
	idleTimeout time.Duration
	dialTimeout time.Duration
}

type Option func(*Options)

func WithInitialCap (initialCap int) Option {
	return func(o *Options) {
		o.initialCap = initialCap
	}
}

func WithMaxCap (maxCap int) Option {
	return func(o *Options) {
		o.maxCap = maxCap
	}
}


func WithMaxIdle (maxIdle int) Option {
	return func(o *Options) {
		o.maxIdle = maxIdle
	}
}

func WithIdleTimeout(idleTimeout time.Duration) Option {
	return func(o *Options) {
		o.idleTimeout = idleTimeout
	}
}

func WithDialTimeout(dialTimeout time.Duration) Option {
	return func(o *Options) {
		o.dialTimeout = dialTimeout
	}
}
