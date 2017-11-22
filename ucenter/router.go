package ucenter

import (
	ucapicontroller "hengzhu/ucenter/controllers"
	"github.com/astaxie/beego"
)


func router()  {
	beego.Router("/api/uc",&ucapicontroller.ApiController{},"get,post:Index")
}