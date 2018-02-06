package controllers

import (
	"encoding/json"
	"hengzhu/models"
	"strconv"
	"hengzhu/tool"
	"errors"
	"time"
)

// TypesController operations for Types
type TypesController struct {
	BaseController
}

func (c *TypesController) List() {
	c.Data["pageTitle"] = "类型列表"
	c.display()
}

func (c *TypesController) Add() {
	c.Data["pageTitle"] = "增加类型"
	c.display()
}

func (c *TypesController) AjaxSave() {
	types := models.Types{}
	types.Name = c.GetString("name")
	types.Default = 2
	charge_mode, _ := c.GetInt("charge_mode", 1)
	types.ChargeMode = charge_mode
	toll_time, _ := c.GetInt("toll_time", 1)
	if charge_mode == 3 {
		types.TollTime = 0
	} else {
		types.TollTime = toll_time
	}
	price, _ := c.GetFloat("price", 0)
	types.Price = price
	unit, _ := c.GetInt("unit", 0)
	types.Unit = unit
	types.CreateTime = time.Now().Unix()

	if _, err := models.AddType(&types); err != nil {
		c.ajaxMsg(err.Error(), MSG_ERR)
	}
	c.ajaxMsg("添加成功", MSG_OK)
}

func (c *TypesController) Table() {
	filter, _ := tool.BuildFilter(c.Controller, 20)

	filter.Order = []string{"asc"}
	filter.Sortby = []string{"id"}

	ss := []models.Types{}
	total, _ := tool.GetAllByFilterWithTotal(new(models.Types), &ss, filter)
	models.AddTypesInfo(ss)

	c.ajaxList("成功", MSG_OK, total, ss)
}

func (c *TypesController) Default() {
	id, _ := c.GetInt("id")
	if id == 0 {
		c.ajaxMsg("参数错误", MSG_ERR)
	}

	err := models.SetDefault(id)
	if err != nil {
		c.ajaxMsg(err.Error(), MSG_ERR)
	}

	c.ajaxMsg("修改成功", MSG_OK)
}

func (c *TypesController) Delete() {
	id, _ := c.GetInt("id")
	if id == 0 {
		c.ajaxMsg(errors.New("参数错误"), MSG_ERR)
	}

	err := models.DeleteType(id)
	if err != nil {
		c.ajaxMsg(err.Error(), MSG_ERR)
	}

	if err = models.UpdateCabinetType(id); err != nil {
		c.ajaxMsg(errors.New("删除失败"), MSG_ERR)
	}

	c.ajaxMsg("修改成功", MSG_OK)
}

// Post ...
// @Title Post
// @Description create Types
// @Param	body		body 	models.Types	true		"body for Types content"
// @Success 201 {int} models.Types
// @Failure 403 body is empty
// @router / [post]
func (c *TypesController) Post() {
	var v models.Types
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if _, err := models.AddType(&v); err == nil {
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

// Put ...
// @Title Put
// @Description update the Types
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Types	true		"body for Types content"
// @Success 200 {object} models.Types
// @Failure 403 :id is not int
// @router /:id [put]
func (c *TypesController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.Types{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateTypeById(&v); err == nil {
			c.Data["json"] = "OK"
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
