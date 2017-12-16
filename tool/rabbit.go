package tool

import (
	"github.com/bysir-zl/bygo/mq"
	"github.com/astaxie/beego"
	"github.com/streadway/amqp"
	"hengzhu/models/bean"
	"encoding/json"
	"time"
	"hengzhu/models"
	"errors"
	"sync"
)

var Queues = make(map[string]string)
var Rabbit *mq.Rabbit
var WaitGroup *sync.WaitGroup

func init() {
	url := beego.AppConfig.String("rabbitmq_url")
	Rabbit = mq.NewRabbit(url)
}

func handleInfo(msg amqp.Delivery) (error) {
	result := bean.RabbitMqMessage{}
	if len(msg.Body) == 0 {
		return errors.New("没有数据")
	}
	err := json.Unmarshal(msg.Body, &result)
	if err != nil {
		return err
	}
	err = models.HandleCabinetFromHardWare(&result)
	if err != nil {
		return err
	}
	return nil
}

func GetMessageFromHardWare(queues map[string]string) {
	for {
		for _, v := range queues {
			go func(s string) {
				err := Rabbit.Receive(s, handleInfo)
				if err != nil {
					beego.Error(err)
				}
			}(v)
		}

		time.Sleep(time.Second * 2)
	}
}
