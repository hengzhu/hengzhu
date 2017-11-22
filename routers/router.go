package routers

import (
	"github.com/astaxie/beego"
	"hengzhu/admin"
	"hengzhu/controllers"
	"hengzhu/ucenter"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	//beego.Include(&controllers.UserController{})
	beego.Router("/login", &controllers.UserController{}, "*:Login")
	beego.Router("/logout", &controllers.UserController{}, "*:Logout")
	beego.Router("/reg", &controllers.UserController{}, "*:Register")
	//beego.Router("/account", &controllers.UserController{}, "*:Index")

	beego.Router("/account", &controllers.AccountController{}, "*:Index")
	beego.Router("/account/changepwd", &controllers.AccountController{}, "*:ChangePwd")
	beego.Router("/account/bind", &controllers.AccountController{}, "*:Bind")
	beego.Router("/account/updatedetail", &controllers.AccountController{}, "*:UpdateDetail")
	beego.Router("/sendvcode", &controllers.UserController{}, "post:SendVCode")
	beego.Router("/service", &controllers.ServiceController{}, "*:Service")
	beego.Router("/detail", &controllers.ServiceController{}, "*:Detail")
	admin.Run()
	ucenter.Run()
}
