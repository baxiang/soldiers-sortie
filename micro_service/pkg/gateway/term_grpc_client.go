package gateway

import (
	"context"
	"userService/pkg/kit"
	"userService/pkg/pb"

	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

type TermEndpoints struct {
	ListTermInfoEndpoint            endpoint.Endpoint
	SaveTermEndpoint                endpoint.Endpoint
	ListTermRiskEndpoint            endpoint.Endpoint
	ListTermActivationStateEndpoint endpoint.Endpoint
	UpdateTermInfoEndpoint          endpoint.Endpoint
	QueryTermInfoEndpoint           endpoint.Endpoint
}

func (t *TermEndpoints) QueryTermInfo(ctx context.Context, in *pb.QueryTermInfoRequest) (*pb.QueryTermInfoReply, error) {
	res, err := t.QueryTermInfoEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.QueryTermInfoReply), nil
}

func (t *TermEndpoints) UpdateTermInfo(ctx context.Context, in *pb.UpdateTermInfoRequest) (*pb.UpdateTermInfoReply, error) {
	res, err := t.UpdateTermInfoEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.UpdateTermInfoReply), nil
}

func (t *TermEndpoints) ListTermActivationState(ctx context.Context, in *pb.ListTermActivationStateRequest) (*pb.ListTermActivationStateReply, error) {
	res, err := t.ListTermActivationStateEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.ListTermActivationStateReply), nil
}

func (t *TermEndpoints) ListTermRisk(ctx context.Context, in *pb.ListTermRiskRequest) (*pb.ListTermRiskReply, error) {
	res, err := t.ListTermRiskEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.ListTermRiskReply), nil
}

func (t *TermEndpoints) SaveTerm(ctx context.Context, in *pb.SaveTermRequest) (*pb.SaveTermReply, error) {
	res, err := t.SaveTermEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.SaveTermReply), nil
}

func (t *TermEndpoints) ListTermInfo(ctx context.Context, in *pb.ListTermInfoRequest) (*pb.ListTermInfoReply, error) {
	res, err := t.ListTermInfoEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	reply, ok := res.(*pb.ListTermInfoReply)
	if !ok {
		return nil, kit.ErrReplyTypeInvalid
	}
	return reply, nil
}

func NewTermServiceClient(conn *grpc.ClientConn, tracer kitgrpc.ClientOption) *TermEndpoints {
	endpoints := new(TermEndpoints)
	options := make([]kitgrpc.ClientOption, 0)
	if tracer != nil {
		options = append(options, tracer)
	}

	{
		endpoint := grpctransport.NewClient(
			conn,
			"pb.Term",
			"ListTermInfo",
			encodeRequest,
			decodeResponse,
			pb.ListTermInfoReply{},
			options...,
		).Endpoint()
		endpoints.ListTermInfoEndpoint = endpoint
	}

	{
		endpoint := grpctransport.NewClient(
			conn,
			"pb.Term",
			"SaveTerm",
			encodeRequest,
			decodeResponse,
			pb.SaveTermReply{},
			options...,
		).Endpoint()
		endpoints.SaveTermEndpoint = endpoint
	}

	{
		endpoint := grpctransport.NewClient(
			conn,
			"pb.Term",
			"ListTermRisk",
			encodeRequest,
			decodeResponse,
			pb.ListTermRiskReply{},
			options...,
		).Endpoint()
		endpoints.ListTermRiskEndpoint = endpoint
	}

	{
		endpoint := grpctransport.NewClient(
			conn,
			"pb.Term",
			"ListTermActivationState",
			encodeRequest,
			decodeResponse,
			pb.ListTermActivationStateReply{},
			options...,
		).Endpoint()
		endpoints.ListTermActivationStateEndpoint = endpoint
	}

	{
		endpoint := grpctransport.NewClient(
			conn,
			"pb.Term",
			"UpdateTermInfo",
			encodeRequest,
			decodeResponse,
			pb.UpdateTermInfoReply{},
			options...,
		).Endpoint()
		endpoints.UpdateTermInfoEndpoint = endpoint
	}

	{
		endpoint := grpctransport.NewClient(
			conn,
			"pb.Term",
			"QueryTermInfo",
			encodeRequest,
			decodeResponse,
			pb.QueryTermInfoReply{},
			append(options, grpctransport.ClientBefore(setUserInfoMD))...,
		).Endpoint()
		endpoints.QueryTermInfoEndpoint = endpoint
	}

	return endpoints
}
