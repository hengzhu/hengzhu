package controllers

import (
	"github.com/astaxie/beego"
	"hengzhu/ucenter/uc_client/data/cache"
)

type UcBaseController struct {
	beego.Controller
	app cache.Item
	cache []cache.Item
}




