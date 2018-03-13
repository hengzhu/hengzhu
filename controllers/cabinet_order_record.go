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
	"fmt"
	"crypto/md5"
	"encoding/hex"
	"strings"
	"net/url"
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

var CabinetTimeStamp = make(map[string]int)

const (
	Wx_Pay    = 1 //微信
	Al_Pay    = 2 //支付宝
	FirstIn   = 1 //存付款
	FirstOut  = 2 //取付款
	ForTime   = "fortime"
	NoForTime = "nofortime"
)

// URLMapping ...
func (c *OrderController) URLMapping() {

}

// @Title Post
// @Description 预下单
// @Param	pay_type		query 	int	true		"1.微信 ,2.支付宝"
// @Param	cabinet_id		query 	int	true		"上报的柜子id"
// @Param	timestamp		query 	int	true		"时间戳"
// @Param	total_fee		query 	string	true		"先存后付跳转参数:支付金额"
// @Param	openid		query 	string	true		"先存后付跳转参数:用户标识"
// @Success 201 {int}
// @Failure 403 body is empty
// @router /ReOrder [post]
func (c *OrderController) ReOrder() {
	var flag bool //后下单标志
	var action_type int
	var cd *models.CabinetDetail
	var err, err2 error
	var price float64
	cabinet_mac := c.GetString("cabinet_id")
	pay_type, _ := c.GetInt8("pay_type")
	timestamp, _ := c.GetInt("timestamp", 0)
	if timestamp == 0 {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = "支付参数错误"
		c.ServeJSON()
		return
	}
	CabinetTimeStamp[cabinet_mac] = timestamp

	t_fee := c.GetString("total_fee")
	open_id := c.GetString("open_id")
	if pay_type != Al_Pay && pay_type != Wx_Pay {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = "支付参数错误"
		c.ServeJSON()
		return
	}
	cab, _ := models.GetCabinetByMac(cabinet_mac)
	t, _ := models.GetTypeById(cab.TypeId)
	cabinet_id := cab.Id
	//查找计费类型
	if t_fee == "" && open_id == "" {
		//去oauth校验是否重复存(计次先付后存)
		if t.ChargeMode == 1 && t.TollTime == 1 {
			if pay_type == Wx_Pay {
				c.Ctx.Output.SetStatus(201)
				c.Ctx.WriteString(beego.AppConfig.String("wx_oauth_url") + strconv.Itoa(cabinet_id) + "_" + strconv.Itoa(FirstIn) + "#wechat_redirect")
				return
			}
			c.Ctx.Output.SetStatus(201)
			c.Ctx.WriteString(beego.AppConfig.String("ali_oauth_url") + strconv.Itoa(cabinet_id) + "_" + strconv.Itoa(FirstIn))
			return
		}
	}
	//免费charge_mode为3
	if t.ChargeMode == 3 {
		if err != nil {
			beego.Warn(err)
			c.Data["json"] = "服务器异常"
			c.ServeJSON()
			return
		}
		if pay_type == Wx_Pay {
			//获取code,重定向到微信授权回调
			c.Ctx.Output.SetStatus(201)
			c.Ctx.WriteString(beego.AppConfig.String("wx_oauth_url") + strconv.Itoa(cabinet_id) + "#wechat_redirect")
			return
		}
		c.Ctx.Output.SetStatus(201)
		c.Ctx.WriteString(beego.AppConfig.String("ali_oauth_url") + strconv.Itoa(cabinet_id))
		return
	} else {
		action_type = t.TollTime
	}

	price = t.Price
	//如果为计时收费
	if t.ChargeMode == 2 {
		action_type = 2
	}
	//先存后付类型取物时带过来的参数
	if t_fee != "" && open_id != "" {
		price, _ = strconv.ParseFloat(t_fee, 64)
		cd, err = models.GetCabinetDetailByOpenId(open_id, cabinet_id)
		if err != nil && err != orm.ErrNoRows {
			beego.Error(err.Error())
			return
		}
		flag = true
		if err == orm.ErrNoRows || cd == nil {
			goto A
		}
		goto B
	}
	//先存后付
	if action_type == 2 {
		if pay_type == Wx_Pay {
			//获取code,重定向到微信授权回调
			c.Ctx.Output.SetStatus(201)
			c.Ctx.WriteString(beego.AppConfig.String("wx_oauth_url") + strconv.Itoa(cabinet_id) + "#wechat_redirect")
			return
		}
		c.Ctx.Output.SetStatus(201)
		c.Ctx.WriteString(beego.AppConfig.String("ali_oauth_url") + strconv.Itoa(cabinet_id))
		return
	}
A:
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
		c.NewOrder(cabinet_id, price, open_id, cd, flag)
		c.ServeJSON()
		return
	}
	//创建订单号
	order_no, _ := models.CreateOrderNo()

	//alipay预下单
	b_pri := []byte(pri)
	b_pub := []byte(pub)
	var client = a.New(beego.AppConfig.String("APPID"), beego.AppConfig.String("alipay_partner"), b_pub, b_pri, true)
	var p = a.AliPayTradePreCreate{}
	//加密是rsa1
	client.SignType = a.K_SIGN_TYPE_RSA

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
	err = addOrder(cd.Id, pay_type, price, order_no)
	if err != nil {
		beego.Warn(err)
		c.Data["json"] = "服务器异常"
		c.ServeJSON()
		return
	}
	//省略添加失败再重新请求
	if flag {
		//重定向调起支付
		c.redirect(result.AliPayPreCreateResponse.QRCode)
	}
	return
}

