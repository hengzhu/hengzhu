package pay_init

import (
	"github.com/astaxie/beego"
	"strconv"
	"hengzhu/tool/payment"
	"errors"
)

var (
	aliPayNotifyUrl = ""
	bbnPayNotifyUrl = ""

	bbnPay        *payment.BbnPay
	bbnPayGoodsId = 0
)

type AliPayNotify struct {
	TradeNO, OutTradeNO string
	Amount              float64

	Type string // 类型: recharge
}

type BbnPayNotify struct {
	TradeNO, OutTradeNO string
	Amount              float64

	Type string // 类型: recharge
}

const (
	Type_Recharge = "re_"
)

// 支付宝调起支付准备
func CreateAliPayPayInfo(tradeNo, subject, totalFee, body string) string {
	q := payment.NewAPPayReqForApp()
	q.OutTradeNO = Type_Recharge + tradeNo
	q.Subject = subject
	q.TotalFee = totalFee
	q.Body = body

	q.NotifyURL = aliPayNotifyUrl

	return q.String()
}

func CheckAliPayNotify(request []byte) (data AliPayNotify, err error) {
	q, err := payment.NewAPPayResultNotifyArgs(request)
	if err != nil {
		return
	}
	amount, _ := strconv.ParseFloat(q.TotalFee, 64)
	t := ""
	myTradeNo := q.OutTradeNO[3:]
	switch q.OutTradeNO[:3] {
	case Type_Recharge:
		t = Type_Recharge
	default:
		err = errors.New("error type," + q.OutTradeNO[:3])
		return
	}

	data = AliPayNotify{
		TradeNO:    q.TradeNO,
		Amount:     amount,
		OutTradeNO: myTradeNo,
		Type:       t,
	}

	return
}

// 微信调起支付准备
func CreateBbnPayInfo(GoodsName, PcorderId, PcuserId string, money string) (payInfo string, err error) {
	m, _ := strconv.ParseFloat(money, 64)
	mInt := int(m * 100)

	i := payment.BbnPayPlaceOrder{
		Money:     mInt,
		GoodsId:   bbnPayGoodsId,
		GoodsName: GoodsName,
		PcorderId: Type_Recharge + PcorderId,
		NotifyUrl: bbnPayNotifyUrl,
		PcuserId:  PcuserId,
	}
	payInfo, err = bbnPay.PlaceOrder(&i)
	if err != nil {
		return
	}

	return
}

func CheckBbnPayNotify(data, sign string) (response BbnPayNotify, err error) {
	bp, err := bbnPay.Notify(data, sign)
	if err != nil {
		return
	}

	t := ""
	myTradeNo := bp.Cporderid[3:]
	switch bp.Cporderid[:3] {
	case Type_Recharge:
		t = Type_Recharge
	default:
		err = errors.New("error type," + bp.Cporderid[:3])
		return
	}
	if bp.Result != 1 {
		err = errors.New("error no pay")
		return
	}
	response.TradeNO = bp.Transid
	response.OutTradeNO = myTradeNo
	response.Type = t
	response.Amount = float64(bp.Money) / 100
	return
}

func init() {
	config := payment.APKeyConfig{
		ALIPAY_KEY:   beego.AppConfig.String("alipay_key"),
		PARTNER_ID:   beego.AppConfig.String("alipay_partner"),
		SELLER_EMAIL: beego.AppConfig.String("alipay_seller_email"),
		PARTNET_PRIVATE_KEY: `-----BEGIN RSA PRIVATE KEY-----
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
`,
		ALIPAY_PUBLIC_KEY: `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDrnmBAGqftFloprbmm3dqPjI3r
yVZWqwNFm+UniokVp1U/gU2lyZNXLOXPUVb9Klje4DzIjtGFCxG2dvHM1u66s63R
/rlgiXPaNNRBDEE/J8d+EBmKm0szQ2Svfon4lVrCVQ7zOnlow71/QI4dBUR8oHEN
UJrUvJvWukvR5hy0KwIDAQAB
-----END PUBLIC KEY-----`,
	}
	payment.InitAPKey(config)

	aliPayNotifyUrl = beego.AppConfig.String("alipay_notify_url")
	bbnPayNotifyUrl = beego.AppConfig.String("bbnpay_notify_url")

	bbnPayGoodsId, _ = beego.AppConfig.Int("bbnpay_goodsid")

	key := beego.AppConfig.String("bbnpay_key")
	appid := beego.AppConfig.String("bbnpay_appid")
	c := payment.BbnPayConfig{
		Key:   key,
		AppId: appid,
	}
	bbnPay = payment.NewBbnPay(c)
}
