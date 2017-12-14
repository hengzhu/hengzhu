package tool

import (
	"github.com/bysir-zl/bygo/mq"
	"github.com/astaxie/beego"
)

var Rabbit *mq.Rabbit

func init() {
	url := beego.AppConfig.String("rabbitmq_url")
	Rabbit = mq.NewRabbit(url)
}
