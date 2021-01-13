package util

import "github.com/baxiang/soldiers-sortie/go-gin-example/pkg/setting"

func Setup() {
	jwtSecret = []byte(setting.AppSetting.JwtSecret)
}