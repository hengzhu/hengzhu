package controllers

import (
	"hengzhu/models"
	"github.com/astaxie/beego"
	"github.com/guidao/gopay/client"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	a "github.com/smartwalle/alipay"
	"time"
	"github.com/astaxie/beego/orm"
	"errors"
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
	Wx_Pay = 1 //微信
	Al_Pay = 2 //支付宝
)

// URLMapping ...
func (c *OrderController) URLMapping() {

}

// @Title Post
// @Description 预下单
// @Param	pay_type		query 	int	true		"1.微信 ,2.支付宝"
// @Param	cabinet_id		query 	int	true		"1.微信 ,2.支付宝"
// @Success 201 {int}
// @Failure 403 body is empty
// @router /ReOrder [post]
func (c *OrderController) ReOrder() {
	var v models.CabinetOrderRecord
	pay_type, _ := c.GetInt8("pay_type")
	cabinet_id, _ := c.GetInt("cabinet_id")
	if pay_type != Al_Pay {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = "[支付宝]:请求失败"
		c.ServeJSON()
		return
	}
	cd, err := models.GetFreeDoorByCabinetId(cabinet_id)
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
	//alipay预下单
	b_pri := []byte(pri)
	b_pub := []byte(pub)
	var client = a.New(beego.AppConfig.String("APPID"), beego.AppConfig.String("alipay_partner"), b_pub, b_pri, true)
	//加密是rsa1
	client.SignType = a.K_SIGN_TYPE_RSA
	var p = a.AliPayTradePreCreate{}
	order_no, _ := models.CreateOrderNo()
	p.OutTradeNo = order_no
	p.NotifyURL = beego.AppConfig.String("alipay_notify_url")
	p.Subject = beego.AppConfig.String("ali_subject")
	p.TotalAmount = beego.AppConfig.String("ali_fee")
	//预下单到支付宝服务器
	result, err := client.TradePreCreate(p)

	if err != nil || !result.IsSuccess() {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = "[支付宝]:网络错误"
		c.ServeJSON()
		return
	}

	v = models.CabinetOrderRecord{
		CabinetDetailId: cd.Id,
		PayType:         pay_type,
		Fee:             50,
		CreateDate:      int(time.Now().Unix()),
		OrderNo:         order_no,
	}
	if _, err := models.AddCabinetOrderRecord(&v); err == nil {
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

func initClient() {
	//加载商户私钥
	if block, _ := pem.Decode([]byte(pri)); block != nil {
		if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
			panic(err)
		} else {
			PARTNET_PRIVATE_KEY = key
		}
	} else {
		panic("load PARTNET_PRIVATE_KEY failed")
	}
	//加载支付宝公钥
	if block, _ := pem.Decode([]byte(pub)); block != nil {
		if pub, err := x509.ParsePKIXPublicKey(block.Bytes); err != nil {
			panic(err)
		} else {
			ALIPAY_PUBLIC_KEY = pub.(*rsa.PublicKey)
		}
	} else {
		panic("load ALIPAY_PUBLIC_KEY failed")
	}
	client.InitAliWebClient(&client.AliWebClient{
		PartnerID:  "2088821824088465111",
		SellerID:   "12",
		AppID:      beego.AppConfig.String("APPID"),
		PrivateKey: PARTNET_PRIVATE_KEY,
		PublicKey:  ALIPAY_PUBLIC_KEY,
	})
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
}
