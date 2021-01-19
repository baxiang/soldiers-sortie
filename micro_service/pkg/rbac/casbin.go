package rbac

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/casbin/casbin"

	"userService/pkg/gormadapter"
	"userService/pkg/model"
)

func NewCasbin(fileName string, options *model.Options) *casbin.Enforcer {
	adapter := gormadapter.NewAdapter(
		"mysql",
		fmt.Sprintf("%s:%s@tcp(%s)/", options.User, options.Password, options.Addr),
	)
	e := casbin.NewEnforcer(fileName, adapter)
	return e
}

func Split(value string) (string, int64) {
	res := strings.Split(value, ":")
	if len(res) <= 1 {
		return "", 0
	}
	i, _ := strconv.Atoi(res[1])
	return res[0], int64(i)

}
