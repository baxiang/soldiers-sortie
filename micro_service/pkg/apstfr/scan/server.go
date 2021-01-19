package scan

import (
	"context"
	"userService/pkg/apstfr/apstfrpb"

	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type server struct{}

func (s server) Pay(ctx context.Context, in *apstfrpb.PayRequest) (*apstfrpb.PayReply, error) {
	panic("implement me")
}

func New(tracer grpctransport.ServerOption) apstfrpb.ScanServer {
	return &server{}
}
