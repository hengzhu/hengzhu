package controllers

import (
	"hengzhu/tool/pay_init"
	"hengzhu/tool/payment"
)

// 支付回调
type PayNotifyController struct {
	BaseController
}

const (
	Ali_Pay = 1 //支付宝
	Bbn_Pay = 2 //微信
)

// @Title 支付宝回调
// @Description 支付宝回调
// @router /ali [post]
func (c *PayNotifyController) AliNotify() {
	ap, err := pay_init.CheckAliPayNotify(c.Ctx.Input.RequestBody)
	if err != nil {
		c.Ctx.WriteString(err.Error())
		return
	}

	switch ap.Type {
	case pay_init.Type_Recharge:
		//err := models.UpdateRechargeSuccessByNo(ap.OutTradeNO, ap.TradeNO, Ali_Pay)
		//if err != nil {
		//	c.Ctx.WriteString(err.Error())
		//	return
		}
	//}

	c.Ctx.WriteString("success")
}

// eg: transdata=%7B%22transtype%22%3A0%2C%22cporderid%22%3A%22re_4ba3YbGUo1%22%2C%22transid%22%3A%220001191495174433775563781837%22%2C%22pcuserid%22%3A%22263%22%2C%22appid%22%3A%221032017051111958%22%2C%22goodsid%22%3A%22153%22%2C%22feetype%22%3A1%2C%22money%22%3A1%2C%22fact_money%22%3A1%2C%22currency%22%3A%22CHY%22%2C%22result%22%3A1%2C%22transtime%22%3A%2220170519141414%22%2C%22pc_priv_info%22%3A%22%22%2C%22paytype%22%3A%221%22%7D&sign=4047a3826502b339b7f2a55145b99291&signtype=MD5
// @Title 微信回调
// @Description 微信回调
// @router /wx [post]
func (c *PayNotifyController) WxNotify() {
	transdata := c.GetString("transdata")
	sign := c.GetString("sign")

	ap, err := pay_init.CheckBbnPayNotify(transdata, sign)
	if err != nil {
		c.Ctx.WriteString(err.Error())
		return
	}

	switch ap.Type {
	case pay_init.Type_Recharge:
		//err := models.UpdateRechargeSuccessByNo(ap.OutTradeNO, ap.TradeNO, Bbn_Pay)
		//if err != nil {
		//	c.Ctx.WriteString(err.Error())
		//	return
		//}
	}

	c.Ctx.WriteString(payment.BbnResponse_Success)
}
