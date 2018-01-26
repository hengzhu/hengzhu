package controllers

import (
	"encoding/json"
	"errors"
	"hengzhu/models"
	"strconv"
	"strings"
	"hengzhu/tool"
	"fmt"
	"time"
	"hengzhu/models/bean"
	"github.com/astaxie/beego"
	"hengzhu/models/admin"
)

// CabinetController operations for Cabinet
type CabinetController struct {
	BaseController
}

func (self *CabinetController) List() {
	self.Data["pageTitle"] = "状态列表"
	self.display()
	//self.TplName = "admin/list.html"
}

func (c *CabinetController) Table() {
	filter, _ := tool.BuildFilter(c.Controller, 20)

	filter.Order = []string{"asc"}
	filter.Sortby = []string{"id"}

	ss := []models.Cabinet{}
	total, _ := tool.GetAllByFilterWithTotal(new(models.Cabinet), &ss, filter)
	models.AddOtherInfo(&ss)

	c.ajaxList("成功", MSG_OK, total, ss)
}

// 获取某个柜子的详情
func (c *CabinetController) Detail() {
	id, _ := c.GetInt("id", 0)
	if id == 0 {
		return
	}
	cabinet, _ := models.GetCabinetById(id)
	types := models.GetAllTypes()
	models.AddInfo(cabinet)
	//models.AddDetails(cabinet)
	c.Data["cabinet"] = cabinet
	c.Data["types"] = types
	c.Data["pageTitle"] = "状态详情"
	c.display()
	//c.TplName = "/state/detail.html"
}

//
func (c *CabinetController) Open() {
	id, _ := c.GetInt("id")
	if id == 0 {
		c.ajaxMsg(errors.New("参数错误"), MSG_ERR)
	}
	cabinetID, doorId := models.GetOpenMsg(id)

	rmm := bean.RabbitMqMessage{
		CabinetId: cabinetID,
		Door:      int(doorId),
		UserId:    strconv.Itoa(c.userId),
		DoorState: OpenDoor,
	}
	bs, _ := json.Marshal(&rmm)
	err := tool.Rabbit.Publish("cabinet_"+cabinetID, bs)
	if err != nil {
		beego.Error("[rabbitmq err:] ", err.Error())
		c.Ctx.WriteString(err.Error())
		return
	}

	//detail, _ := models.GetDetail(id, int(doorId))
	//detail, _ := models.GetCabinetDetailById(id)
	//detail.OpenState = 2
	//models.UpdateCabinetDetailById(detail)

	user, _ := admin.GetAdminById(c.userId)
	log := models.Log{
		CabinetDetailId: id,
		Action:          OpenDoor,
		User:            user.RealName,
		Time:            time.Now(),
	}
	models.AddLog(&log)

	c.ajaxMsg("成功", MSG_OK)
}

// 清除柜子的存物状态
func (c *CabinetController) Flush() {
	fmt.Printf("user:%v\n", c)
	fmt.Printf("user:%v\n", c.userId)
}

func (c *CabinetController) AjaxSave() {
	id, _ := c.GetInt("id")
	if id == 0 {
		c.ajaxMsg(errors.New("参数错误"), MSG_ERR)
	}

	cabinet, _ := models.GetCabinetById(id)
	typeId, _ := c.GetInt("type", 1)
	cabinet.TypeId = typeId
	cabinet.Address = c.GetString("address")
	cabinet.Number = c.GetString("number")
	cabinet.Desc = c.GetString("desc")
	cabinet.UpdateTime = time.Now()

	if err := models.UpdateCabinetById(cabinet); err != nil {
		c.ajaxMsg(err.Error(), MSG_ERR)
	}
	c.ajaxMsg("修改成功", MSG_OK)
}

// Post ...
// @Title Post
// @Description create Cabinet
// @Param	body		body 	models.Cabinet	true		"body for Cabinet content"
// @Success 201 {int} models.Cabinet
// @Failure 403 body is empty
// @router / [post]
func (c *CabinetController) Post() {
	var v models.Cabinet
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if _, err := models.AddCabinet(&v); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = v
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get Cabinet by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Cabinet
// @Failure 403 :id is empty
// @router /:id [get]
func (c *CabinetController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetCabinetById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get Cabinet
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Cabinet
// @Failure 403
// @router / [get]
func (c *CabinetController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Data["json"] = errors.New("Error: invalid query key/value pair")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := models.GetAllCabinet(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Cabinet
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Cabinet	true		"body for Cabinet content"
// @Success 200 {object} models.Cabinet
// @Failure 403 :id is not int
// @router /:id [put]
func (c *CabinetController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.Cabinet{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateCabinetById(&v); err == nil {
			c.Data["json"] = "OK"
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Cabinet
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *CabinetController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteCabinet(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
