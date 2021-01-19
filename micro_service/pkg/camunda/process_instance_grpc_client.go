package camunda

import (
	"context"
	"userService/pkg/camunda/pb"
	"userService/pkg/kit"

	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

type ProcessInstanceEndpoints struct {
	ListEndpoint endpoint.Endpoint
}

func NewProcessInstanceClient(conn *grpc.ClientConn, tracer kitgrpc.ClientOption) *ProcessInstanceEndpoints {
	endpoints := new(ProcessInstanceEndpoints)

	options := make([]kitgrpc.ClientOption, 0)
	if tracer != nil {
		options = append(options, tracer)
	}

	{
		e := grpctransport.NewClient(
			conn,
			"pb.ProcessInstance",
			"List",
			encodeRequest,
			decodeResponse,
			pb.ProcessInstanceListResp{},
			options...,
		).Endpoint()
		endpoints.ListEndpoint = e
	}

	return endpoints
}

func (p *ProcessInstanceEndpoints) List(ctx context.Context, in *pb.ProcessInstanceListReq) (*pb.ProcessInstanceListResp, error) {
	res, err := p.ListEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	reply, ok := res.(*pb.ProcessInstanceListResp)
	if !ok {
		return nil, kit.ErrReplyTypeInvalid
	}
	return reply, nil
}
