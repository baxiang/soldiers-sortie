package api

import (
	"github.com/baxiang/soldiers-sortie/go-gin-example/service/auth_service"
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"

	"github.com/baxiang/soldiers-sortie/go-gin-example/pkg/app"
	"github.com/baxiang/soldiers-sortie/go-gin-example/pkg/util"
)

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

// @Summary Get Auth
// @Produce  json
// @Param username query string true "userName"
// @Param password query string true "password"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /auth [get]
func GetAuth(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}

	username := c.PostForm("username")
	password := c.PostForm("password")

	a := auth{Username: username, Password: password}
	ok, _ := valid.Valid(&a)

	if !ok {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, app.INVALID_PARAMS, nil)
		return
	}

	authService := auth_service.Auth{Username: username, Password: password}
	isExist, err := authService.Check()
	if err != nil {
		appG.Response(http.StatusInternalServerError, app.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}

	if !isExist {

		appG.Response(http.StatusUnauthorized, app.ERROR_AUTH, nil)
		return
	}

	token, err := util.GenerateToken(username, password)
	if err != nil {
		appG.Response(http.StatusInternalServerError, app.ERROR_AUTH_TOKEN, nil)
		return
	}

	appG.Response(http.StatusOK, app.SUCCESS, map[string]string{
		"token": token,
	})
}