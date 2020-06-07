package gorpc

import (
	"context"
	"fmt"
	"github.com/baxiang/gorpc/interceptor"
)

//服务
type Handler func(interface{},context.Context,func(interface{})error,[]interceptor.ServerInterceptor)(interface{},error)

//服务通用接口
type Service interface {
	Register(string,Handler)
	Serve(*ServerOptions)
	Close()
}

type service struct {
	serviceName string
	svr interface{}
	ctx context.Context //上下文
	cancel context.CancelFunc
	handlers map[string]Handler
	opts *ServerOption
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


func(s *service)Register(handleName string,handler Handler){
	if s.handlers==nil{
		s.handlers = make(map[string]Handler)
	}
	s.handlers[handleName]= handler
}

func (s *service)Serve(opts *ServerOptions){

}


func(s *service)Close(){
	s.closing = true
	if s.cancel !=nil{
		s.cancel()
	}
	fmt.Println("service closing")
}