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
MIICWwIBAAKBgQC9dLK9/ajiJkj4S/CIqfG0/SdafVNs6ldQFIcWb+ow5XOi+H5w
ZiNxgKHeOgtoSoWZtWdCOnxGIdgLlqnjHFgIbqkhETavC7NUWH82qbtMqJUmBjCg
1yGP91JeWHM6OfqnSWB09LgNQjuDRMOtKaEzb2Toa+8q/6O/ftPcUF3dZwIDAQAB
AoGADzD3UCKx0whs23QDYoH1/qQ57piUAuy7eZFbz6HDro4Heq7gPJUEDIra793J
omAvXEbec8IKyvjVwQAguTRBnrAL1mJViCsB5ZAXc30TC0x4YKUlOUcTV0Ncy/H4
B2pQcfddjD5E7GZXoqY4SHBFwX/shURzYSFyM1x3TKDT0eECQQDuISCRB15tZsB8
oYugahMDPDzz3ymacsu9CY/xjvhoWTiwVWc8ZR39wDgz4Ycri72S+M95xO6Y7sj/
GL4hXjCxAkEAy6x4CkwhShEd607azPtF2TJF/sophW77Pe9l96BtKLEnSZyidRak
FibFsJ/17pU0HFYFHsHC9kFe0Kn3KAK1lwJAAXBZzgaJX4fbaeVf/pwleUOH6sFS
cwh2irHgGMmQXrELUqVxdj/2Km5a6JVYR78UairutgGmn23x8PipTXJQQQJALsM3
gG3ASugpLWiadevPOrH/PiOeauNzTeIUUEmGJoyeD5ml9younGNkikv/xDp/j230
mP41zCJwKYqMk6QjmwJAZPAGVJINvOKvETdFEvSJH12aOxMh5lFy5LGYnrpyOTyz
LWIdIY6Ipdn+jXG88I7Un8b9AN6Mk1m3Mt36Af/7mg==
-----END RSA PRIVATE KEY-----
`,
		ALIPAY_PUBLIC_KEY: `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCnxj/9qwVfgoUh/y2W89L6BkRA
FljhNhgPdyPuBV64bfQNN1PjbCzkIM6qRdKBoLPXmKKMiFYnkd6rAoprih3/PrQE
B/VsW8OoM8fxn67UDYuyBTqA23MML9q1+ilIZwBC2AQ2UBVOrFXfFl75p6/B5Ksi
NG9zpgmLCUYuLkxpLQIDAQAB
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
