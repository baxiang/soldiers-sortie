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

type ExternalTaskEndpoints struct {
	FetchAndLockEndpoint endpoint.Endpoint
	GetEndpoint          endpoint.Endpoint
	GetListEndpoint      endpoint.Endpoint
	CompleteEndpoint     endpoint.Endpoint
}

func (e *ExternalTaskEndpoints) Get(ctx context.Context, in *pb.GetExternalTaskReq) (*pb.GetExternalTaskResp, error) {
	res, err := e.GetEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	reply, ok := res.(*pb.GetExternalTaskResp)
	if !ok {
		return nil, kit.ErrReplyTypeInvalid
	}
	return reply, nil
}

func (e *ExternalTaskEndpoints) GetList(ctx context.Context, in *pb.GetQuery) (*pb.GetListExternalTaskResp, error) {
	res, err := e.GetListEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	reply, ok := res.(*pb.GetListExternalTaskResp)
	if !ok {
		return nil, kit.ErrReplyTypeInvalid
	}
	return reply, nil
}

func (e *ExternalTaskEndpoints) FetchAndLock(ctx context.Context, in *pb.FetchAndLockExternalTaskReq) (*pb.FetchAndLockExternalTaskResp, error) {
	res, err := e.FetchAndLockEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	reply, ok := res.(*pb.FetchAndLockExternalTaskResp)
	if !ok {
		return nil, kit.ErrReplyTypeInvalid
	}
	return reply, nil
}

func (e *ExternalTaskEndpoints) Complete(ctx context.Context, in *pb.CompleteExternalTaskReq) (*pb.CompleteExternalTaskResp, error) {
	res, err := e.CompleteEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	reply, ok := res.(*pb.CompleteExternalTaskResp)
	if !ok {
		return nil, kit.ErrReplyTypeInvalid
	}
	return reply, nil
}

func NewExternalTaskClient(conn *grpc.ClientConn, tracer kitgrpc.ClientOption) *ExternalTaskEndpoints {
	endpoints := new(ExternalTaskEndpoints)

	options := make([]kitgrpc.ClientOption, 0)
	if tracer != nil {
		options = append(options, tracer)
	}

	{
		e := grpctransport.NewClient(
			conn,
			"pb.ExternalTask",
			"Get",
			encodeRequest,
			decodeResponse,
			pb.GetExternalTaskResp{},
			options...,
		).Endpoint()
		endpoints.GetEndpoint = e
	}

	{
		e := grpctransport.NewClient(
			conn,
			"pb.ExternalTask",
			"GetList",
			encodeRequest,
			decodeResponse,
			pb.GetListExternalTaskResp{},
			options...,
		).Endpoint()
		endpoints.GetListEndpoint = e
	}

	{
		e := grpctransport.NewClient(
			conn,
			"pb.ExternalTask",
			"FetchAndLock",
			encodeRequest,
			decodeResponse,
			pb.FetchAndLockExternalTaskResp{},
			options...,
		).Endpoint()
		endpoints.FetchAndLockEndpoint = e
	}

	{
		e := grpctransport.NewClient(
			conn,
			"pb.ExternalTask",
			"Complete",
			encodeRequest,
			decodeResponse,
			pb.CompleteExternalTaskResp{},
			options...,
		).Endpoint()
		endpoints.CompleteEndpoint = e
	}

	return endpoints
}
