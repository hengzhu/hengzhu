package controllers

import (
	"github.com/gorilla/websocket"
	"github.com/astaxie/beego"
)

//长连接操作中心
type WsHubController struct {
	BaseController
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var num, _ = beego.AppConfig.Int("cabinet_num")
//连接池
var connections = make(map[int]*websocket.Conn, num)

// @Title 初始化长连接
// @Description 初始化websocket
// @Param	cabinet_id		query 	string	true		"the id for cabinet"
// @Success 201
// @Failure 403 tcp connect failed
// @router /hub [get]
func (c *WsHubController) InitWsHub() {
	connections = make(map[int]*websocket.Conn, num)
	cabinet_id, _ := c.GetInt("cabinet_id")
	if cabinet_id == 0 {
		c.Data["json"] = "[参数错误]"
		c.Ctx.Output.Status = 400
		c.ServeJSON()
		return
	}
	conn, err := upgrader.Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil)
	if err != nil {
		c.Ctx.Output.Status = 502
		c.Data["json"] = "[建立连接失败:]" + err.Error()
		c.ServeJSON()
		return
	}
	//conn.SetReadDeadline(time.Now().Add(time.Second * 180))
	//conn.SetWriteDeadline(time.Now().Add(time.Second * 180))

	connections[cabinet_id] = conn
	c.Ctx.Output.Status = 200
}
