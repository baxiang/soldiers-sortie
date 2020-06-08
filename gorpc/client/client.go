package client

import (
	"context"
	"github.com/baxiang/gorpc/codec"
	"github.com/baxiang/gorpc/interceptor"
)

type Client interface {
	Invoke(ctx context.Context,req, rsp interface{},path string,opts ...Option)error
}

// use a global client
var DefaultClient = New()
var New = func() *defaultClient {
	return &defaultClient{
		opts : &Options{
			protocol : "proto",
		},
	}
}

type defaultClient struct {
	opts *Options
}

func (c *defaultClient)Call(ctx context.Context, servicePath string, req interface{}, rsp interface{},
	opts ...Option) error {
	callOpts := make([]Option, 0, len(opts)+1)
	callOpts = append(callOpts,opts...)
	callOpts = append(callOpts,WithSerializationType(codec.MsgPack))


}

func(c *defaultClient)Invoke(ctx context.Context,req, rsp interface{},path string,opts...Option)error{
	for _, o := range opts {
		o(c.opts)
	}

	if c.opts.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.opts.timeout)
		defer cancel()
	}

	// set serviceName, method
	newCtx, clientStream := stream.NewClientStream(ctx)

	serviceName, method , err := utils.ParseServicePath(path)
	if err != nil {
		return err
	}

	c.opts.serviceName = serviceName
	c.opts.method = method

	// TODO : delete or not
	clientStream.WithServiceName(serviceName)
	clientStream.WithMethod(method)

	// execute the interceptor first
	return interceptor.ClientIntercept(newCtx, req, rsp, c.opts.interceptors, c.invoke)
}