package utils

import (
	bc "github.com/bysir-zl/bygo/cache"
	"github.com/astaxie/beego"
)

const (
	//prefix
	PAY     = "pay_"
	NOPAY   = "nopay_"
	MANAGER = "manager_"
	LOCKED    = "locked_"
)

var Redis *bc.BRedis

func init() {
	redisHost := beego.AppConfig.String("redis_host")
	if redisHost == "" {
		redisHost = "116.62.167.76:6379"
	}
	Redis = bc.NewRedis(redisHost)
}
