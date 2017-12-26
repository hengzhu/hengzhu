package tool

import (
	"github.com/astaxie/beego"
	"github.com/streadway/amqp"
	"hengzhu/models/bean"
	"encoding/json"
	"time"
	"hengzhu/models"
	"errors"
)

//var Queues = make(map[string]string)
//var RabbitStarted = make(map[string]bool)
var NewCabinet = "new"
var Rabbit *models.Rabbit
var CabinetDoors = make(map[string]int)

func init() {
	url := beego.AppConfig.String("rabbitmq_url")
	Rabbit = models.NewRabbit(url)
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

//func GetMessageFromHardWare() {
//	// 由于各种原因服务器重启后，之前的队列只需要重新启动一次就可以
//	queues := models.GetCabinetQueues()
//	//for {
//	for _, v := range queues {
//		//if RabbitStarted[v] == true {
//		//	// 该协程已经启动，无需再次启动
//		//	continue
//		//}
//		go func(s string) {
//			err := Rabbit.Receive(s, handleInfo)
//			if err != nil {
//				beego.Error(err)
//			}
//			//RabbitStarted[s] = true
//		}("cabinet_" + v)
//	}
//
//	//	time.Sleep(time.Second * 2)
//	//}
//}

func GetMsg(name string) {
	err := Rabbit.Receive(name, handleMsgInfo)
	if err != nil {
		beego.Error(err)
	}
}

// 处理柜子信息
func handleMsgInfo(msg amqp.Delivery) (error) {
	result := bean.RabbitMqMessage{}
	if len(msg.Body) == 0 {
		return errors.New("没有数据")
	}
	err := json.Unmarshal(msg.Body, &result)
	if err != nil {
		return err
	}

	if result.CabinetId == "" {
		return errors.New("参数错误")
	}

	if len(result.DoorStatus) != 0 {
		err = createOrAddCabinet(&result)
		if err != nil {
			return err
		}
	} else {
		err = handleHeartbeat(&result)
		if err != nil {
			return err
		}
	}

	return nil
}

// 初始化或者扩展柜子
func createOrAddCabinet(result *bean.RabbitMqMessage) (err error) {
	flag := models.CheckIfAdd(result.CabinetId)
	if flag == true {
		// 初始化柜子
		typ := models.GetDefaultType()
		c := models.Cabinet{
			CabinetID:  result.CabinetId,
			Desc:       result.Desc,
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
			LastTime:   time.Now(),
			TypeId:     typ.Id,
		}

		id, err := models.AddCabinet(&c)
		if err != nil {
			return err
		}

		// 初始化时默认所有的柜子都可用、空闲且关闭状态
		for i, _ := range result.DoorStatus {
			cd := &models.CabinetDetail{
				CabinetId: int(id),
				Door:      i + 1,
				OpenState: 1,
				Using:     1,
				UseState:  1,
			}
			_, err = models.AddCabinetDetail(cd)
			if err != nil {
				return err
			}
		}

		// 放在缓存中，方便下次快速判断
		CabinetDoors[result.CabinetId] = len(result.DoorStatus)
	} else {
		// 不需要初始化，检测是否需要扩展
		oldDoors := CabinetDoors[result.CabinetId]
		cabinet, _ := models.GetCabinetByMac(result.CabinetId)

		cabinet.LastTime = time.Now()
		err = models.UpdateCabinetById(cabinet)
		if err != nil {
			return err
		}

		if oldDoors == 0 {
			oldDoors = models.GetTotalDoors(cabinet.Id)
			CabinetDoors[result.CabinetId] = oldDoors
		}

		if len(result.DoorStatus) > oldDoors {
			for i := oldDoors + 1; i <= len(result.DoorStatus); i++ {
				cd := &models.CabinetDetail{
					CabinetId: int(cabinet.Id),
					Door:      i,
					OpenState: 1,
					Using:     1,
					UseState:  1,
				}
				_, err = models.AddCabinetDetail(cd)
				if err != nil {
					return err
				}
			}

			// 更新缓存
			CabinetDoors[result.CabinetId] = result.Door
		}
	}
	return nil
}

// 处理心跳信息
func handleHeartbeat(result *bean.RabbitMqMessage) (err error) {
	cabinet, _ := models.GetCabinetByMac(result.CabinetId)

	cabinet.LastTime = time.Now()
	err = models.UpdateCabinetById(cabinet)
	if err != nil {
		return err
	}

	err = models.HandleCabinetFromHardWare(result)
	if err != nil {
		return err
	}

	return nil
}
