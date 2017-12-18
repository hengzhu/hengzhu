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
	"reflect"
)

var Queues = make(map[string]string)
var Rabbit *mq.Rabbit

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
	//是否初始化柜子
	InitInfo := reflect.ValueOf(result.InitInfo)
	if len(InitInfo.Bytes()) > 0 {
		typ := models.GetDefaultType()
		c := models.Cabinet{
			CabinetID:  result.InitInfo.CabinetID,
			Number:     result.InitInfo.Number,
			Desc:       result.InitInfo.Desc,
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
			LastTime:   time.Now(),
			TypeId:     typ.Id,
		}
		//Doors:      result.InitInfo.Doors,
		//OnUse:      result.InitInfo.OnUse,
		//Close:      result.InitInfo.Close,
		models.AddCabinet(&c)
		for i := 0; i < result.InitInfo.Doors; i++ {
			cd := &models.CabinetDetail{
				Door:      i + 1,
				OpenState: 1,
				Using:     1,
				UseState:  1,
			}
			models.AddCabinetDetail(cd)
		}
		return nil
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
