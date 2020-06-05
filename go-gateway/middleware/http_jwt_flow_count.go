package middleware

import (
	"github.com/baxiang/go-gateway/dao"
	"github.com/baxiang/go-gateway/pkg"
	"github.com/gin-gonic/gin"
	"errors"
	"fmt"
)

func HTTPJwtFlowCountMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		appInterface, ok := c.Get("app")
		if !ok {
			c.Next()
			return
		}
		appInfo := appInterface.(*dao.App)
		appCounter, err := pkg.FlowCounterHandler.GetCounter(pkg.FlowAppPrefix + appInfo.AppID)
		if err != nil {
			ResponseError(c, 2002, err)
			c.Abort()
			return
		}
		appCounter.Increase()
		if appInfo.Qpd>0 && appCounter.TotalCount>appInfo.Qpd{
			ResponseError(c, 2003, errors.New(fmt.Sprintf("租户日请求量限流 limit:%v current:%v",appInfo.Qpd,appCounter.TotalCount)))
			c.Abort()
			return
		}
		c.Next()
	}
}