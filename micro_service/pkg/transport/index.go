package transport

import (
	"context"
	"userService/pkg/camunda/pb"

	"github.com/go-kit/kit/endpoint"
)

func MakeGetProcessDefinitionEndpoint(service pb.ProcessDefinitionServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.GetProcessDefinitionReq)
		if !ok {
			return nil, ErrRequestTypeInvalid
		}
		return service.Get(ctx, req)
	}
}

func MakeStartProcessDefinitionEndpoint(service pb.ProcessDefinitionServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.StartProcessDefinitionReq)
		if !ok {
			return nil, ErrRequestTypeInvalid
		}
		return service.Start(ctx, req)
	}
}

func MakeListProcessInstanceEndpoint(service pb.ProcessInstanceServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.ProcessInstanceListReq)
		if !ok {
			return nil, ErrRequestTypeInvalid
		}
		return service.List(ctx, req)
	}
}
func MakeTaskGetListEndpoint(service pb.TaskServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.GetListTaskReq)
		if !ok {
			return nil, ErrRequestTypeInvalid
		}
		return service.GetList(ctx, req)
	}
}
func MakeTaskGetEndpoint(service pb.TaskServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.GetTaskReq)
		if !ok {
			return nil, ErrRequestTypeInvalid
		}
		return service.Get(ctx, req)
	}
}
func MakeTaskCompleteEndpoint(service pb.TaskServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.CompleteTaskReq)
		if !ok {
			return nil, ErrRequestTypeInvalid
		}
		return service.Complete(ctx, req)
	}
}

func MakeGetTaskFormValueEndpoint(s pb.TaskServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.GetFormValueRequest)
		return s.GetFormValue(ctx, req)
	}
}

func MakeExternalTaskGetEndpoint(service pb.ExternalTaskServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return service.Get(ctx, request.(*pb.GetExternalTaskReq))
	}
}

func MakeExternalTaskGetListEndpoint(service pb.ExternalTaskServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return service.GetList(ctx, request.(*pb.GetQuery))
	}
}

func MakeExternalTaskFetchAndLockEndpoint(service pb.ExternalTaskServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return service.FetchAndLock(ctx, request.(*pb.FetchAndLockExternalTaskReq))
	}
}

func MakeExternalTaskCompleteEndpoint(service pb.ExternalTaskServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return service.Complete(ctx, request.(*pb.CompleteExternalTaskReq))
	}
}
