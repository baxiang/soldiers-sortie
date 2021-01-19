package gateway

import (
	"context"
	"userService/pkg/apstfr/apstfrpb"

	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

type ScanEndpoints struct {
	PayEndpoint endpoint.Endpoint
}

func (s *ScanEndpoints) Pay(ctx context.Context, in *apstfrpb.PayRequest) (*apstfrpb.PayReply, error) {
	res, err := s.PayEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*apstfrpb.PayReply), nil
}

func NewScanGRPCClient(conn *grpc.ClientConn, tracer kitgrpc.ClientOption) *ScanEndpoints {
	endpoints := new(ScanEndpoints)
	options := make([]kitgrpc.ClientOption, 0)
	if tracer != nil {
		options = append(options, tracer)
	}

	{
		e := grpctransport.NewClient(
			conn,
			"apstfrpb.Scan",
			"Pay",
			encodeRequest,
			decodeResponse,
			apstfrpb.PayReply{},
			options...,
		).Endpoint()
		endpoints.PayEndpoint = e
	}

	return endpoints
}
