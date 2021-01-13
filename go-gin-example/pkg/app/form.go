package app

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"net/http"
)

// BindAndValid binds and validates data
func BindAndValid(c *gin.Context, form interface{}) (int, int) {
	err := c.Bind(form)
	if err != nil {
		return http.StatusBadRequest, INVALID_PARAMS
	}

	valid := validation.Validation{}
	check, err := valid.Valid(form)
	if err != nil {
		return http.StatusInternalServerError, ERROR
	}
	if !check {
		MarkErrors(valid.Errors)
		return http.StatusBadRequest, INVALID_PARAMS
	}

	return http.StatusOK, SUCCESS
}