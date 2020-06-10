package gorpc

import (
	"context"
	"fmt"
	"errors"
	"github.com/baxiang/gorpc/codec"
	"github.com/baxiang/gorpc/codes"
	"github.com/baxiang/gorpc/interceptor"
	"github.com/baxiang/gorpc/metadata"
	"github.com/baxiang/gorpc/protocol"
	"github.com/baxiang/gorpc/transport"
	"github.com/baxiang/gorpc/utils"
	"github.com/baxiang/gorpc/log"
	"github.com/golang/protobuf/proto"
)



//服务通用接口
type Service interface {
	Register(string,Handler)
	Serve(*ServerOptions)
	Close()
}

type service struct {
	svr interface{}
	ctx context.Context //上下文
	cancel context.CancelFunc
	serviceName string
	handlers map[string]Handler
	opts *ServerOptions
	closing bool
}
type ServiceDesc struct {
	Svr interface{}
	ServiceName string
	Methods []*MethodDesc
	HandlerType interface{}
}

type MethodDesc struct {
	MethodName string
	Handler Handler
}

type Handler func (interface{}, context.Context, func(interface{}) error, []interceptor.ServerInterceptor) (interface{}, error)
func(s *service)Register(handleName string,handler Handler){
	if s.handlers==nil{
		s.handlers = make(map[string]Handler)
	}
	s.handlers[handleName]= handler
}

func (s *service)Serve(opts *ServerOptions){
	s.opts = opts
	transportOpts := []transport.ServerTransportOption {
		transport.WithServerAddress(s.opts.address),
		transport.WithServerNetwork(s.opts.network),
		transport.WithHandler(s),
		transport.WithServerTimeout(s.opts.timeout),
		transport.WithSerializationType(s.opts.serializationType),
		transport.WithProtocol(s.opts.protocol),
	}

	serverTransport := transport.GetServerTransport(s.opts.protocol)

	s.ctx, s.cancel = context.WithCancel(context.Background())

	if err := serverTransport.ListenAndServe(s.ctx, transportOpts ...); err != nil {
		log.Errorf("%s serve error, %v", s.opts.network, err)
		return
	}

	fmt.Printf("%s service serving at %s ... \n",s.opts.protocol, s.opts.address)

	<- s.ctx.Done()
}


func(s *service)Close(){
	s.closing = true
	if s.cancel !=nil{
		s.cancel()
	}
	fmt.Println("service closing")
}

func(s *service)Handle(ctx context.Context,reqbuf []byte)([]byte,error){
	request :=&protocol.Request{}
	if err := proto.Unmarshal(reqbuf, request); err != nil {
		return nil, err
	}
	ctx = metadata.WithServerMetadata(ctx, request.Metadata)

	serverSerialization := codec.GetSerialization(s.opts.serializationType)

	dec := func(req interface {}) error {

		if err := serverSerialization.Unmarshal(request.Payload, req); err != nil {
			return err
		}
		return nil
	}

	if s.opts.timeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, s.opts.timeout)
		defer cancel()
	}


	_, method , err := utils.ParseServicePath(string(request.ServicePath))
	if err != nil {
		return nil, codes.New(codes.ClientMsgErrorCode, "method is invalid")
	}

	handler := s.handlers[method]
	if handler == nil {
		return nil, errors.New("handlers is nil")
	}

	rsp, err := handler(s.svr, ctx, dec, s.opts.interceptors)
	if err != nil {
		return nil, err
	}

	rspbuf, err := serverSerialization.Marshal(rsp)
	if err != nil {
		return nil, err
	}

	return rspbuf, nil

}