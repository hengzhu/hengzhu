package controllers

import (
	"encoding/json"
	"hengzhu/models"
)

//扫码支付预下单操作
type ReOrderController struct {
	BaseController
}

// @Title Post
// @Description 支付宝预下单
// @Success 201 {int} 'success'
// @Failure 403 body is empty
// @router /ali [post]
func (c *ReOrderController) Post() {
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
