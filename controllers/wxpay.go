package controllers

import (
	//"hengzhu/tool/payment"
	//"github.com/astaxie/beego"
	//"encoding/xml"
	//"hengzhu/models"
)

type WxPayController struct {
	BaseController
}
//
//func (c *WxPayController) NewOrder() {
//	cabinetId := c.Input().Get("cabinet_id")
//	ip := c.Ctx.Input.IP() //不知道你们的机制是不是这样获得ip
//	cid, err := strconv.ParseInt(cabinetId, 10, 64)
//	if err != nil {
//		beego.Error("[WxPay] NewOrder err in cabinet to int:", err)
//	}
//	cabdetail := models.GetIdleDoorByCabinetId(cid) //根据用户当前扫码的柜子获得一个空闲的门
//	order_no, _ := models.CreateOrderNo()           //这里最好定义好一个订单生成规则 我们内部通过订单号就可以区分支付渠道最好
//	nonstr := ""
//
//	//根据参数创建一个新的订单并且向微信下单获得微信返回的结果
//	wxOrderReq := payment.WXUnifiedorderRequest{//参数可选,签名可以自动生成
//		AppId: "",            //*必填 注意 这两个参数是在payment/wxpay 一开始就配置的
//		MchId: "",            //*必填 注意 这两个参数是在payment/wxpay 一开始就配置的
//		DeviceInfo: "",       // 选填 设备号
//		NonceStr: nonstr,     //*必填 随机字符串
//		Sign: "",             //*必填 但是post方法会自己生成签名,因此可以不用人工填写
//		Body: "",             //*必填 商品描述
//		Detail: "",           // 选填 商品详情
//		Attach: "",           // 选填 附加数据
//		OutTradeNo: order_no, //*必填 商户系统内部订单号 这个重要
//		FeeType: "",          //*选填 币种
//		TotalFee: "标价规则?统一?", //*必填 商品标价
//		SpBillCreateIp: ip,
//		//不知道你们这里是不是填这个//*必填 终端ip地址
//		TimeStart: "",                          // 选填 交易起始时间
//		TimeExpire: "",                         // 选填 交易结束时间
//		GoodsTag: "",                           // 选填 订单优惠标记
//		NotifyURL: "payback 的回调路由",             //*必填 支付结果通知地址 非常重要
//		TradeType: "NATIVE",                    //*必填 交易类型 这里应为native 扫码支付
//		ProductId: string(cabdetail.CabinetId), //*必填 商品id原本为选填,但是在扫码支付下必须填写
//		LimitPay: "",                           // 选填 限定支付方式
//		OpenId: "",                             // 选填 在扫码支付的情况下不用填写
//	}
//	ok := models.CreateNewWxOrder(wxOrderReq) //创建一个本地订单
//	if !ok {
//		beego.Error("[WxPay]: CreateNewWxOrder fail")
//		//创建一个订单失败
//		return
//	}
//	res, err := wxOrderReq.Post()
//	if err != nil {
//		beego.Error("[WxPay]: NewOrder post err and order:", wxOrderReq)
//		//返回一个失败的结果
//		c.Ctx.WriteString(err.Error())
//		return
//	}
//	ok = res.SignValid() //校验返回结果的签名
//	if !ok {
//		beego.Error("[WxPay]: NewOrder post response sign err,order:", wxOrderReq, "res:", res)
//		//签名错误 此处要返回结果
//		c.Ctx.WriteString("verify sign error")
//		return
//	}
//	if res.ReturnCode != "SUCCESS" { //通信结果
//		beego.Error("[WxPay]: NewOrder post response communication err,order:", wxOrderReq, "res:", res) //通信失败
//		c.Ctx.WriteString("communication error")
//		return
//	}
//	if res.ResultCode != "SUCCESS" { //业务结果 下单成功或者失败
//		beego.Error("[WxPay]: NewOrder post response order fail,order:", wxOrderReq, "res:", res)
//		c.Ctx.WriteString(" NewOrder post response order fail")
//		return
//	}
//	beego.Debug("[WxPay]: NewOrder success and code:", res.CodeURL)
//	//c.TplName = ""
//	c.Data["json"] = res.CodeURL
//	c.ServeJSON()
//}

//func (c *WxPayController) PayBack() {
//	//reqbyte := c.Ctx.Input.CopyBody(1<<13)//这里设置8kb 因为复制微信的实例下来是1kb多所以8kb应该足够
//	notify := payment.WXPayResultNotifyArgs{}
//	//restowx := payment.WXPayResultResponse{}
//	err := xml.Unmarshal(c.Ctx.Input.RequestBody, &notify)
//	if err != nil {
//		beego.Error("[WxPay]: PayBack err in Unmarshal:", err)
//		c.Data["xml"] = payment.WXPayResultResponse{ReturnCode: "FAIL", ReturnMsg: "参数格式校验错误"}
//		return
//	}
//	ok := notify.SignValid()
//	if !ok {
//		c.Data["xml"] = payment.WXPayResultResponse{ReturnCode: "FAIL", ReturnMsg: "签名失败"}
//		return
//	}
//
//	//如果上面的步骤进行正常说明微信返回正常然后接下来就是自己的业务逻辑处理
//	/*
//	处理业务逻辑
//	go func(){}()这里最好异步处理 需要同步给微信返回结果
//	*/
//	detailId := models.GetCabDIdByOrderNo(notify.OutTradeNo)
//	if detailId == 0 {
//		//说明这个订单有问题
//	}
//	ok = models.WxPaySuccess(notify, detailId)
//	if !ok {
//		//处理失败
//	}
//	c.Data["xml"] = payment.WXPayResultResponse{ReturnCode: "SUCCESS", ReturnMsg: ""}
//
//}
