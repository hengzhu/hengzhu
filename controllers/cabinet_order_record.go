package controllers

import (
	"encoding/json"
	"hengzhu/models"
)

// CabinetOrderRecordController operations for CabinetOrderRecord
type CabinetOrderRecordController struct {
	BaseController
}

// URLMapping ...
func (c *CabinetOrderRecordController) URLMapping() {
}

// @Title Post
// @Description 预下单
// @Param	body		body 	models.CabinetOrderRecord	true		"body for CabinetOrderRecord content"
// @Success 201 {int} models.CabinetOrderRecord
// @Failure 403 body is empty
// @router / [post]
func (c *CabinetOrderRecordController) Post() {
	var v models.CabinetOrderRecord
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if _, err := models.AddCabinetOrderRecord(&v); err == nil {
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
