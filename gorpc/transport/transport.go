package transport

import (
	"context"
	"net"
)

const DefaultPayloadLength  = 1024
const MaxPayloadLength  =  4*1024*1014

// 服务传输层主要提供一种监听和处理请求的能力
type ServerTransport interface {
	// monitoring and processing of requests
	ListenAndServe(context.Context, ...ServerTransportOption) error
}


// 客户端传输层主要提供一种向下游发送请求的能力
type ClientTransport interface {
	Send(context.Context,[]byte,...ClientTransportOption)([]byte,error)
}

type Framer interface {
	ReadFrame(net.Conn)([]byte,error)
}