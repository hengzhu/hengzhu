package controllers

import (
	// "github.com/astaxie/beego/orm"
	"hengzhu/admin/src/rbac"
)

// UserController oprations for User
type ServiceController struct {
	rbac.CommonController
}

func (this *ServiceController) Service() {
	// games := GetGameInfo()
	// this.Data["games"] = games
	this.TplName = "servicecontroller/service.tpl"
	// this.TplName = "usercontroller/login.tpl"
}

func (this *ServiceController) Detail() {
	this.TplName = "servicecontroller/detail.tpl"
}
