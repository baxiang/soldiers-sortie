package plugin

import "github.com/opentracing/opentracing-go"

type Options struct {
	SvrAddr string
	Services []string
	SelectorSvrAddr string
	TracingSvrAddr string
}
type Option func(*Options)
func WithSvrAddr(addr string) Option {
	return func(o *Options) {
		o.SvrAddr = addr
	}
}

// WithSvrAddr allows you to set Services of Options
func WithServices(services []string) Option {
	return func(o *Options) {
		o.Services = services
	}
}

// WithSvrAddr allows you to set SelectorSvrAddr of Options
func WithSelectorSvrAddr(addr string) Option {
	return func(o *Options) {
		o.SelectorSvrAddr = addr
	}
}

// WithSvrAddr allows you to set TracingSvrAddr of Options
func WithTracingSvrAddr(addr string) Option {
	return func(o *Options) {
		o.TracingSvrAddr = addr
	}
}




type Plugin interface {

}

type ResolverPlugin interface {
	Init(... Option)error
}

type TracingPlugin interface {
	Init(...Option) (opentracing.Tracer, error)
}

var PluginMap = make(map[string]Plugin)

// Register opens an entry point for all plug-ins to register
func Register(name string, plugin Plugin) {
	if PluginMap == nil {
		PluginMap = make(map[string]Plugin)
	}
	PluginMap[name] = plugin
}
