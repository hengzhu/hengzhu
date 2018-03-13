package controllers

import (
	"encoding/json"
	"errors"
	"hengzhu/models"
	"strconv"
	"strings"
	"time"
	"hengzhu/models/admin"
	"fmt"
	"github.com/astaxie/beego"
	"hengzhu/utils"
)

// CabinetDetailController operations for CabinetDetail
type CabinetDetailController struct {
	BaseController
}

func (c *CabinetDetailController) Table() {
	id, _ := c.GetInt("id", 0)
	if id == 0 {
		return
	}
	details, err := models.GetDetailsByCabinetId(id)
	models.AddAllInfo(details)

	if err != nil {
		c.ajaxList("失败", MSG_ERR, 0, details)
	}
	c.ajaxList("成功", MSG_OK, 0, details)
}

// 更改柜子的启用状态
func (c *CabinetDetailController) ChangeUse() {
	id, _ := c.GetInt("id")
	if id == 0 {
		c.ajaxMsg(errors.New("参数错误"), MSG_ERR)
	}

	cabinetDetail, _ := models.GetCabinetDetailById(id)

	useState, _ := c.GetInt("use")
	if useState != 1 && useState != 2 {
		c.ajaxMsg(errors.New("状态错误"), MSG_ERR)
	}
	cabinetDetail.UseState = useState

	if err := models.ChangeCabinetState(cabinetDetail); err != nil {
		c.ajaxMsg(err.Error(), MSG_ERR)
	}
	c.ajaxMsg("修改成功", MSG_OK)
}

// 清除柜子的储物状态
func (c *CabinetDetailController) Clear() {
	id, _ := c.GetInt("id")
	if id == 0 {
		c.ajaxMsg(errors.New("参数错误"), MSG_ERR)
	}

	cabinetDetail, _ := models.GetCabinetDetailById(id)
	cabinetDetail.Using = 1
	//cabinetDetail.UserID = ""
	cabinetDetail.StoreTime = 0

	if err := models.UpdateCabinetDetailWithNUll(cabinetDetail); err != nil {
		c.ajaxMsg(err.Error(), MSG_ERR)
	}
	beego.Warn("删除缓存:", "*_"+strconv.Itoa(id))
	utils.Redis.DEL(utils.MANAGER + strconv.Itoa(id))
	utils.Redis.DEL(utils.NOPAY + strconv.Itoa(id))
	utils.Redis.DEL(utils.PAY + strconv.Itoa(id))
	utils.Redis.DEL(utils.LOCKED + strconv.Itoa(id))
	utils.Redis.DEL(utils.FREE + strconv.Itoa(id))

	user, _ := admin.GetAdminById(c.userId)

	log := models.Log{}
	log.CabinetDetailId = id
	log.Action = "清除"
	log.Time = time.Now()
	log.User = user.RealName

	models.AddLog(&log)

	c.ajaxMsg("修改成功", MSG_OK)
}

// 获取某个柜子门的历史记录
func (c *CabinetDetailController) Record() {
	id, _ := c.GetInt("id", 0)
	if id == 0 {
		return
	}

	//startTime, _ := c.GetInt64("begin", 0)
	//endTime, _ := c.GetInt64("end", time.Now().Unix())
	begin := c.GetString("begin", "2017-12-01 00:00:00")
	end := c.GetString("end", time.Now().Format("2006-01-02 15:04:05"))

	cabinetDetail, _ := models.GetCabinetDetailById(id)
	models.AddIDInfo(cabinetDetail)

	fmt.Printf("begin:%v---end:%v\n", begin, end)
	models.AddLogInfo(cabinetDetail, begin, end)

	c.Data["record"] = cabinetDetail
	c.Data["pageTitle"] = "历史记录"
	c.TplName = "cabinet/record.html"
	//c.display()
}

// Post ...
// @Title Post
// @Description create CabinetDetail
// @Param	body		body 	models.CabinetDetail	true		"body for CabinetDetail content"
// @Success 201 {int} models.CabinetDetail
// @Failure 403 body is empty
// @router / [post]
func (c *CabinetDetailController) Post() {
	var v models.CabinetDetail
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if _, err := models.AddCabinetDetail(&v); err == nil {
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
// @Description get CabinetDetail by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.CabinetDetail
// @Failure 403 :id is empty
// @router /:id [get]
func (c *CabinetDetailController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetCabinetDetailById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get CabinetDetail
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.CabinetDetail
// @Failure 403
// @router / [get]
func (c *CabinetDetailController) GetAll() {
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

	l, err := models.GetAllCabinetDetail(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the CabinetDetail
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.CabinetDetail	true		"body for CabinetDetail content"
// @Success 200 {object} models.CabinetDetail
// @Failure 403 :id is not int
// @router /:id [put]
func (c *CabinetDetailController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.CabinetDetail{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateCabinetDetailById(&v); err == nil {
			c.Data["json"] = "OK"
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// @Title 更新柜子门的开关状态
// @Description 更新柜子门的开关状态
// @Param	cabinet_id		query 	string	true		"The id for cabinet you want to update"
// @Param	open_state		query 	string	true		""
// @Param	door_no		    query 	string	true		""
// @Success 200 {object}
// @Failure 403
// @router /ChangeOpen [put]
func (c *CabinetDetailController) ChangeOpen() {
	cid, _ := c.GetInt("cabinet_id")
	if cid == 0 {
		c.ajaxMsg(errors.New("参数错误"), MSG_ERR)
	}
	door_no, _ := c.GetInt("door_no")
	if door_no == 0 {
		c.ajaxMsg(errors.New("门号错误"), MSG_ERR)
	}

	cabinetDetail, _ := models.GetCabinetDetail(cid, door_no)

	openState, _ := c.GetInt("open_state")
	if openState != 1 && openState != 2 {
		c.ajaxMsg(errors.New("状态错误"), MSG_ERR)
	}
	cabinetDetail.UseState = openState

	if err := models.UpdateCabinetDetailById(cabinetDetail); err != nil {
		c.ajaxMsg(err.Error(), MSG_ERR)
	}
	err := connections[cid].WriteMessage(4, []byte("ok"))
	if err != nil {
		beego.Error(err)
	}
	c.ajaxMsg("修改成功", MSG_OK)
}

// Delete ...
// @Title Delete
// @Description delete the CabinetDetail
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *CabinetDetailController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteCabinetDetail(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
