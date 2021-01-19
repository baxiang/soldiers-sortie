package gateway

import (
	"context"
	"net/http"
	"userService/pkg/pb"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ClientEndpoints struct {
	UserEndpoints        *UserEndpoints
	StaticEndpoints      *StaticEndpoints
	InstitutionEndpoints *InstitutionEndpoints
	MerchantEndpoints    *MerchantEndpoints
	TermEndpoints        *TermEndpoints
	WorkflowEndpoints    *WorkflowEndpoints
	ScanEndpoints        *ScanEndpoints
}

func NewHttpHandler(c *ClientEndpoints) http.Handler {
	engine := gin.New()
	engine.Use(logRequest)
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))
	RegisterUserHandler(engine, c.UserEndpoints)
	RegisterStaticHandler(engine, c.StaticEndpoints)
	RegisterInstitutionHandler(engine, c.InstitutionEndpoints)
	RegisterMerchantHandler(engine, c.MerchantEndpoints)
	RegisterTermHandler(engine, c.TermEndpoints)
	RegisterWorkflowHandler(engine, c.WorkflowEndpoints)
	RegisterScanHandler(engine, c.ScanEndpoints)
	return engine
}

func encodeRequest(ctx context.Context, request interface{}) (interface{}, error) {
	return request, nil
}

func decodeResponse(ctx context.Context, response interface{}) (interface{}, error) {
	return response, nil
}

type StatusError interface {
	GetErr() *pb.Error
}

func logRequest(c *gin.Context) {
	logrus.Debugln(c.Request.URL.String())
	c.Next()
}
