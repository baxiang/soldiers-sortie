package gateway

import (
	"userService/pkg/pb"
	"userService/pkg/userservice"

	stdjwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	httptransport "github.com/go-kit/kit/transport/http"
)

func RegisterWorkflowHandler(engine *gin.Engine, endpoints *WorkflowEndpoints) {
	group := engine.Group("/workflow")

	group.POST("/start",
		userservice.JwtMiddleware(keyFunc, stdjwt.SigningMethodHS256, userservice.UserClaimFactory),
		convertHttpHandlerToGinHandler(httptransport.NewServer(
			endpoints.StartEndpoint,
			decodeHttpRequest(&pb.StartWorkflowRequest{}),
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(errorEncoder),
			httptransport.ServerBefore(setUserInfoContext),
		)))

	group.POST("/listTask",
		userservice.JwtMiddleware(keyFunc, stdjwt.SigningMethodHS256, userservice.UserClaimFactory),
		convertHttpHandlerToGinHandler(httptransport.NewServer(
			endpoints.ListTaskEndpoint,
			decodeHttpRequest(&pb.ListTaskRequest{}),
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(errorEncoder),
			httptransport.ServerBefore(setUserInfoContext),
		)))

	group.POST("/handleTask",
		userservice.JwtMiddleware(keyFunc, stdjwt.SigningMethodHS256, userservice.UserClaimFactory),
		convertHttpHandlerToGinHandler(httptransport.NewServer(
			endpoints.HandleTaskEndpoint,
			decodeHttpRequest(&pb.HandleTaskRequest{}),
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(errorEncoder),
			httptransport.ServerBefore(setUserInfoContext),
		)))

	group.POST("/listRemark",
		convertHttpHandlerToGinHandler(httptransport.NewServer(
			endpoints.ListRemarkEndpoint,
			decodeHttpRequest(&pb.ListRemarkRequest{}),
			encodeHttpResponse,
			httptransport.ServerErrorEncoder(errorEncoder),
		)))
}