func (c *OrderController) NewOrder(cid int, fee float64, open_id string, cd *models.CabinetDetail, flag bool) {
	var trade_type = payment.TRADE_TYPE_NATIVE
	var err error
	var cabdetail *models.CabinetDetail
	if flag {
		trade_type = payment.TRADE_TYPE_JSAPI
	}
	fee = fee * 100.00
	total_fee := strconv.FormatFloat(fee, 'f', 0, 64)
	if open_id == "" {
		//非取物时下单
		cabdetail, err = models.GetFreeDoorByCabinetId(cid) //根据用户当前扫码的柜子获得一个空闲的门
		if err == orm.ErrNoRows {
			c.Ctx.Output.SetStatus(404)
			c.Data["json"] = errors.New("没有空闲的门可分配").Error()
			c.ServeJSON()
			return
		}
		if err != nil {
			c.Ctx.Output.SetStatus(500)
			c.Data["json"] = errors.New("服务器崩溃").Error()
			c.ServeJSON()
			return
		}
	} else {
		cabdetail, err = models.GetCabinetDetailByOpenId(open_id, cid)
		if err != nil {
			beego.Error(err.Error())
			return
		}
		if cabdetail == nil {
			cabdetail = cd
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
		SpBillCreateIp: "116.62.167.76",
		//*必填 终端ip地址
		TimeStart: "",                                      // 选填 交易起始时间
		TimeExpire: "",                                     // 选填 交易结束时间
		GoodsTag: "",                                       // 选填 订单优惠标记
		NotifyURL: beego.AppConfig.String("wx_notify_url"), //*必填 支付结果通知地址 非常重要
		TradeType: trade_type,                              //*必填 交易类型 这里应为native 扫码支付
		ProductId: strconv.Itoa(cabdetail.CabinetId),       //*必填 商品id原本为选填,但是在扫码支付下必须填写
		LimitPay: "",                                       // 选填 限定支付方式
		OpenId: open_id,                                    // 选填 在扫码支付的情况下不用填写
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
	beego.Debug("[WxPay]: NewOrder success and code:", res.PrePayId)
	if flag {
		//	微信重定向调起支付
		queryStr := fmt.Sprintf("appId=%s&nonceStr=%s&package=%s&signType=%s&timeStamp=%s&key=%s", beego.AppConfig.String("WxAPPID"), nonstr, "prepay_id="+res.PrePayId, "MD5", strconv.Itoa(int(time.Now().Unix())), beego.AppConfig.String("WxKey"))
		hash := md5.New()
		hash.Write([]byte(queryStr))
		cipherStr := hash.Sum(nil)

		queryStr = fmt.Sprintf("appId=%s&nonceStr=%s&package=%s&signType=%s&timeStamp=%s&key=%s", beego.AppConfig.String("WxAPPID"), nonstr, url.QueryEscape("prepay_id="+res.PrePayId), "MD5", strconv.Itoa(int(time.Now().Unix())), beego.AppConfig.String("WxKey"))
		queryStr = queryStr + "&paySign=" + strings.ToUpper(hex.EncodeToString(cipherStr))

		c.redirect(beego.AppConfig.String("domain") + "/middle/payingPage.html?" + queryStr)
	}
	return
}

// @Title Get
// @Description 取物
// @Param	pay_type		query 	int	true		"取物扫码方式：1.微信 ,2.支付宝"
// @Param	cabinet_id		query 	string	true		"上报的柜子id"
// @Param	timestamp		query 	int	true		"时间戳"
// @Success 201 {int}
// @Failure 403 body is empty
// @router /takeout [get]
func (c *OrderController) TakeOut() {
	var str = strconv.Itoa(FirstOut)
	cabinet_mac := c.GetString("cabinet_id")
	pay_type, _ := c.GetInt8("pay_type")
	if pay_type != Al_Pay && pay_type != Wx_Pay {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = "支付参数错误"
		c.ServeJSON()
		return
	}

	timestamp, _ := c.GetInt("timestamp", 0)
	if timestamp == 0 {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = "支付参数错误"
		c.ServeJSON()
		return
	}
	CabinetTimeStamp[cabinet_mac] = timestamp

	cab, _ := models.GetCabinetByMac(cabinet_mac)
	cabinet_id := cab.Id
	t, _ := models.GetTypeById(cab.TypeId)
	//如果为先存后付
	if t.TollTime == 2 {
		str = NoForTime
	}
	//如果为计时付费
	if t.ChargeMode == 2 {
		str = ForTime
	}
	if pay_type == Wx_Pay {
		//获取code,重定向到微信授权回调
		c.Ctx.Output.SetStatus(201)
		c.Ctx.WriteString(beego.AppConfig.String("wx_oauth_url") + strconv.Itoa(cabinet_id) + "_" + str + "#wechat_redirect")
		return
	}
	c.Ctx.Output.SetStatus(201)
	c.Ctx.WriteString(beego.AppConfig.String("ali_oauth_url") + strconv.Itoa(cabinet_id) + "_" + str)
	return
}

func addOrder(cdid int, pay_type int8, price float64, order_no string) (err error) {
	v := models.CabinetOrderRecord{
		CabinetDetailId: cdid,
		PayType:         pay_type,
		Fee:             price,
		CreateDate:      int(time.Now().Unix()),
		OrderNo:         order_no,
	}
	if _, err = models.AddCabinetOrderRecord(&v); err != nil {
		return
	}
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
	config.APP_SECRET = beego.AppConfig.String("APPSECRET")
	payment.InitWXKey(config)

}
