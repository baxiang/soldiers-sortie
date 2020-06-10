package transport

import (
	"context"
	"encoding/binary"
	"github.com/baxiang/gorpc/codec"
	"github.com/baxiang/gorpc/codes"
	"io"
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



func NewFramer() Framer{
	return &framer{
		buffer: make([]byte,DefaultPayloadLength),
	}
}
type framer struct {
	buffer []byte
	counter int
}

func(f *framer)Resize(){
	f.buffer = make([]byte,len(f.buffer)*2)
}

func(f *framer)ReadFrame(conn net.Conn)([]byte,error){
	frameHeader := make([]byte, codec.FrameHeadLen)
	if num, err := io.ReadFull(conn, frameHeader); num != codec.FrameHeadLen || err != nil {
		return nil, err
	}

	// validate magic
	if magic := uint8(frameHeader[0]); magic != codec.Magic {
		return nil, codes.NewFrameworkError(codes.ClientMsgErrorCode, "invalid magic...")
	}

	length := binary.BigEndian.Uint32(frameHeader[7:11])

	if length > MaxPayloadLength {
		return nil, codes.NewFrameworkError(codes.ClientMsgErrorCode, "payload too large...")
	}

	for uint32(len(f.buffer)) < length && f.counter <= 12 {
		f.buffer = make([]byte, len(f.buffer)*2)
		f.counter++
	}

	if num, err := io.ReadFull(conn, f.buffer[:length]); uint32(num) != length || err != nil {
		return nil, err
	}

	return append(frameHeader, f.buffer[:length]...), nil
}