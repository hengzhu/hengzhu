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

var NewCabinet = "new"
var Rabbit *models.Rabbit

func init() {
	url := beego.AppConfig.String("rabbitmq_url")
	Rabbit = models.NewRabbit(url)
}

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

	beego.Info("msg:", result)
	//if len(result.DoorStatus) != 0 {
	err = createOrUpdateCabinet(&result)
	if err != nil {
		return err
	}
	//} else {
	err = handleHeartbeat(&result)
	if err != nil {
		return err
	}
	//}

	return nil
}

// 初始化或者扩展柜子,每次都要根据上传的状态修改数据库
func createOrUpdateCabinet(result *bean.RabbitMqMessage) (err error) {
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

		// 初始化时默认所有的柜子都空闲,启用
		for i, door := range result.DoorStatus {
			OpenState := 1
			if door.Locked == false {
				OpenState = 2
			}
			WireConnected := 1
			if door.WireConnected == false {
				WireConnected = 2
			}
			cd := &models.CabinetDetail{
				CabinetId:     int(id),
				Door:          i + 1,
				OpenState:     OpenState,
				Using:         1,
				UseState:      1,
				WireConnected: WireConnected,
			}
			_, err = models.AddCabinetWithNull(cd)
			if err != nil {
				return err
			}
		}
	} else {
		// 不需要初始化
		cabinet, _ := models.GetCabinetByMac(result.CabinetId)

		// 判断柜子门数，是否需要删除多余的柜子门
		if len(result.DoorStatus) < models.GetTotalDoors(cabinet.Id) {
			err := models.DeleteDoor(cabinet.Id, len(result.DoorStatus))
			if err != nil {
				beego.Error("wrong to delete door!")
			}
		}

		cabinet.LastTime = time.Now()
		err = models.UpdateCabinetById(cabinet)
		if err != nil {
			return err
		}

		// 循环判断每个门是否需要更新，或者是否需要扩展
		for _, door := range result.DoorStatus {
			cabinetDetail, _ := models.GetCabinetDetail(cabinet.Id, door.Door)

			// 当前数据库没有这个门，需要扩展
			if cabinetDetail == nil {
				OpenState := 1
				if door.Locked == false {
					OpenState = 2
				}
				WireConnected := 1
				if door.WireConnected == false {
					WireConnected = 2
				}
				cd := &models.CabinetDetail{
					CabinetId:     int(cabinet.Id),
					Door:          door.Door,
					OpenState:     OpenState,
					Using:         1,
					UseState:      1,
					WireConnected: WireConnected,
				}
				_, err = models.AddCabinetWithNull(cd)
				if err != nil {
					return err
				}
			} else {
				// 当前数据库有此门，只需要判断是否需要更新状态
				OpenState := 1
				if door.Locked == false {
					OpenState = 2
				}
				WireConnected := 1
				if door.WireConnected == false {
					WireConnected = 2
				}

				// 需要更新门状态,修改数据库
				if cabinetDetail.OpenState != OpenState || cabinetDetail.WireConnected != WireConnected {
					cabinetDetail.OpenState = OpenState
					cabinetDetail.WireConnected = WireConnected
					models.UpdateCabinetDetail(cabinetDetail)
				}
			}
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
	return nil
}
