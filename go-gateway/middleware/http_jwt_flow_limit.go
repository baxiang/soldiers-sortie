package middleware

import (
	"github.com/baxiang/go-gateway/dao"
	"github.com/baxiang/go-gateway/pkg"
	"github.com/gin-gonic/gin"
	"errors"
	"fmt"
)

func HTTPJwtFlowLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		appInterface, ok := c.Get("app")
		if !ok {
			c.Next()
			return
		}
		appInfo := appInterface.(*dao.App)
		if appInfo.Qps > 0 {
			clientLimiter, err := pkg.FlowLimiterHandler.GetLimiter(
				pkg.FlowAppPrefix+appInfo.AppID+"_"+c.ClientIP(),
				float64(appInfo.Qps))
			if err != nil {
				ResponseError(c, 5001, err)
				c.Abort()
				return
			}
			if !clientLimiter.Allow() {
				ResponseError(c, 5002, errors.New(fmt.Sprintf("%v flow limit %v", c.ClientIP(), appInfo.Qps), ))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}