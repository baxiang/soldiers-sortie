package gateway

import (
	"userService/pkg/pb"
	"userService/pkg/userservice"

	stdjwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	httptransport "github.com/go-kit/kit/transport/http"
)

func RegisterMerchantHandler(engine *gin.Engine, endpoints *MerchantEndpoints) {
	group := engine.Group("/merchant")

	group.POST("/listMerchant",
		userservice.JwtMiddleware(keyFunc, stdjwt.SigningMethodHS256, userservice.UserClaimFactory),
		convertHttpHandlerToGinHandler(httptransport.NewServer(
			endpoints.ListMerchantEndpoint,
			decodeHttpRequest(&pb.ListMerchantRequest{}),
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(errorEncoder),
			httptransport.ServerBefore(setUserInfoContext),
		)))

	group.POST("/listGroup", convertHttpHandlerToGinHandler(httptransport.NewServer(
		endpoints.ListGroupMerchantEndpoint,
		decodeHttpRequest(&pb.ListGroupMerchantRequest{}),
		encodeHttpResponse,
		httptransport.ServerErrorEncoder(errorEncoder),
	)))

	group.POST("/saveMerchant", convertHttpHandlerToGinHandler(httptransport.NewServer(
		endpoints.SaveMerchantEndpoint,
		decodeHttpRequest(&pb.SaveMerchantRequest{}),
		encodeHttpResponse,
		httptransport.ServerErrorEncoder(errorEncoder),
	)))

	group.POST("/saveBankAccount", convertHttpHandlerToGinHandler(httptransport.NewServer(
		endpoints.SaveMerchantBankAccountEndpoint,
		decodeHttpRequest(&pb.SaveMerchantBankAccountRequest{}),
		encodeHttpResponse,
		httptransport.ServerErrorEncoder(errorEncoder),
	)))

	group.POST("/saveMerchantBusiness", convertHttpHandlerToGinHandler(httptransport.NewServer(
		endpoints.SaveMerchantBusinessEndpoint,
		decodeHttpRequest(&pb.SaveMerchantBusinessRequest{}),
		encodeHttpResponse,
		httptransport.ServerErrorEncoder(errorEncoder),
	)))

	group.POST("/saveMerchantPicture", convertHttpHandlerToGinHandler(httptransport.NewServer(
		endpoints.SaveMerchantPictureEndpoint,
		decodeHttpRequest(&pb.SaveMerchantPictureRequest{}),
		encodeHttpResponse,
		httptransport.ServerErrorEncoder(errorEncoder),
	)))

	group.POST("/getMerchantBankAccount", convertHttpHandlerToGinHandler(httptransport.NewServer(
		endpoints.GetMerchantBankAccountEndpoint,
		decodeHttpRequest(&pb.GetMerchantBankAccountRequest{}),
		encodeHttpResponse,
		httptransport.ServerErrorEncoder(errorEncoder),
	)))

	group.POST("/getMerchantBizDeal", convertHttpHandlerToGinHandler(httptransport.NewServer(
		endpoints.GetMerchantBizDealEndpoint,
		decodeHttpRequest(&pb.GetMerchantBizDealRequest{}),
		encodeHttpResponse,
		httptransport.ServerErrorEncoder(errorEncoder),
	)))

	group.POST("/getMerchantBizFee", convertHttpHandlerToGinHandler(httptransport.NewServer(
		endpoints.GetMerchantBizFeeEndpoint,
		decodeHttpRequest(&pb.GetMerchantBizFeeRequest{}),
		encodeHttpResponse,
		httptransport.ServerErrorEncoder(errorEncoder),
	)))

	group.POST("/getMerchantBusiness", convertHttpHandlerToGinHandler(httptransport.NewServer(
		endpoints.GetMerchantBusinessEndpoint,
		decodeHttpRequest(&pb.GetMerchantBusinessRequest{}),
		encodeHttpResponse,
		httptransport.ServerErrorEncoder(errorEncoder),
	)))

	group.POST("/getMerchantPicture", convertHttpHandlerToGinHandler(httptransport.NewServer(
		endpoints.GetMerchantPictureEndpoint,
		decodeHttpRequest(&pb.GetMerchantPictureRequest{}),
		encodeHttpResponse,
		httptransport.ServerErrorEncoder(errorEncoder),
	)))

	group.POST("/getMerchantById", convertHttpHandlerToGinHandler(httptransport.NewServer(
		endpoints.GetMerchantByIdEndpoint,
		decodeHttpRequest(&pb.GetMerchantByIdRequest{}),
		encodeHttpResponse,
		httptransport.ServerErrorEncoder(errorEncoder),
	)))

	group.POST("/saveMerchantBizDealAndFee", convertHttpHandlerToGinHandler(httptransport.NewServer(
		endpoints.SaveMerchantBizDealAndFeeEndpoint,
		decodeHttpRequest(&pb.SaveMerchantBizDealAndFeeRequest{}),
		encodeHttpResponse,
		httptransport.ServerErrorEncoder(errorEncoder),
	)))

	group.POST("/generateMchtCd", convertHttpHandlerToGinHandler(httptransport.NewServer(
		endpoints.GenerateMchtCdEndpoint,
		decodeHttpRequest(&pb.GenerateMchtCdRequest{}),
		encodeHttpResponse,
		httptransport.ServerErrorEncoder(errorEncoder),
	)))

	group.POST("/merchantInfoQuery",
		userservice.JwtMiddleware(keyFunc, stdjwt.SigningMethodHS256, userservice.UserClaimFactory),
		convertHttpHandlerToGinHandler(httptransport.NewServer(
			endpoints.MerchantInfoQueryEndpoint,
			decodeHttpRequest(&pb.MerchantInfoQueryRequest{}),
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(errorEncoder),
			httptransport.ServerBefore(setUserInfoContext),
		)))

	group.POST("/merchantForceChangeStatus",
		convertHttpHandlerToGinHandler(httptransport.NewServer(
			endpoints.MerchantForceChangeStatusEndpoint,
			decodeHttpRequest(&pb.MerchantForceChangeStatusRequest{}),
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(errorEncoder),
		)))
}
