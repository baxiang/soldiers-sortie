package util

import "github.com/baxiang/soldiers-sortie/blog/pkg/setting"

func Setup() {
	jwtSecret = []byte(setting.AppSetting.JwtSecret)
}