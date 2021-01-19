package common

import (
	"github.com/casbin/casbin"
	"github.com/go-redis/redis"
	"github.com/hashicorp/consul/api"
	"github.com/jinzhu/gorm"
)

var (
	RedisClient  *redis.Client
	DB           *gorm.DB
	Enforcer     *casbin.Enforcer
	ConsulClient *api.Client
	SignKey      = "huiepay"
	InvalidParam = "InvalidParamError"
)
