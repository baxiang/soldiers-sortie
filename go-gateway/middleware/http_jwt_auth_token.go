package middleware

import (
	"github.com/baxiang/go-gateway/dao"
	"github.com/baxiang/go-gateway/pkg"
	"github.com/gin-gonic/gin"
	"strings"
	"errors"
)

//jwt auth token
func HTTPJwtAuthTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			ResponseError(c, 2001, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)

		//fmt.Println("serviceDetail",serviceDetail)
		// decode jwt token
		// app_id 与  app_list 取得 appInfo
		// appInfo 放到 gin.context
		token:=strings.ReplaceAll(c.GetHeader("Authorization"),"Bearer ","")
		//fmt.Println("token",token)
		appMatched:=false
		if token!=""{
			claims,err:=pkg.JwtDecode(token)
			if err!=nil{
				ResponseError(c, 2002, err)
				c.Abort()
				return
			}
			//fmt.Println("claims.Issuer",claims.Issuer)
			appList:=dao.AppManagerHandler.GetAppList()
			for _,appInfo:=range appList{
				if appInfo.AppID==claims.Issuer{
					c.Set("app",appInfo)
					appMatched = true
					break
				}
			}
		}
		if serviceDetail.AccessControl.OpenAuth==1 && !appMatched{
			ResponseError(c, 2003, errors.New("not match valid app"))
			c.Abort()
			return
		}
		c.Next()
	}
}