package controllers

import (
	"encoding/json"
	"errors"
	"hengzhu/models"
	"strconv"
	"strings"
	"hengzhu/admin/src/rbac"
	"hengzhu/tool"
	"hengzhu/models/bean"
	"fmt"
)

// CabinetController operations for Cabinet
type CabinetController struct {
	rbac.CommonController
}

// 获取状态信息
func (c *CabinetController) State() {
	if c.IsAjax() {
		filter, err := tool.BuildFilter(c.Controller, 20)
		if err != nil {
			c.RespJSON(bean.CODE_Params_Err, err.Error())
			return
		}

		filter.Order = []string{"asc"}
		filter.Sortby = []string{"id"}

		ss := []models.Cabinet{}

		total, err := tool.GetAllByFilterWithTotal(new(models.Cabinet), &ss, filter)
		if err != nil {
			c.RespJSON(bean.CODE_Not_Found, err.Error())
			return
		}

		models.AddOtherInfo(&ss)
		//c.RespJSONDataWithTotal(ss, total)
		c.Data["json"] = &map[string]interface{}{"total": total, "rows": ss}
		c.ServeJSON()
		return
	} else {
		c.TplName = c.GetTemplatetype() + "/state/index.tpl"
	}
}

// 获取某个柜子的详情
func (c *CabinetController) Detail() {
	id, _ := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if id == 0 {
		return
	}
	cabinet, _ := models.GetCabinetById(id)
	models.AddInfo(cabinet)
	models.AddDetails(cabinet)
	c.Data["cabinet"] = cabinet
	c.TplName = c.GetTemplatetype() + "/state/detail.tpl"
}

// 清除柜子的存物状态
func (c *CabinetController) Flush() {

	fmt.Printf("user:%v\n", c)
	fmt.Printf("user:%v\n", c.Uid())
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
