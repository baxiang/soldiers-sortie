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

type ProcessDefinitionEndpoints struct {
	GetEndpoint   endpoint.Endpoint
	StartEndpoint endpoint.Endpoint
}

func NewProcessDefinitionClient(conn *grpc.ClientConn, tracer kitgrpc.ClientOption) *ProcessDefinitionEndpoints {
	endpoints := new(ProcessDefinitionEndpoints)

	options := make([]kitgrpc.ClientOption, 0)
	if tracer != nil {
		options = append(options, tracer)
	}

	{
		e := grpctransport.NewClient(
			conn,
			"pb.ProcessDefinition",
			"Get",
			encodeRequest,
			decodeResponse,
			pb.GetProcessDefinitionResp{},
			options...,
		).Endpoint()
		endpoints.GetEndpoint = e
	}

	{
		e := grpctransport.NewClient(
			conn,
			"pb.ProcessDefinition",
			"Start",
			encodeRequest,
			decodeResponse,
			pb.StartProcessDefinitionResp{},
			options...,
		).Endpoint()
		endpoints.StartEndpoint = e
	}

	return endpoints
}

func (p *ProcessDefinitionEndpoints) Get(ctx context.Context, in *pb.GetProcessDefinitionReq) (*pb.GetProcessDefinitionResp, error) {
	res, err := p.GetEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	reply, ok := res.(*pb.GetProcessDefinitionResp)
	if !ok {
		return nil, kit.ErrReplyTypeInvalid
	}
	return reply, nil
}
func (p *ProcessDefinitionEndpoints) Start(ctx context.Context, in *pb.StartProcessDefinitionReq) (*pb.StartProcessDefinitionResp, error) {
	res, err := p.StartEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	reply, ok := res.(*pb.StartProcessDefinitionResp)
	if !ok {
		return nil, kit.ErrReplyTypeInvalid
	}
	return reply, nil
}
