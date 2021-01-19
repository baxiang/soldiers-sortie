package staticservice

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

//MakeSyncDataEndpoint .
func MakeSyncDataEndpoint(s pb.StaticServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.StaticSyncDataReq)
		return s.SyncData(ctx, req)
	}
}

//MakeGetDictionaryItemEndpoint .
func MakeGetDictionaryItemEndpoint(s pb.StaticServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.StaticGetDictionaryItemReq)
		return s.GetDictionaryItem(ctx, req)
	}
}

//MakeGetDicByProdAndBizEndpoint .
func MakeGetDicByProdAndBizEndpoint(s pb.StaticServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.StaticGetDicByProdAndBizReq)
		return s.GetDicByProdAndBiz(ctx, req)
	}
}

func MakeCheckValuesEndpoint(s pb.StaticServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.StaticCheckValuesReq)
		return s.CheckValues(ctx, req)
	}
}

func MakeGetDictionaryLayerItemEndpoint(s pb.StaticServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.GetDictionaryLayerItemReq)
		return s.GetDictionaryLayerItem(ctx, req)
	}
}

func MakeGetDictionaryItemByPkEndpoint(s pb.StaticServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.GetDictionaryItemByPkReq)
		return s.GetDictionaryItemByPk(ctx, req)
	}
}

func MakeGetUnionPayBankListEndpoint(s pb.StaticServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.GetUnionPayBankList(ctx, request.(*pb.GetUnionPayBankListRequest))
	}
}

func MakeFindUnionPayMccListEndpoint(s pb.StaticServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.FindUnionPayMccList(ctx, request.(*pb.FindUnionPayMccListRequest))
	}
}

func MakeGetInsProdBizFeeMapInfoEndpoint(s pb.StaticServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return s.GetInsProdBizFeeMapInfo(ctx, request.(*pb.GetInsProdBizFeeMapInfoRequest))
	}
}

func MakeListTransMapEndpoint(service pb.StaticServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return service.ListTransMap(ctx, request.(*pb.ListTransMapRequest))
	}
}

func MakeListFeeMapEndpoint(service pb.StaticServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return service.ListFeeMap(ctx, request.(*pb.ListFeeMapRequest))
	}
}

func MakeFindAreaEndpoint(service pb.StaticServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return service.FindArea(ctx, request.(*pb.FindAreaRequest))
	}
}

func MakeFindMerchantFirstThreeCodeEndpoint(service pb.StaticServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return service.FindMerchantFirstThreeCode(ctx, request.(*pb.FindMerchantFirstThreeCodeRequest))
	}
}

func MakeSaveOrgDictionaryItemEndpoint(service pb.StaticServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return service.SaveOrgDictionaryItem(ctx, request.(*pb.SaveOrgDictionaryItemRequest))
	}
}

func MakeListOrgDictionaryItemEndpoint(service pb.StaticServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return service.ListOrgDictionaryItem(ctx, request.(*pb.ListOrgDictionaryItemRequest))
	}
}

func MakeSaveInsProdBizFeeMapInfoEndpoint(service pb.StaticServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return service.SaveInsProdBizFeeMapInfo(ctx, request.(*pb.SaveInsProdBizFeeMapInfoRequest))
	}
}
