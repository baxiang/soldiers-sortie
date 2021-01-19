package gateway

import (
	"userService/pkg/pb"

	"github.com/gin-gonic/gin"
	httptransport "github.com/go-kit/kit/transport/http"
)

func RegisterInstitutionHandler(engine *gin.Engine, endpoints *InstitutionEndpoints) {
	engine.POST("/institution/tnxHisDownload", convertHttpHandlerToGinHandler(httptransport.NewServer(
		endpoints.TnxHisDownloadEndpoint,
		decodeHttpRequest(&pb.InstitutionTnxHisDownloadReq{}),
		encodeHttpResponse,
		httptransport.ServerErrorEncoder(errorEncoder),
	)))

	engine.POST("/institution/getTfrTrnLogs", convertHttpHandlerToGinHandler(httptransport.NewServer(
		endpoints.GetTfrTrnLogsEndpoint,
		decodeHttpRequest(&pb.GetTfrTrnLogsReq{}),
		encodeHttpResponse,
		httptransport.ServerErrorEncoder(errorEncoder),
	)))

	engine.POST("/institution/getTfrTrnLog", convertHttpHandlerToGinHandler(httptransport.NewServer(
		endpoints.GetTfrTrnLogEndpoint,
		decodeHttpRequest(&pb.GetTfrTrnLogReq{}),
		encodeHttpResponse,
		httptransport.ServerErrorEncoder(errorEncoder),
	)))

	engine.POST("/institution/downloadTfrTrnLogs", convertHttpHandlerToGinHandler(httptransport.NewServer(
		endpoints.DownloadTfrTrnLogsEndpoint,
		decodeHttpRequest(&pb.DownloadTfrTrnLogsReq{}),
		encodeHttpResponse,
		httptransport.ServerErrorEncoder(errorEncoder),
	)))

	engine.POST("/institution/listGroups", convertHttpHandlerToGinHandler(httptransport.NewServer(
		endpoints.ListGroupsEndpoint,
		decodeHttpRequest(&pb.ListGroupsRequest{}),
		encodeHttpResponse,
		httptransport.ServerErrorEncoder(errorEncoder),
	)))

	engine.POST("/institution/listInstitutions", convertHttpHandlerToGinHandler(httptransport.NewServer(
		endpoints.ListInstitutionsEndpoint,
		decodeHttpRequest(&pb.ListInstitutionsRequest{}),
		encodeHttpResponse,
		httptransport.ServerErrorEncoder(errorEncoder),
	)))

	engine.POST("/institution/saveInstitution", convertHttpHandlerToGinHandler(httptransport.NewServer(
		endpoints.SaveInstitutionEndpoint,
		decodeHttpRequest(&pb.SaveInstitutionRequest{}),
		encodeHttpResponse,
		httptransport.ServerErrorEncoder(errorEncoder),
	)))

	engine.POST("/institution/getInstitutionById", convertHttpHandlerToGinHandler(httptransport.NewServer(
		endpoints.GetInstitutionByIdEndpoint,
		decodeHttpRequest(&pb.GetInstitutionByIdRequest{}),
		encodeHttpResponse,
		httptransport.ServerErrorEncoder(errorEncoder),
	)))

	engine.POST("/institution/saveInstitutionFeeControlCash", convertHttpHandlerToGinHandler(httptransport.NewServer(
		endpoints.SaveInstitutionFeeControlCashEndpoint,
		decodeHttpRequest(&pb.SaveInstitutionFeeControlCashRequest{}),
		encodeHttpResponse,
		httptransport.ServerErrorEncoder(errorEncoder),
	)))

	engine.POST("/institution/getInstitutionControl", convertHttpHandlerToGinHandler(httptransport.NewServer(
		endpoints.GetInstitutionControlEndpoint,
		decodeHttpRequest(&pb.GetInstitutionControlRequest{}),
		encodeHttpResponse,
		httptransport.ServerErrorEncoder(errorEncoder),
	)))

	engine.POST("/institution/getInstitutionCash", convertHttpHandlerToGinHandler(httptransport.NewServer(
		endpoints.GetInstitutionCashEndpoint,
		decodeHttpRequest(&pb.GetInstitutionCashRequest{}),
		encodeHttpResponse,
		httptransport.ServerErrorEncoder(errorEncoder),
	)))

	engine.POST("/institution/getInstitutionFee", convertHttpHandlerToGinHandler(httptransport.NewServer(
		endpoints.GetInstitutionFeeEndpoint,
		decodeHttpRequest(&pb.GetInstitutionFeeRequest{}),
		encodeHttpResponse,
		httptransport.ServerErrorEncoder(errorEncoder),
	)))

	engine.POST("/institution/saveGroup", convertHttpHandlerToGinHandler(httptransport.NewServer(
		endpoints.SaveGroupEndpoint,
		decodeHttpRequest(&pb.SaveGroupRequest{}),
		encodeHttpResponse,
		httptransport.ServerErrorEncoder(errorEncoder),
	)))

	engine.POST("/institution/bindGroup", convertHttpHandlerToGinHandler(httptransport.NewServer(
		endpoints.BindGroupEndpoint,
		decodeHttpRequest(&pb.BindGroupRequest{}),
		encodeHttpResponse,
		httptransport.ServerErrorEncoder(errorEncoder),
	)))

	engine.POST("/institution/listBindGroup", convertHttpHandlerToGinHandler(httptransport.NewServer(
		endpoints.ListBindGroupEndpoint,
		decodeHttpRequest(&pb.ListBindGroupRequest{}),
		encodeHttpResponse,
		httptransport.ServerErrorEncoder(errorEncoder),
	)))

	engine.POST("/institution/removeBindGroup", convertHttpHandlerToGinHandler(httptransport.NewServer(
		endpoints.RemoveBindGroupEndpoint,
		decodeHttpRequest(&pb.RemoveBindGroupRequest{}),
		encodeHttpResponse,
		httptransport.ServerErrorEncoder(errorEncoder),
	)))

}
