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

	// 判断是否需要拓展柜子
	//oldDoors := CabinetDoors[result.CabinetId]
	//if oldDoors == 0 {
	//	cabinet, _ := models.GetCabinetByMac(result.CabinetId)
	//	oldDoors = models.GetTotalDoors(cabinet.Id)
	//	CabinetDoors[result.CabinetId] = oldDoors
	//}
	//if result.Door > oldDoors {
	//	cabinet, _ := models.GetCabinetByMac(result.CabinetId)
	//	for i := CabinetDoors[result.CabinetId] + 1; i <= result.Door; i++ {
	//		cd := &models.CabinetDetail{
	//			CabinetId: int(cabinet.Id),
	//			Door:      i,
	//			OpenState: 1,
	//			Using:     1,
	//			UseState:  1,
	//		}
	//		models.AddCabinetDetail(cd)
	//	}
	//
	//	CabinetDoors[result.CabinetId] = result.Door
	//}

	err = models.HandleCabinetFromHardWare(&result)
	if err != nil {
		return err
	}
	return nil
}

func GetMessageFromHardWare() {
	// 由于各种原因服务器重启后，之前的队列只需要重新启动一次就可以
	queues := models.GetCabinetQueues()
	//for {
	for _, v := range queues {
		//if RabbitStarted[v] == true {
		//	// 该协程已经启动，无需再次启动
		//	continue
		//}
		go func(s string) {
			err := Rabbit.Receive(s, handleInfo)
			if err != nil {
				beego.Error(err)
			}
			//RabbitStarted[s] = true
		}("cabinet_" + v)
	}

	//	time.Sleep(time.Second * 2)
	//}
}

// 初始化柜子相关信息
func GetNewCabinet(name string) {
	err := Rabbit.Receive(name, handleNewInfo)
	if err != nil {
		beego.Error(err)
	}
}

// 处理初始化柜子信息
func handleNewInfo(msg amqp.Delivery) (error) {
	result := bean.CabinetInfo{}
	if len(msg.Body) == 0 {
		return errors.New("没有数据")
	}
	err := json.Unmarshal(msg.Body, &result)
	if err != nil {
		return err
	}

	if result.CabinetID == "" || result.Doors == 0 {
		return errors.New("参数错误")
	}

	// 判断是否已有这个柜子，是否初始化柜子
	// 如果已有，则跳过
	flag := models.CheckIfAdd(result.CabinetID)
	if flag == false {
		return errors.New("already have this cabinet")
	}

	typ := models.GetDefaultType()
	c := models.Cabinet{
		CabinetID:  result.CabinetID,
		Number:     result.Number,
		Desc:       result.Desc,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
		LastTime:   time.Now(),
		TypeId:     typ.Id,
	}
	//Doors:      result.InitInfo.Doors,
	//OnUse:      result.InitInfo.OnUse,
	//Close:      result.InitInfo.Close,
	id, _ := models.AddCabinet(&c)
	for i := 0; i < result.Doors; i++ {
		cd := &models.CabinetDetail{
			CabinetId: int(id),
			Door:      i + 1,
			OpenState: 1,
			Using:     1,
			UseState:  1,
		}
		models.AddCabinetDetail(cd)
	}

	CabinetDoors[result.CabinetID] = result.Doors
	value := "cabinet_" + result.CabinetID
	//Queues[strconv.FormatInt(id, 10) ] = value
	go func(s string) {
		err := Rabbit.Receive(s, handleInfo)
		if err != nil {
			beego.Error(err)
		}
		//RabbitStarted[s] = true
	}(value)

	return nil
}
