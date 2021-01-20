package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/baxiang/soldiers-sortie/lorem-grpc/endpoints"
	"github.com/baxiang/soldiers-sortie/lorem-grpc/transport"
	"log"
	"strconv"
	"time"

	"github.com/baxiang/soldiers-sortie/lorem-grpc/pb"
	"github.com/baxiang/soldiers-sortie/lorem-grpc/services"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

func main() {
	var (
		grpcAddr = flag.String("addr", ":8081",
			"gRPC address")
	)
	flag.Parse()
	ctx := context.Background()
	conn, err := grpc.Dial(*grpcAddr, grpc.WithInsecure(),
		grpc.WithTimeout(1*time.Second))

	if err != nil {
		log.Fatalln("gRPC dial:", err)
	}
	defer conn.Close()

	loremService := NewGrpcClient(conn)
	args := flag.Args()
	var cmd string
	cmd, args = pop(args)

	switch cmd {
	case "lorem":
		var requestType, minStr, maxStr string

		requestType, args = pop(args)
		minStr, args = pop(args)
		maxStr, args = pop(args)

		min, _ := strconv.Atoi(minStr)
		max, _ := strconv.Atoi(maxStr)
		lorem(ctx, loremService, requestType, min, max)
	default:
		log.Fatalln("unknown command", cmd)
	}
}

// Return new lorem_grpc service
func NewGrpcClient(conn *grpc.ClientConn) services.Service {
	var loremEndpoint = grpctransport.NewClient(
		conn, "Lorem", "Lorem",
		transport.EncodeGRPCLoremRequest,
		transport.DecodeGRPCLoremResponse,
		pb.LoremResponse{},
	).Endpoint()

	return endpoints.Endpoints{
		LoremEndpoint: loremEndpoint,
	}
}

// parse command line argument one by one
func pop(s []string) (string, []string) {
	if len(s) == 0 {
		return "", s
	}
	return s[0], s[1:]
}

// call lorem service
func lorem(ctx context.Context, service services.Service, requestType string, min int, max int) {
	message, err := service.Lorem(ctx, requestType, min, max)
	if err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Println(message)
}
