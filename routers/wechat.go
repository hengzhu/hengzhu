package routers
//
//import (
//	"github.com/astaxie/beego"
//	"hengzhu/controllers"
//)
//
//func init() {
//	beego.Router("/exp", &controllers.WxPayController{}, "*NewOrder") //统一下单接口,这里需要完成给用户分配一个柜子并且生成一个订单
//
//	beego.Router("/exp", &controllers.WxPayController{}, "*PayBack") // 支付完成回调 打开对用的柜门
//
//	beego.Router("/exp", &controllers.WxUnlockController{}, "*GetCode") //用户扫码来打开对应的柜子
//
//	beego.Router("/exp", &controllers.WxUnlockController{}, "*GetOpenId") //成功获得用户的openid,可以找出用户使用的柜子
//}
