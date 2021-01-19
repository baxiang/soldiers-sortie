package institutionservice

import (
	"userService/pkg/pb"

	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
)

func grpcDecode(_ context.Context, req interface{}) (interface{}, error) {
	return req, nil
}

func grpcEncode(_ context.Context, res interface{}) (interface{}, error) {
	return res, nil
}

func MakeTnxHisDownloadEndpoint(s pb.InstitutionServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.InstitutionTnxHisDownloadReq)
		return s.TnxHisDownload(ctx, req)
	}
}

func MakeGetTfrTrnLogsEndpoint(s pb.InstitutionServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.GetTfrTrnLogsReq)
		return s.GetTfrTrnLogs(ctx, req)
	}
}

func MakeGetTfrTrnLogEndpoint(s pb.InstitutionServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.GetTfrTrnLogReq)
		return s.GetTfrTrnLog(ctx, req)
	}
}

func MakeDownloadTfrTrnLogsEndpoint(s pb.InstitutionServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.DownloadTfrTrnLogsReq)
		return s.DownloadTfrTrnLogs(ctx, req)
	}
}

func MakeListGroupsEndpoint(service pb.InstitutionServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.ListGroupsRequest)
		if !ok {
			return nil, ErrRequestTypeInvalid
		}
		return service.ListGroups(ctx, req)
	}
}

func MakeListInstitutionsEndpoint(service pb.InstitutionServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.ListInstitutionsRequest)
		if !ok {
			return nil, ErrRequestTypeInvalid
		}
		return service.ListInstitutions(ctx, req)
	}
}

func MakeSaveInstitutionEndpoint(service pb.InstitutionServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.SaveInstitutionRequest)
		if !ok {
			return nil, ErrRequestTypeInvalid
		}
		return service.SaveInstitution(ctx, req)
	}
}

func MakeGetInstitutionByIdEndpoint(service pb.InstitutionServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return service.GetInstitutionById(ctx, request.(*pb.GetInstitutionByIdRequest))
	}
}

func MakeSaveInstitutionFeeControlCashEndpoint(service pb.InstitutionServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return service.SaveInstitutionFeeControlCash(ctx, request.(*pb.SaveInstitutionFeeControlCashRequest))
	}
}

func MakeGetInstitutionControlEndpoint(service pb.InstitutionServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return service.GetInstitutionControl(ctx, request.(*pb.GetInstitutionControlRequest))
	}
}

func MakeGetInstitutionCashEndpoint(service pb.InstitutionServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return service.GetInstitutionCash(ctx, request.(*pb.GetInstitutionCashRequest))
	}
}

func MakeGetInstitutionFeeEndpoint(service pb.InstitutionServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return service.GetInstitutionFee(ctx, request.(*pb.GetInstitutionFeeRequest))
	}
}

func MakeSaveGroupEndpoint(service pb.InstitutionServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return service.SaveGroup(ctx, request.(*pb.SaveGroupRequest))
	}
}

func MakeBindGroupEndpoint(service pb.InstitutionServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return service.BindGroup(ctx, request.(*pb.BindGroupRequest))
	}
}

func MakeListBindGroupEndpoint(service pb.InstitutionServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return service.ListBindGroup(ctx, request.(*pb.ListBindGroupRequest))
	}
}

func MakeRemoveBindGroupEndpoint(service pb.InstitutionServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return service.RemoveBindGroup(ctx, request.(*pb.RemoveBindGroupRequest))
	}
}
