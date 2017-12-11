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
// @Success 201 {int}
// @Failure 403 body is empty
// @router /ReOrder [get]
func (c *OrderController) ReOrder() {
	var v models.CabinetOrderRecord
	pay_type, _ := c.GetInt8("pay_type")
	beego.Warn(pay_type)
	//alipay预下单
	initRSA()

	b_pri := []byte(pri)
	b_pub := []byte(pub)
	var client = a.New(beego.AppConfig.String("APPID"), beego.AppConfig.String("alipay_partner"), b_pub, b_pri, true)
	//加密是rsa1
	client.SignType = a.K_SIGN_TYPE_RSA

	var p = a.AliPayTradePreCreate{}
	order_no, _ := models.CreateOrderNo()

	p.NotifyURL = beego.AppConfig.String("ali_notify")
	p.Subject = beego.AppConfig.String("ali_subject")
	p.OutTradeNo = order_no
	p.TotalAmount = "0.1"
	//预下单到支付宝服务器
	result, err := client.TradePreCreate(p)
	if err != nil {
		c.Data["json"] = "[支付宝]:网络错误"
		c.ServeJSON()
	}
	v = models.CabinetOrderRecord{
		CustomerId:   "1",
		PayType:      pay_type,
		ThirdOrderNo: "",
		Fee:          50,
		CreateDate:   int(time.Now().Unix()),
	}
	if _, err := models.AddCabinetOrderRecord(&v); err == nil {
		c.Ctx.Output.SetStatus(201)
	} else {
		c.Data["json"] = err.Error()
		c.ServeJSON()
	}
	_, err = models.AddCabinetOrderRecord(&v)
	if err != nil {
		c.Data["json"] = "网络错误"
		c.ServeJSON()
	}
	//省略添加失败再重新请求
	beego.Warn("+++++++++++: ", result, "\n", result.AliPayPreCreateResponse.QRCode)
	c.Data["json"] = result.AliPayPreCreateResponse.QRCode
	c.ServeJSON()
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

func initRSA() {
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
