package gateway

import (
	"userService/pkg/pb"

	"github.com/gin-gonic/gin"
	httptransport "github.com/go-kit/kit/transport/http"
)

func RegisterStaticHandler(engine *gin.Engine, endpoints *StaticEndpoints) {
	group := engine.Group("/static")

	group.POST("/syncData", convertHttpHandlerToGinHandler(httptransport.NewServer(
		endpoints.SyncDataEndpoint,
		decodeHttpRequest(&pb.StaticSyncDataReq{}),
		encodeHttpResponse,
		httptransport.ServerErrorEncoder(errorEncoder),
	)))

	group.POST("/getDictionaryItem", convertHttpHandlerToGinHandler(httptransport.NewServer(
		endpoints.GetDictionaryItemEndpoint,
		decodeHttpRequest(&pb.StaticGetDictionaryItemReq{}),
		encodeHttpResponse,
		httptransport.ServerErrorEncoder(errorEncoder),
	)))

	group.POST("/getDicByProdAndBiz", convertHttpHandlerToGinHandler(httptransport.NewServer(
		endpoints.GetDicByProdAndBizEndpoint,
		decodeHttpRequest(&pb.StaticGetDicByProdAndBizReq{}),
		encodeHttpResponse,
		httptransport.ServerErrorEncoder(errorEncoder),
	)))

	group.POST("/checkValues", convertHttpHandlerToGinHandler(httptransport.NewServer(
		endpoints.CheckValuesEndpoint,
		decodeHttpRequest(&pb.StaticCheckValuesReq{}),
		encodeHttpResponse,
		httptransport.ServerErrorEncoder(errorEncoder),
	)))

	group.POST("/getDictionaryLayerItem", convertHttpHandlerToGinHandler(httptransport.NewServer(
		endpoints.GetDictionaryLayerItemEndpoint,
		decodeHttpRequest(&pb.GetDictionaryLayerItemReq{}),
		encodeHttpResponse,
		httptransport.ServerErrorEncoder(errorEncoder),
	)))

	group.POST("/getDictionaryItemByPk", convertHttpHandlerToGinHandler(httptransport.NewServer(
		endpoints.GetDictionaryItemByPkEndpoint,
		decodeHttpRequest(&pb.GetDictionaryItemByPkReq{}),
		encodeHttpResponse,
		httptransport.ServerErrorEncoder(errorEncoder),
	)))

	group.POST("/getUnionPayBankList", convertHttpHandlerToGinHandler(httptransport.NewServer(
		endpoints.GetUnionPayBankListEndpoint,
		decodeHttpRequest(&pb.GetUnionPayBankListRequest{}),
		encodeHttpResponse,
		httptransport.ServerErrorEncoder(errorEncoder),
	)))

	group.POST("/findUnionPayMccList", convertHttpHandlerToGinHandler(httptransport.NewServer(
		endpoints.FindUnionPayMccListEndpoint,
		decodeHttpRequest(&pb.FindUnionPayMccListRequest{}),
		encodeHttpResponse,
		httptransport.ServerErrorEncoder(errorEncoder),
	)))

	group.POST("/getInsProdBizFeeMapInfo", convertHttpHandlerToGinHandler(httptransport.NewServer(
		endpoints.GetInsProdBizFeeMapInfoEndpoint,
		decodeHttpRequest(&pb.GetInsProdBizFeeMapInfoRequest{}),
		encodeHttpResponse,
		httptransport.ServerErrorEncoder(errorEncoder),
	)))

	group.POST("/listTransMap", convertHttpHandlerToGinHandler(httptransport.NewServer(
		endpoints.ListTransMapEndpoint,
		decodeHttpRequest(&pb.ListTransMapRequest{}),
		encodeHttpResponse,
		httptransport.ServerErrorEncoder(errorEncoder),
	)))

	group.POST("/listFeeMap", convertHttpHandlerToGinHandler(httptransport.NewServer(
		endpoints.ListFeeMapEndpoint,
		decodeHttpRequest(&pb.ListFeeMapRequest{}),
		encodeHttpResponse,
		httptransport.ServerErrorEncoder(errorEncoder),
	)))

	group.POST("/findArea", convertHttpHandlerToGinHandler(httptransport.NewServer(
		endpoints.FindAreaEndpoint,
		decodeHttpRequest(&pb.FindAreaRequest{}),
		encodeHttpResponse,
		httptransport.ServerErrorEncoder(errorEncoder),
	)))

	group.POST("/findMerchantFirstThreeCode", convertHttpHandlerToGinHandler(httptransport.NewServer(
		endpoints.FindMerchantFirstThreeCodeEndpoint,
		decodeHttpRequest(&pb.FindMerchantFirstThreeCodeRequest{}),
		encodeHttpResponse,
		httptransport.ServerErrorEncoder(errorEncoder),
	)))

	group.POST("/saveOrgDictionaryItem", convertHttpHandlerToGinHandler(httptransport.NewServer(
		endpoints.SaveOrgDictionaryItemEndpoint,
		decodeHttpRequest(&pb.SaveOrgDictionaryItemRequest{}),
		encodeHttpResponse,
		httptransport.ServerErrorEncoder(errorEncoder),
	)))

	group.POST("/listOrgDictionaryItem", convertHttpHandlerToGinHandler(httptransport.NewServer(
		endpoints.ListOrgDictionaryItemEndpoint,
		decodeHttpRequest(&pb.ListOrgDictionaryItemRequest{}),
		encodeHttpResponse,
		httptransport.ServerErrorEncoder(errorEncoder),
	)))

	group.POST("/saveInsProdBizFeeMapInfo", convertHttpHandlerToGinHandler(httptransport.NewServer(
		endpoints.SaveInsProdBizFeeMapInfoEndpoint,
		decodeHttpRequest(&pb.SaveInsProdBizFeeMapInfoRequest{}),
		encodeHttpResponse,
		httptransport.ServerErrorEncoder(errorEncoder),
	)))
}
