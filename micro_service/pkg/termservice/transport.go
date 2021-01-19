package termservice

import (
	"context"
	"userService/pkg/pb"

	"github.com/go-kit/kit/endpoint"
)

func MakeListTermInfoEndpoint(service pb.TermServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return service.ListTermInfo(ctx, request.(*pb.ListTermInfoRequest))
	}
}

func MakeSaveTermEndpoint(service pb.TermServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return service.SaveTerm(ctx, request.(*pb.SaveTermRequest))
	}
}

func MakeListTermRiskEndpoint(service pb.TermServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return service.ListTermRisk(ctx, request.(*pb.ListTermRiskRequest))
	}
}

func MakeListTermActivationStateEndpoint(service pb.TermServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return service.ListTermActivationState(ctx, request.(*pb.ListTermActivationStateRequest))
	}
}

func MakeUpdateTermInfoEndpoint(service pb.TermServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return service.UpdateTermInfo(ctx, request.(*pb.UpdateTermInfoRequest))
	}
}

func MakeQueryTermInfoEndpoint(service pb.TermServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return service.QueryTermInfo(ctx, request.(*pb.QueryTermInfoRequest))
	}
}
