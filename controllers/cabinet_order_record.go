package controllers

import (
	"hengzhu/models"
	"github.com/astaxie/beego"
	"crypto/rsa"
	a "github.com/smartwalle/alipay"
	"time"
	"github.com/astaxie/beego/orm"
	"errors"
	"strconv"
	"hengzhu/tool/payment"
)

// 柜子订单支付
type OrderController struct {
	BaseController
}

var pri, pub string

var (
	PARTNET_PRIVATE_KEY *rsa.PrivateKey = nil
	ALIPAY_PUBLIC_KEY   *rsa.PublicKey  = nil
)

const (
	Wx_Pay    = 1 //微信
	Al_Pay    = 2 //支付宝
	First_In  = 1 //存付款
	First_Out = 2 //取付款
	ForTime   = "fortime"
)

// URLMapping ...
func (c *OrderController) URLMapping() {

}

// @Title Post
// @Description 预下单
// @Param	pay_type		query 	int	true		"1.微信 ,2.支付宝"
// @Param	cabinet_id		query 	int	true		"上报的柜子id"
// @Success 201 {int}
// @Failure 403 body is empty
// @router /ReOrder [post]
func (c *OrderController) ReOrder() {
	var v models.CabinetOrderRecord
	var action_type int
	var cd *models.CabinetDetail
	var err, err2 error
	var price float64
	cabinet_mac := c.GetString("cabinet_id")
	pay_type, _ := c.GetInt8("pay_type")
	t_fee := c.Ctx.Request.FormValue("total_fee")
	open_id := c.Ctx.Request.FormValue("open_id")
	if pay_type != Al_Pay && pay_type != Wx_Pay {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = "支付参数错误"
		c.ServeJSON()
		return
	}
	cab, _ := models.GetCabinetByMac(cabinet_mac)
	cabinet_id := cab.Id
	//查找计费类型
	t, _ := models.GetTypeById(cab.TypeId)
	action_type = t.TollTime
	price = t.Price
	//如果为计时收费
	if t.ChargeMode == 2 {
		action_type = 2
	}
	if t_fee != "" && open_id != "" {
		price, _ = strconv.ParseFloat(t_fee, 64)
		cd, err = models.GetCabinetDetailByOpenId(open_id)
		if err != nil {
			beego.Error(err.Error())
			return
		}
		goto B
	}
	//先存后付
	if action_type == 2 {
		if pay_type == Wx_Pay {
			//获取code,重定向到微信授权回调
			//c.GetCode(cabinet_id)
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = beego.AppConfig.String("wx_oauth_url") + strconv.Itoa(cabinet_id)
			c.ServeJSON()
			return
		}
		c.Ctx.Output.SetStatus(201)
		c.Data["json"] = beego.AppConfig.String("ali_oauth_url") + strconv.Itoa(cabinet_id)
		c.ServeJSON()
		return
	}
	cd, err2 = models.GetFreeDoorByCabinetId(cabinet_id)
	if err2 == orm.ErrNoRows {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = errors.New("没有空闲的门可分配").Error()
		c.ServeJSON()
		return
	}
	if err2 != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = errors.New("服务器崩溃").Error()
		c.ServeJSON()
		return
	}
B:
	if pay_type == Wx_Pay {
		c.NewOrder(cabinet_id, price, open_id)
		return
	}

	order_no, _ := models.CreateOrderNo()

	//alipay预下单
	b_pri := []byte(pri)
	b_pub := []byte(pub)
	var client = a.New(beego.AppConfig.String("APPID"), beego.AppConfig.String("alipay_partner"), b_pub, b_pri, true)
	//加密是rsa1
	client.SignType = a.K_SIGN_TYPE_RSA
	var p = a.AliPayTradePreCreate{}
	p.OutTradeNo = order_no
	p.NotifyURL = beego.AppConfig.String("alipay_notify_url")
	p.Subject = beego.AppConfig.String("ali_subject")
	p.TotalAmount = strconv.FormatFloat(price, 'f', 2, 64)
	//预下单到支付宝服务器
	result, err := client.TradePreCreate(p)

	if err != nil || result.AliPayPreCreateResponse.Code != "10000" {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "[支付宝]:网络错误"
		c.ServeJSON()
		return
	}

	v = models.CabinetOrderRecord{
		CabinetDetailId: cd.Id,
		PayType:         pay_type,
		Fee:             price,
		CreateDate:      int(time.Now().Unix()),
		OrderNo:         order_no,
	}
	if _, err = models.AddCabinetOrderRecord(&v); err == nil {
		c.Ctx.Output.SetStatus(201)
	} else {
		c.Ctx.Output.SetStatus(501)
		c.Data["json"] = "服务器异常"
		beego.Warn(err)
		c.ServeJSON()
		return
	}
	//省略添加失败再重新请求
	c.Data["json"] = result.AliPayPreCreateResponse.QRCode
	c.ServeJSON()
	return
}

