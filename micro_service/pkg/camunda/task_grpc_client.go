package camunda

import (
	"context"
	"userService/pkg/camunda/pb"

	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

type TaskEndpoints struct {
	GetListEndpoint      endpoint.Endpoint
	GetEndpoint          endpoint.Endpoint
	CompleteEndpoint     endpoint.Endpoint
	GetFormValueEndpoint endpoint.Endpoint
}

func (t *TaskEndpoints) GetFormValue(ctx context.Context, in *pb.GetFormValueRequest) (*pb.GetFormValueReply, error) {
	res, err := t.GetFormValueEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.GetFormValueReply), nil
}

func (t *TaskEndpoints) GetList(ctx context.Context, in *pb.GetListTaskReq) (*pb.GetListTaskResp, error) {
	res, err := t.GetListEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.GetListTaskResp), nil
}
func (t *TaskEndpoints) Get(ctx context.Context, in *pb.GetTaskReq) (*pb.GetTaskResp, error) {
	res, err := t.GetEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.GetTaskResp), nil
}
func (t *TaskEndpoints) Complete(ctx context.Context, in *pb.CompleteTaskReq) (*pb.CompleteTaskResp, error) {
	res, err := t.CompleteEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return res.(*pb.CompleteTaskResp), nil
}

func NewTaskClient(conn *grpc.ClientConn, tracer kitgrpc.ClientOption) *TaskEndpoints {
	endpoints := new(TaskEndpoints)

	options := make([]kitgrpc.ClientOption, 0)
	if tracer != nil {
		options = append(options, tracer)
	}

	{
		e := grpctransport.NewClient(
			conn,
			"pb.Task",
			"GetList",
			encodeRequest,
			decodeResponse,
			pb.GetListTaskResp{},
			options...,
		).Endpoint()
		endpoints.GetListEndpoint = e
	}

	{
		e := grpctransport.NewClient(
			conn,
			"pb.Task",
			"Get",
			encodeRequest,
			decodeResponse,
			pb.GetTaskResp{},
			options...,
		).Endpoint()
		endpoints.GetEndpoint = e
	}

	{
		e := grpctransport.NewClient(
			conn,
			"pb.Task",
			"Complete",
			encodeRequest,
			decodeResponse,
			pb.CompleteTaskResp{},
			options...,
		).Endpoint()
		endpoints.CompleteEndpoint = e
	}

	{
		e := grpctransport.NewClient(
			conn,
			"pb.Task",
			"GetFormValue",
			encodeRequest,
			decodeResponse,
			pb.GetFormValueReply{},
			options...,
		).Endpoint()
		endpoints.GetFormValueEndpoint = e
	}

	return endpoints
}
