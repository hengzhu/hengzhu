package controllers

import (
	"encoding/json"
	"hengzhu/models"
	"strconv"
	"hengzhu/tool"
)

// TypesController operations for Types
type TypesController struct {
	BaseController
}

func (self *TypesController) List() {
	self.Data["pageTitle"] = "类型列表"
	self.display()
}

func (c *TypesController) Table() {
	filter, _ := tool.BuildFilter(c.Controller, 20)

	filter.Order = []string{"asc"}
	filter.Sortby = []string{"id"}

	ss := []models.Types{}
	total, _ := tool.GetAllByFilterWithTotal(new(models.Types), &ss, filter)

	c.ajaxList("成功", MSG_OK, total, ss)
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
