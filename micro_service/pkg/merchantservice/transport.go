package merchantservice

import (
	"context"

	"userService/pkg/kit"
	"userService/pkg/pb"

	"github.com/go-kit/kit/endpoint"
)

func MakeListMerchantEndpoint(service pb.MerchantServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.ListMerchantRequest)
		if !ok {
			return nil, kit.ErrRequestTypeInvalid
		}
		return service.ListMerchant(ctx, req)
	}
}

func MakeListGroupMerchantEndpoint(service pb.MerchantServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*pb.ListGroupMerchantRequest)
		if !ok {
			return nil, kit.ErrRequestTypeInvalid
		}
		return service.ListGroupMerchant(ctx, req)
	}
}

func MakeSaveMerchantEndpoint(service pb.MerchantServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return service.SaveMerchant(ctx, request.(*pb.SaveMerchantRequest))
	}
}

func MakeSaveMerchantBankAccountEndpoint(service pb.MerchantServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return service.SaveMerchantBankAccount(ctx, request.(*pb.SaveMerchantBankAccountRequest))
	}
}

func MakeSaveMerchantBusinessEndpoint(service pb.MerchantServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return service.SaveMerchantBusiness(ctx, request.(*pb.SaveMerchantBusinessRequest))
	}
}

func MakeSaveMerchantPictureEndpoint(service pb.MerchantServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return service.SaveMerchantPicture(ctx, request.(*pb.SaveMerchantPictureRequest))
	}
}

func MakeGetMerchantBankAccountEndpoint(service pb.MerchantServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return service.GetMerchantBankAccount(ctx, request.(*pb.GetMerchantBankAccountRequest))
	}
}

func MakeGetMerchantBizDealEndpoint(service pb.MerchantServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return service.GetMerchantBizDeal(ctx, request.(*pb.GetMerchantBizDealRequest))
	}
}

func MakeGetMerchantBizFeeEndpoint(service pb.MerchantServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return service.GetMerchantBizFee(ctx, request.(*pb.GetMerchantBizFeeRequest))
	}
}

func MakeGetMerchantBusinessEndpoint(service pb.MerchantServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return service.GetMerchantBusiness(ctx, request.(*pb.GetMerchantBusinessRequest))
	}
}

func MakeGetMerchantPictureEndpoint(service pb.MerchantServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return service.GetMerchantPicture(ctx, request.(*pb.GetMerchantPictureRequest))
	}
}

func MakeGetMerchantByIdEndpoint(service pb.MerchantServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return service.GetMerchantById(ctx, request.(*pb.GetMerchantByIdRequest))
	}
}

func MakeSaveMerchantBizDealAndFeeEndpoint(service pb.MerchantServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return service.SaveMerchantBizDealAndFee(ctx, request.(*pb.SaveMerchantBizDealAndFeeRequest))
	}
}

func MakeGenerateMchtCdEndpoint(service pb.MerchantServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return service.GenerateMchtCd(ctx, request.(*pb.GenerateMchtCdRequest))
	}
}

func MakeMerchantInfoQueryEndpoint(service pb.MerchantServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return service.MerchantInfoQuery(ctx, request.(*pb.MerchantInfoQueryRequest))
	}
}

func MakeMerchantForceChangeStatusEndpoint(service pb.MerchantServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return service.MerchantForceChangeStatus(ctx, request.(*pb.MerchantForceChangeStatusRequest))
	}
}
