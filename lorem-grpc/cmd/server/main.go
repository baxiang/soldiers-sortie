package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/baxiang/soldiers-sortie/lorem-grpc/endpoints"
	"github.com/baxiang/soldiers-sortie/lorem-grpc/pb"
	"github.com/baxiang/soldiers-sortie/lorem-grpc/services"
	"github.com/baxiang/soldiers-sortie/lorem-grpc/transport"
	"google.golang.org/grpc"
)

func main() {
	var gRPCAddr = flag.String("grpc", ":8081", "gRPC listen address")
	flag.Parse()

	var errChan chan error
	go func() {
		listener, err := net.Listen("tcp", *gRPCAddr)
		if err != nil {
			errChan <- err
			return
		}
		gRPCServer := grpc.NewServer()
		svc := services.LoremService{}
		e := endpoints.Endpoints{LoremEndpoint: endpoints.MakeLoremEndpoint(svc)}
		handler := transport.NewGRPCServer(context.Background(), e)
		pb.RegisterLoremServer(gRPCServer, handler)
		errChan <- gRPCServer.Serve(listener)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()
	fmt.Println(<-errChan)
}