func (c *OrderController) NewOrder(cid int, fee float64, open_id string) {
	//cabinetId := c.Input().Get("cabinet_id")
	//ip := c.Ctx.Input.IP() //不知道你们的机制是不是这样获得ip
	//cid, err := strconv.ParseInt(cabinetId, 10, 64)
	//if err != nil {
	//	beego.Error("[WxPay] NewOrder err in cabinet to int:", err)
	//	c.Data["json"] = err.Error()
	//}
	var err error
	var cabdetail *models.CabinetDetail
	fee = fee * 100.00
	total_fee := strconv.FormatFloat(fee, 'f', 2, 64)
	if open_id == "" {
		//非取物时下单
		cabdetail = models.GetIdleDoorByCabinetId(int64(cid)) //根据用户当前扫码的柜子获得一个空闲的门
	} else {
		cabdetail, err = models.GetCabinetDetailByOpenId(open_id)
		if err != nil {
			beego.Error(err.Error())
			return
		}
	}

	order_no, _ := models.CreateOrderNo() //这里最好定义好一个订单生成规则 我们内部通过订单号就可以区分支付渠道最好
	nonstr := order_no
	//根据参数创建一个新的订单并且向微信下单获得微信返回的结果
	wxOrderReq := payment.WXUnifiedorderRequest{//参数可选,签名可以自动生成
		AppId: beego.AppConfig.String("WxAPPID"),  //*必填 注意 这两个参数是在payment/wxpay 一开始就配置的
		MchId: beego.AppConfig.String("WxMCH_ID"), //*必填 注意 这两个参数是在payment/wxpay 一开始就配置的
		DeviceInfo: "",                            // 选填 设备号
		NonceStr: nonstr,                          //*必填 随机字符串
		//Sign: "",                                        //*必填 但是post方法会自己生成签名,因此可以不用人工填写
		//Body: strconv.Itoa(First_In), //*必填 商品描述
		Body: "恒铸-储物柜:" + strconv.Itoa(cabdetail.Door), //*必填 商品描述
		Detail: "",                                     // 选填 商品详情
		Attach: "",                                     // 选填 附加数据
		OutTradeNo: order_no,                           //*必填 商户系统内部订单号 这个重要
		FeeType: "",                                    //*选填 币种
		TotalFee: total_fee,                            //*必填 商品标价
		SpBillCreateIp: "39.108.53.220",
		//不知道你们这里是不是填这个//*必填 终端ip地址
		TimeStart: "",                                      // 选填 交易起始时间
		TimeExpire: "",                                     // 选填 交易结束时间
		GoodsTag: "",                                       // 选填 订单优惠标记
		NotifyURL: beego.AppConfig.String("wx_notify_url"), //*必填 支付结果通知地址 非常重要
		TradeType: "NATIVE",                                //*必填 交易类型 这里应为native 扫码支付
		ProductId: strconv.Itoa(cabdetail.CabinetId),       //*必填 商品id原本为选填,但是在扫码支付下必须填写
		LimitPay: "",                                       // 选填 限定支付方式
		OpenId: "",                                         // 选填 在扫码支付的情况下不用填写
	}
	ok := models.CreateNewWxOrder(wxOrderReq, cabdetail.Id) //创建一个本地订单
	if !ok {
		beego.Error("[WxPay]: CreateNewWxOrder fail")
		//创建一个订单失败
		c.Data["json"] = "[WxPay]: CreateNewWxOrder fail"
		return
	}
	res, err := wxOrderReq.Post()
	if err != nil {
		beego.Error("[WxPay]: NewOrder post err and order:", wxOrderReq, err.Error())
		//返回一个失败的结果
		c.Data["json"] = err.Error()
		return
	}
	ok = res.SignValid() //校验返回结果的签名
	if !ok {
		beego.Error("[WxPay]: NewOrder post response sign err,order:", wxOrderReq, "res:", res)
		//签名错误 此处要返回结果
		c.Data["json"] = "verify sign error"
		return
	}
	if res.ReturnCode != "SUCCESS" { //通信结果
		beego.Error("[WxPay]: NewOrder post response communication err,order:", wxOrderReq, "res:", res) //通信失败
		c.Data["json"] = "communication error"
		return
	}
	if res.ResultCode != "SUCCESS" { //业务结果 下单成功或者失败
		beego.Error("[WxPay]: NewOrder post response order fail,order:", wxOrderReq, "res:", res)
		c.Data["json"] = " NewOrder post response order fail"
		return
	}
	beego.Debug("[WxPay]: NewOrder success and code:", res.CodeURL)
	//c.TplName = ""
	c.Data["json"] = res.CodeURL
	c.ServeJSON()

}

