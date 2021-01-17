package middleware

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"github.com/baxiang/soldiers-sortie/blog/pkg/app"
	"github.com/baxiang/soldiers-sortie/blog/pkg/util"
)

// JWT is jwt middleware
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}

		code = app.SUCCESS
		token := c.Query("token")
		if token == "" {
			code = app.INVALID_PARAMS
		} else {
			_, err := util.ParseToken(token)
			if err != nil {
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					code = app.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
				default:
					code = app.ERROR_AUTH_CHECK_TOKEN_FAIL
				}
			}
		}

		if code != app.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  app.GetMsg(code),
				"data": data,
			})

			c.Abort()
			return
		}

		c.Next()
	}
}