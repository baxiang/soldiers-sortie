package main

import (
	pb "github.com/baxiang/soldiers-sortie/go-grpc-example/proto"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
	"time"
)

var (
	Address =":8000"
)

type StreamService struct {
	pb.UnimplementedStreamServerServer
}

// ListValue 实现ListValue方法
func (s *StreamService) ListValue(req *pb.SimpleReq, srv pb.StreamServer_ListValueServer) error {
	for n := 0; n < 5; n++ {
		// 向流中发送消息， 默认每次send送消息最大长度为`math.MaxInt32`bytes
		err := srv.Send(&pb.StreamResp{
			Val: req.Para + strconv.Itoa(n),
		})
		if err != nil {
			return err
		}
		// log.Println(n)
		time.Sleep(1 * time.Second)
	}
	return nil
}
func main() {
	listener, err := net.Listen("tcp", Address)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}
	log.Println(Address + " net.Listing...")
	// 新建gRPC服务器实例
	// 默认单次接收最大消息长度为`1024*1024*4`bytes(4M)，单次发送消息最大长度为`math.MaxInt32`bytes
	// grpcServer := grpc.NewServer(grpc.MaxRecvMsgSize(1024*1024*4), grpc.MaxSendMsgSize(math.MaxInt32))
	grpcServer := grpc.NewServer()
	// 在gRPC服务器注册我们的服务
	pb.RegisterStreamServerServer(grpcServer, &StreamService{})

	//用服务器 Serve() 方法以及我们的端口信息区实现阻塞等待，直到进程被杀死或者 Stop() 被调用
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("grpcServer.Serve err: %v", err)
	}
}