// @Title Get
// @Description 取物
// @Param	pay_type		query 	int	true		"取物扫码方式：1.微信 ,2.支付宝"
// @Param	cabinet_id		query 	int	true		"上报的柜子id"
// @Success 201 {int}
// @Failure 403 body is empty
// @router /TakeOut [get]
func (c *OrderController) TakeOut() {
	var str string
	cabinet_mac := c.GetString("cabinet_id")
	pay_type, _ := c.GetInt8("pay_type")
	if pay_type != Al_Pay && pay_type != Wx_Pay {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = "支付参数错误"
		c.ServeJSON()
		return
	}
	cab, _ := models.GetCabinetByMac(cabinet_mac)
	cabinet_id := cab.Id
	t, _ := models.GetTypeById(cab.TypeId)
	//如果为计时付费
	if t.ChargeMode == 2 {
		str = ForTime
	}
	if pay_type == Wx_Pay {
		//获取code,重定向到微信授权回调
		//c.GetCode(cabinet_id)
		c.Ctx.Output.SetStatus(201)
		c.Data["json"] = beego.AppConfig.String("wx_oauth_url") + strconv.Itoa(cabinet_id) + str
		c.ServeJSON()
		return
	}
	c.Ctx.Output.SetStatus(201)
	c.Data["json"] = beego.AppConfig.String("ali_oauth_url") + strconv.Itoa(cabinet_id) + str
	c.ServeJSON()
	return
}

func init() {
	pri = `-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQDrnmBAGqftFloprbmm3dqPjI3ryVZWqwNFm+UniokVp1U/gU2l
yZNXLOXPUVb9Klje4DzIjtGFCxG2dvHM1u66s63R/rlgiXPaNNRBDEE/J8d+EBmK
m0szQ2Svfon4lVrCVQ7zOnlow71/QI4dBUR8oHENUJrUvJvWukvR5hy0KwIDAQAB
AoGAP5Wv99y5sJu1nUXKsiNw1ghiTF07NYxVB7X4c2FJeVR9BvRIFhN99aqiIf6b
cRq6fPsarC0Okc7Y6trSiir+pVM3EpbwwOG0KK2OoUMJdfipHoV1/NX3ZhrWNAa7
f8y3QBKWrhYTjV12YNfwWrV0YUitc0dALsND28kZ3hNP5pECQQD8y6wrlIJXymc7
ZeV8TLR6izVfe0PkLs+IiFiF5qtfnSrNIc4XBYVL36yDpBjVaBLgKkVMhP++ODd9
kvRsB9UVAkEA7pr2wLZ1EmkPXQb3ojY+C7Xw/l5/DTQC4/5QW3MBl0GPoU3a/O0b
pPk0d3nbV5BaKTUm2B5uB0vtGOqrjQs0PwJBAMKP30sLWeZHmXxVyHIKdz15tvJt
5KrSfFgQ2FD2YB+Oz0piIkQFs7nZxOTsf1CAcUamQf/KvSqiCdNUL1qWDKECQCk9
MU6nel5/N/+NF7m6hEjD3m4oaO8gQSukpcDYhLrewvNPIH08gd2mkLHhps5gjaS3
ogoSYFP0hHsc/B95g0MCQQCZ36tOM9VzeDjpJbXKNDmQmRkE6rcVvxFn6HqyNP6z
81qxGn+fqK4YMt4ZA6Z33H6dQvsMtPbB8H9Cg2xoVDYq
-----END RSA PRIVATE KEY-----
`
	pub = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDrnmBAGqftFloprbmm3dqPjI3r
yVZWqwNFm+UniokVp1U/gU2lyZNXLOXPUVb9Klje4DzIjtGFCxG2dvHM1u66s63R
/rlgiXPaNNRBDEE/J8d+EBmKm0szQ2Svfon4lVrCVQ7zOnlow71/QI4dBUR8oHEN
UJrUvJvWukvR5hy0KwIDAQAB
-----END PUBLIC KEY-----`

	config := payment.WXKeyConfig{}
	config.APP_ID = beego.AppConfig.String("WxAPPID")
	config.MCH_ID = beego.AppConfig.String("WxMCH_ID")
	config.MCH_KEY = beego.AppConfig.String("WxKey")
	payment.InitWXKey(config)

}
