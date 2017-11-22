package payment

import (
	"errors"
	"github.com/bysir-zl/bygo/util"
	"github.com/bysir-zl/bygo/util/encoder"
	"strconv"
	"github.com/google/go-querystring/query"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strings"
	"net/url"
)

const (
	zpayHost       = "http://pay.csl2016.cn:8000"
	zpayPlaceOrder = zpayHost + "/createOrder.e"
	zpayQueryOrder = zpayHost + "/queryOrder.e"

	ZpayResponse_Success = 0
)

// H5网页收银台 第三方支付
type Zpay struct {
	ZpayConfig
}

type ZpayConfig struct {
	PartnerId, Key string
	AppId          int
}

type ZpayPlaceOrder struct {
	PartnerId  string `json:"partner_id"` // 在商户后台注册后⾃动⽣成
	AppId      int    `json:"app_id"`
	WapType    int    `json:"wap_type"`     // 支付方式 1微信h5 2支付宝h5 3银联H5 4微信二维码 5微信公众号 6QQ钱包
	Money      string `json:"money"`        // 支付金额
	OutTradeNo string `json:"out_trade_no"` // 商户订单编号
	Subject    string `json:"subject"`      // 商品名称
	Qn         string `json:"qn"`           // 商户渠道代码
	ReturnUrl  string `json:"return_url"`   // 商户服务端接收支付结果通知的地址
	Sign       string `json:"sign"`
}

type ZpayQueryOrder struct {
	PartnerId  string `json:"partner_id" url:"partner_id"` // 在商户后台注册后⾃动⽣成
	AppId      int    `json:"app_id" url:"app_id"`
	OutTradeNo string `json:"out_trade_no" url:"out_trade_no"` // 商户订单编号
	Sign       string `json:"sign" url:"sign"`
}

type ZpayNotify struct {
	Code        int    `json:"code" form:"code"` // 0 交易成功 1 交易失败
	AppId       int    `json:"app_id,omitempty" form:"app_id"`
	PayWay      int    `json:"pay_way,omitempty" form:"pay_way"`             // ⽀付⽅式 1 微信 2 ⽀付宝 3 QQ钱包	9 银联
	OutTradeNo  string `json:"out_trade_no" form:"out_trade_no"`             // 商户订单编号
	InvoiceNo   string `json:"invoice_no,omitempty" form:"invoice_no"`       // 平台订单编号
	UpInvoiceNo string `json:"up_invoice_no,omitempty" form:"up_invoice_no"` // 银⾏或微信⽀付流⽔号
	Money       int    `json:"money,omitempty" form:"money"`                 // 本次交易的金额（请务必严格校验商品金额与交易的金额是否一致）
	Qn          string `json:"qn,omitempty" form:"qn"`                       // ⽀付请求同名参数透传
	Sign        string `json:"sign,omitempty" form:"sign"`
}

type ZpanyResp struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// 下单
func (p *Zpay) PlaceOrder(i *ZpayPlaceOrder) (payurl string, err error) {
	i.PartnerId = p.PartnerId
	i.AppId = p.AppId
	signKv := util.OrderKV{}
	signKv.Add("partner_id", i.PartnerId)
	signKv.Add("app_id", strconv.Itoa(i.AppId))
	signKv.Add("wap_type", strconv.Itoa(i.WapType))
	signKv.Add("money", i.Money)
	signKv.Add("out_trade_no", i.OutTradeNo)
	signKv.Add("subject", url.QueryEscape(i.Subject))
	signKv.Add("qn", i.Qn)
	signKv.Add("return_url", i.ReturnUrl)
	signKv.Sort()
	// 签名
	signStr := signKv.EncodeStringWithoutEscape() + "&key=" + p.Key
	signKv.Add("sign", strings.ToUpper(encoder.Md5String(signStr)))
	querystr := signKv.EncodeStringWithoutEscape()
	return zpayPlaceOrder + "?" + querystr, nil
}

func (p *Zpay) Notify(resp *ZpayNotify) (err error) {
	//r.ParseForm()

	signKv := util.OrderKV{}
	signKv.Add("code", strconv.Itoa(resp.Code))
	signKv.Add("app_id", strconv.Itoa(resp.AppId))
	signKv.Add("pay_way", strconv.Itoa(resp.PayWay))
	signKv.Add("out_trade_no", resp.OutTradeNo)
	signKv.Add("invoice_no", resp.InvoiceNo)
	signKv.Add("up_invoice_no", resp.UpInvoiceNo)
	signKv.Add("money", strconv.Itoa(resp.Money))
	signKv.Add("qn", resp.Qn)
	signKv.Sort()
	sign := resp.Sign
	// 签名
	signStr := signKv.EncodeStringWithoutEscape() + "&key=" + p.Key
	signT := strings.ToUpper(encoder.Md5String(signStr))

	if sign != signT {
		err = errors.New("sign error")
		return
	}
	if resp.AppId != p.AppId {
		err = errors.New("appid  error")
		return
	}
	//判断是否支付
	if resp.Code != ZpayResponse_Success {
		err = errors.New("交易未完成")
		return
	}

	zq := ZpayQueryOrder{}
	zq.PartnerId = p.PartnerId
	zq.AppId = p.AppId
	zq.OutTradeNo = resp.OutTradeNo
	zq.Sign = strings.ToUpper(encoder.Md5String("app_id=" + strconv.Itoa(zq.AppId) + "&out_trade_no=" + zq.OutTradeNo + "&" +
		"partner_id=" + zq.PartnerId + "&key=" + p.Key))
	querystr, _ := query.Values(zq)
	queryUrl := zpayQueryOrder + "?" + querystr.Encode()
	response, err := http.Get(queryUrl)
	if err != nil {
		// handle error
		return
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		// handle error
		return
	}
	zpayrp := ZpanyResp{}
	err = json.Unmarshal(body, &zpayrp);
	if err != nil {
		return
	}
	if zpayrp.Code != strconv.Itoa(ZpayResponse_Success) {
		err = errors.New("交易未完成:" + zpayrp.Message)
		return
	}
	return
}
