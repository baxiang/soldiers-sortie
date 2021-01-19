package userservice

import (
	"net/http"
	"net/url"
	"time"
	"userService/pkg/common"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ClaimsFactory func() jwt.Claims

type UserClaims struct {
	UserId string
	jwt.StandardClaims
}

func genToken(userId string, exTime time.Time) (string, error) {
	claims := UserClaims{
		UserId: userId,
	}
	claims.ExpiresAt = exTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(common.SignKey))
}

func UserClaimFactory() jwt.Claims {
	return &UserClaims{}
}

func JwtMiddleware(keyFunc jwt.Keyfunc, method jwt.SigningMethod, newClaims ClaimsFactory) gin.HandlerFunc {
	return func(c *gin.Context) {

		tokenString := c.GetHeader("Authorization")

		token, err := jwt.ParseWithClaims(tokenString, newClaims(), func(token *jwt.Token) (interface{}, error) {
			if token.Method != method {
				return nil, ErrUnexpectedSigningMethod
			}

			return keyFunc(token)
		})
		if err != nil {
			logrus.Errorln(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"err":  "Unauthorized",
				"desc": "auth 认证失败",
			})
			return
		}
		if !token.Valid {
			logrus.Errorln(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"err":  "Unauthorized",
				"desc": "auth 认证失败",
			})
			return
		}

		claims, ok := token.Claims.(*UserClaims)
		if !ok || claims.UserId == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"err":  "Unauthorized",
				"desc": "auth 认证失败",
			})
			return
		}
		c.Request.Form = url.Values{}
		c.Request.Form.Set("userid", claims.UserId)
		c.Next()
	}
}
