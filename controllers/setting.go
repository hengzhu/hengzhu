package controllers

import (
	"hengzhu/models"
	"time"
)

// SettingController operations for Setting
type SettingController struct {
	BaseController
}

func (c *SettingController) Get() {
	c.Data["pageTitle"] = "设置"

	setting, _ := models.GetSettingById(1)
	c.Data["setting"] = setting

	c.display("setting/set")
}

func (c *SettingController) AjaxSave() {
	setting := models.Setting{}

	id, _ := c.GetInt("id")
	if id == 0 {
		c.ajaxMsg("参数错误", MSG_ERR)
	}
	setting.Id = id

	log_time, _ := c.GetInt("log_time", 30)
	setting.LogTime = log_time

	customer := c.GetString("customer", "")
	if customer == "" {
		c.ajaxMsg("参数错误", MSG_ERR)
	}
	setting.Customer = customer

	setting.UpdateTime = time.Now().Unix()

	if err := models.UpdateSettingById(&setting); err != nil {
		c.ajaxMsg(err.Error(), MSG_ERR)
	}
	c.ajaxMsg("修改成功", MSG_OK)
}
