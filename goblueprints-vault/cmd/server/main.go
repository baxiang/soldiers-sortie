package main

import (
	"flag"
	"fmt"
	"context"
	ep "github.com/baxiang/soldiers-sortie/goblueprints-vault/endpoints"
	"github.com/baxiang/soldiers-sortie/goblueprints-vault/pb"
	"github.com/baxiang/soldiers-sortie/goblueprints-vault/services"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	httptransport "github.com/go-kit/kit/transport/http"
)

func main() {
	var (
		httpAddr = flag.String("http", ":8080", "http listen address")
		gRPCAddr = flag.String("grpc", ":8081", "gRPC listen address")
	)
	flag.Parse()
	ctx := context.Background()
	srv := services.NewService()
	errChan := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	hashEndpoint := ep.MakeHashEndpoint(srv)
	validateEndpoint := ep.MakeValidateEndpoint(srv)
	endpoints := ep.VaultEndpoints{
		HashEndpoint:     hashEndpoint,
		ValidateEndpoint: validateEndpoint,
	}

	// HTTP transport
	go func() {
		log.Println("http:", *httpAddr)
		handler := NewHTTPServer(ctx, endpoints)
		errChan <- http.ListenAndServe(*httpAddr, handler)
	}()

	// gRPC transport
	go func() {
		listener, err := net.Listen("tcp", *gRPCAddr)
		if err != nil {
			errChan <- err
			return
		}
		log.Println("grpc:", *gRPCAddr)
		handler := NewGRPCServer(ctx, endpoints)
		gRPCServer := grpc.NewServer()
		pb.RegisterVaultServer(gRPCServer, handler)
		errChan <- gRPCServer.Serve(listener)
	}()

	log.Fatalln(<-errChan)
}
// NewHTTPServer makes a new Vault HTTP service.
func NewHTTPServer(ctx context.Context, endpoints ep.VaultEndpoints) http.Handler {
	m := http.NewServeMux()
	m.Handle("/hash", httptransport.NewServer(
		endpoints.HashEndpoint,
		ep.DecodeHashRequest,
		ep.EncodeResponse,
	))
	m.Handle("/validate", httptransport.NewServer(
		endpoints.ValidateEndpoint,
		ep.DecodeValidateRequest,
		ep.EncodeResponse,
	))
	return m
}

// NewGRPCServer gets a new pb.VaultServer.
func NewGRPCServer(ctx context.Context, endpoints ep.VaultEndpoints) pb.VaultServer {
	return &grpcServer{
		hash: grpctransport.NewServer(
			endpoints.HashEndpoint,
			DecodeGRPCHashRequest,
			EncodeGRPCHashResponse,
		),
		validate: grpctransport.NewServer(
			endpoints.ValidateEndpoint,
			DecodeGRPCValidateRequest,
			EncodeGRPCValidateResponse,
		),
	}
}

type grpcServer struct {
	hash     grpctransport.Handler
	validate grpctransport.Handler
}

func (s *grpcServer) Hash(ctx context.Context, r *pb.HashRequest) (*pb.HashResponse, error) {
	_, resp, err := s.hash.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.HashResponse), nil
}

func (s *grpcServer) Validate(ctx context.Context, r *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	_, resp, err := s.validate.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.ValidateResponse), nil
}



func DecodeGRPCHashRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.HashRequest)
	return ep.HashRequest{Password: req.Password}, nil
}

func EncodeGRPCHashResponse(ctx context.Context, r interface{}) (interface{}, error) {
	res := r.(ep.HashResponse)
	return &pb.HashResponse{Hash: res.Hash, Err: res.Err}, nil
}



func DecodeGRPCValidateRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.ValidateRequest)
	return ep.ValidateRequest{Password: req.Password, Hash: req.Hash}, nil
}

func EncodeGRPCValidateResponse(ctx context.Context, r interface{}) (interface{}, error) {
	res := r.(ep.ValidateResponse)
	return &pb.ValidateResponse{Valid: res.Valid}, nil
}