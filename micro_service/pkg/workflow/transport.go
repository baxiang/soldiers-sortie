package workflow

import (
	"context"
	"userService/pkg/pb"

	"github.com/go-kit/kit/endpoint"
)

func decodeRequest(ctx context.Context, request interface{}) (interface{}, error) {
	return request, nil
}

func encodeResponse(ctx context.Context, response interface{}) (interface{}, error) {
	return response, nil
}

func MakeListTaskEndpoint(service pb.WorkflowServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return service.ListTask(ctx, request.(*pb.ListTaskRequest))
	}
}

func MakeHandleTaskEndpoint(service pb.WorkflowServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return service.HandleTask(ctx, request.(*pb.HandleTaskRequest))
	}
}

func MakeStartWorkflowEndpoint(service pb.WorkflowServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return service.Start(ctx, request.(*pb.StartWorkflowRequest))
	}
}

func MakeListRemarkEndpoint(service pb.WorkflowServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return service.ListRemark(ctx, request.(*pb.ListRemarkRequest))
	}
}
