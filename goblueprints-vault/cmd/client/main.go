package main

import (
	"flag"
	"fmt"
	"context"
	"github.com/baxiang/soldiers-sortie/goblueprints-vault/endpoints"
	"github.com/baxiang/soldiers-sortie/goblueprints-vault/pb"
	"github.com/baxiang/soldiers-sortie/goblueprints-vault/services"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
	"os"
	"time"
	"log"
)

func main() {
	var (
		grpcAddr = flag.String("addr", ":8081", "gRPC address")
	)
	flag.Parse()
	ctx := context.Background()
	conn, err := grpc.Dial(*grpcAddr, grpc.WithInsecure(), grpc.WithTimeout(1*time.Second))
	if err != nil {
		log.Fatalln("gRPC dial:", err)
	}
	defer conn.Close()
	vaultService :=New(conn)
	args := flag.Args()
	var cmd string
	cmd, args = pop(args)
	switch cmd {
	case "hash":
		var password string
		password, args = pop(args)
		hash(ctx, vaultService, password)
	case "validate":
		var password, hash string
		password, args = pop(args)
		hash, args = pop(args)
		validate(ctx, vaultService, password, hash)
	default:
		log.Fatalln("unknown command", cmd)
	}
}

// New makes a new vault.Service client.
func New(conn *grpc.ClientConn) services.Service {
	var hashEndpoint = grpctransport.NewClient(
		conn, "pb.Vault", "Hash",
		EncodeGRPCHashRequest,
		DecodeGRPCHashResponse,
		pb.HashResponse{},
	).Endpoint()
	var validateEndpoint = grpctransport.NewClient(
		conn, "pb.Vault", "Validate",
		EncodeGRPCValidateRequest,
		DecodeGRPCValidateResponse,
		pb.ValidateResponse{},
	).Endpoint()
	return endpoints.VaultEndpoints{
		HashEndpoint:     hashEndpoint,
		ValidateEndpoint: validateEndpoint,
	}
}
func EncodeGRPCHashRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(endpoints.HashRequest)
	return &pb.HashRequest{Password: req.Password}, nil
}

func DecodeGRPCHashResponse(ctx context.Context, r interface{}) (interface{}, error) {
	res := r.(*pb.HashResponse)
	return endpoints.HashResponse{Hash: res.Hash, Err: res.Err}, nil
}

func EncodeGRPCValidateRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(endpoints.ValidateRequest)
	return &pb.ValidateRequest{Password: req.Password, Hash: req.Hash}, nil
}
func DecodeGRPCValidateResponse(ctx context.Context, r interface{}) (interface{}, error) {
	res := r.(*pb.ValidateResponse)
	return endpoints.ValidateResponse{Valid: res.Valid}, nil
}

func hash(ctx context.Context, service services.Service, password string) {
	h, err := service.Hash(ctx, password)
	if err != nil {
		fmt.Printf("hash error  %v",err)
		//log.Fatalln(err.Error())
	}
	fmt.Println(h)
}

func validate(ctx context.Context, service services.Service, password, hash string) {
	valid, err := service.Validate(ctx, password, hash)
	if err != nil {
		log.Fatalln(err.Error())
	}
	if !valid {
		fmt.Println("invalid")
		os.Exit(1)
	}
	fmt.Println("valid")
}

func pop(s []string) (string, []string) {
	if len(s) == 0 {
		return "", s
	}
	return s[0], s[1:]
}