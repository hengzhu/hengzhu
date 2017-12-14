package tool

import (
	"github.com/bysir-zl/bygo/mq"
	"github.com/bysir-zl/bygo/config"
)

var Rabbit *mq.Rabbit

func init() {
	url := config.GetString("url", "rabbit")
	Rabbit = mq.NewRabbit(url)
}
