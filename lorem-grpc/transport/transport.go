package transport

import (
	"github.com/baxiang/soldiers-sortie/lorem-grpc/endpoints"
	"golang.org/x/net/context"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/baxiang/soldiers-sortie/lorem-grpc/pb"
)

type grpcServer struct {
	lorem grpctransport.Handler
}

// implement LoremServer Interface in lorem.pb.go
func (s *grpcServer) Lorem(ctx context.Context, r *pb.LoremRequest) (*pb.LoremResponse, error) {
	_, resp, err := s.lorem.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.LoremResponse), nil
}

// create new grpc server
func NewGRPCServer(_ context.Context, endpoint endpoints.Endpoints) pb.LoremServer {
	return &grpcServer{
		lorem: grpctransport.NewServer(
			endpoint.LoremEndpoint,
			DecodeGRPCLoremRequest,
			EncodeGRPCLoremResponse,
		),
	}
}

//Encode and Decode Lorem Request
func EncodeGRPCLoremRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(endpoints.LoremRequest)
	return &pb.LoremRequest{
		RequestType: req.RequestType,
		Max: req.Max,
		Min: req.Min,
	} , nil
}

func DecodeGRPCLoremRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.LoremRequest)
	return endpoints.LoremRequest{
		RequestType: req.RequestType,
		Max: req.Max,
		Min: req.Min,
	}, nil
}

// Encode and Decode Lorem Response
func EncodeGRPCLoremResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(endpoints.LoremResponse)
	return &pb.LoremResponse{
		Message: resp.Message,
		Err: resp.Err,
	}, nil
}

func DecodeGRPCLoremResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pb.LoremResponse)
	return endpoints.LoremResponse{
		Message: resp.Message,
		Err: resp.Err,
	}, nil
}