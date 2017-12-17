package controllers

import (
	"hengzhu/tool/payment"
	"github.com/smartwalle/alipay"
	"github.com/astaxie/beego"
	"hengzhu/models"
	"hengzhu/tool"
	"hengzhu/models/bean"
	"strconv"
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"errors"
	"encoding/xml"
	"time"
)

const (
	GrantType = "authorization_code"
	OpenDoor  = "open"
	CloseDoor = "close"
)

var oauth_pri, oauth_pub string

// 支付回调
type PayNotifyController struct {
	BaseController
}

// @Title 支付宝支付回调
// @Description 支付宝支付回调
// @router /alinotify [post]
func (c *PayNotifyController) AliNotify() {
	b_pri := []byte(pri)
	b_pub := []byte(pub)
	client := alipay.New(beego.AppConfig.String("APPID"), beego.AppConfig.String("alipay_partner"), b_pub, b_pri, true)
	client.SignType = alipay.K_SIGN_TYPE_RSA
	client.AliPayPublicKey = b_pub
	var noti *alipay.TradeNotification
	//忽略验签
	noti, err := client.GetTradeNotification(c.Ctx.Request)
	if err != nil {
		beego.Error(err)
		c.Ctx.WriteString(err.Error())
		return
	}
	if noti.TradeStatus != "TRADE_SUCCESS" {
		c.Ctx.WriteString("noti.TradeStatus")
		return
	}

	cd, err := models.UpdateOrderSuccessByNo(noti.TradeNo, noti.OutTradeNo, noti.BuyerId)
	if err != nil {
		c.Ctx.WriteString(err.Error())
		return
	}
	//添加日志记录
	m := models.Log{
		CabinetDetailId: cd.Id,
		User:            noti.BuyerId,
		Time:            time.Now(),
		Action:          OpenDoor,
	}
	models.AddLog(&m)
	//cor, err := models.GetOrderRecordByOerderNo(noti.OutTradeNo)
	//if err != nil {
	//	c.Ctx.WriteString(err.Error())
	//	return
	//}
	//door_no := []byte{uint8(cd.Door)}
	//err = connections[cabinet.CabinetId].WriteMessage(len([]byte{uint8(cd.Door)}), door_no)
	cab, _ := models.GetCabinetById(cd.CabinetId)
	rmm := bean.RabbitMqMessage{
		CabinetId: cab.CabinetID,
		Door:      cd.Door,
		UserId:    noti.BuyerId,
		DoorState: OpenDoor,
	}
	bs, _ := json.Marshal(&rmm)
	err = tool.Rabbit.Publish("cabinet_"+strconv.Itoa(cd.CabinetId), bs)
	if err != nil {
		beego.Error("[rabbitmq err:] ", err.Error())
		c.Ctx.WriteString(err.Error())
		return
	}
	tool.Queues[strconv.Itoa(cd.CabinetId)] = "cabinet_" + strconv.Itoa(cd.CabinetId)
	c.Ctx.WriteString("success")

}

// @Title 支付宝授权用户信息
// @Description 支付宝授权用户信息
// @router /oauthnotify [post]
func (c *PayNotifyController) OauthNotify() {
	var cid, door_no int
	var cdid int
	auth_code := c.Ctx.Input.Query("auth_code")
	cabinet_id, _ := strconv.Atoi(c.Ctx.Input.Query("state"))

	o_pri := []byte(oauth_pri)
	o_pub := []byte(oauth_pub)
	client := alipay.New(beego.AppConfig.String("APPID"), beego.AppConfig.String("alipay_partner"), o_pub, o_pri, true)
	//忽略验签
	//client.AliPayPublicKey = o_pub
	ao := bean.AliOauthClient{}
	ao.Ali = client

	param := bean.AliOauth{}
	param.GrantType = GrantType
	param.Code = auth_code

	reults, err := ao.Oauth(param)
	if err != nil || reults == nil {
		beego.Error(reults, err)
		c.Ctx.WriteString(err.Error())
		return
	}
	openid := reults.AlipaySystemOauthTokenResponse.UserId
	//先存后付授权开门
	if cabinet_id != 0 {
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
		//先绑定openid,上传关门信息时才修改为被占用
		err, door := models.BindOpenIdForCabinetDoor(openid, cd.Id)
		if err != nil {
			c.Ctx.Output.SetStatus(501)
			c.Data["json"] = errors.New("系统错误").Error()
			c.ServeJSON()
			return
		}
		cdid = cd.Id
		cid = cabinet_id
		door_no = door
		goto A
	}
	//根据扫码用户的user_id获取已经支付并正在使用的柜子和门
	cid, door_no, cdid, err = models.GetCabinetAndDoorByUserId(openid)
	if err == orm.ErrNoRows {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = errors.New("未使用已经支付的柜子").Error()
		c.ServeJSON()
		return
	}
	if err != nil {
		beego.Error(err)
		c.Ctx.WriteString(err.Error())
		return
	}
A:
//添加日志记录
	m := models.Log{
		CabinetDetailId: cdid,
		User:            openid,
		Time:            time.Now(),
		Action:          OpenDoor,
	}
	models.AddLog(&m)
	cab, _ := models.GetCabinetById(cid)
	rmm := bean.RabbitMqMessage{
		CabinetId: cab.CabinetID,
		Door:      door_no,
		UserId:    openid,
		DoorState: OpenDoor,
	}
	bs, _ := json.Marshal(&rmm)
	//下发开门信息
	err = tool.Rabbit.Publish("cabinet_"+strconv.Itoa(cid), bs)
	if err != nil {
		beego.Error("[rabbitmq err:] ", err.Error())
		c.Ctx.WriteString(err.Error())
		return
	}
	tool.Queues[strconv.Itoa(cid)] = "cabinet_" + strconv.Itoa(cid)
	c.Ctx.WriteString("success")

}

// eg: transdata=%7B%22transtype%22%3A0%2C%22cporderid%22%3A%22re_4ba3YbGUo1%22%2C%22transid%22%3A%220001191495174433775563781837%22%2C%22pcuserid%22%3A%22263%22%2C%22appid%22%3A%221032017051111958%22%2C%22goodsid%22%3A%22153%22%2C%22feetype%22%3A1%2C%22money%22%3A1%2C%22fact_money%22%3A1%2C%22currency%22%3A%22CHY%22%2C%22result%22%3A1%2C%22transtime%22%3A%2220170519141414%22%2C%22pc_priv_info%22%3A%22%22%2C%22paytype%22%3A%221%22%7D&sign=4047a3826502b339b7f2a55145b99291&signtype=MD5
// @Title 微信支付回调
// @Description 微信支付回调
// @router /wxnotify [post]
func (c *PayNotifyController) WxNotify() {
	notify := payment.WXPayResultNotifyArgs{}
	beego.Warn("======== ", c.Ctx.Request.Form, c.Ctx.Request.URL)
	err := xml.Unmarshal(c.Ctx.Input.RequestBody, &notify)
	if err != nil {
		beego.Error("[WxPay]: PayBack err in Unmarshal:", err)
		c.Data["xml"] = payment.WXPayResultResponse{ReturnCode: "FAIL", ReturnMsg: "参数格式校验错误"}
		c.ServeXML()
		return
	}
	ok := notify.SignValid()
	if !ok {
		c.Data["xml"] = payment.WXPayResultResponse{ReturnCode: "FAIL", ReturnMsg: "签名失败"}
		c.ServeXML()
		return
	}
	//go func(){}()这里最好异步处理 需要同步给微信返回结果

	detailId := models.GetCabDIdByOrderNo(notify.OutTradeNo)
	if detailId == 0 {
		//说明这个订单有问题
		beego.Error("[WxPay]: get cabinet by out_order_no fail")
	}
	//ok = models.WxPaySuccess(notify, detailId)
	//if !ok {
	//	//处理失败
	//	beego.Error("[WxPay]: WxPay fail")
	//}
	cd, err := models.UpdateOrderSuccessByNo(notify.TransactionId, notify.OutTradeNo, notify.OpenId)
	if err != nil {
		c.Ctx.WriteString(err.Error())
		return
	}
	//添加日志记录
	m := models.Log{
		CabinetDetailId: cd.Id,
		User:            notify.OpenId,
		Time:            time.Now(),
		Action:          OpenDoor,
	}
	models.AddLog(&m)
	cab, _ := models.GetCabinetById(cd.CabinetId)
	rmm := bean.RabbitMqMessage{
		CabinetId: cab.CabinetID,
		Door:      cd.Door,
		UserId:    notify.OpenId,
		DoorState: OpenDoor,
	}
	bs, _ := json.Marshal(&rmm)
	err = tool.Rabbit.Publish("cabinet_"+strconv.Itoa(cd.CabinetId), bs)
	if err != nil {
		beego.Error("[rabbitmq err:] ", err.Error())
		c.Ctx.WriteString(err.Error())
		return
	}
	tool.Queues[strconv.Itoa(cd.CabinetId)] = "cabinet_" + strconv.Itoa(cd.CabinetId)
	c.Data["xml"] = payment.WXPayResultResponse{ReturnCode: "SUCCESS", ReturnMsg: ""}

}

// @Title 微信授权用户信息
// @Description 微信授权用户信息
// @router /wxoauthnotify [post]
func (c *PayNotifyController) WxOauthNotify() {
	var cid, door_no int
	var cdid int
	code := c.Input().Get("code")
	cabinet_id, err := strconv.Atoi(c.Input().Get("state"))
	wxastoken := payment.WXOAuth2AccessTokenRequest{
		Code:      code,
		GrantType: GrantType,
	}
	res, err := wxastoken.Get()
	if err != nil {
		beego.Error("[WxUnlock] GetOpenId err in wxastoken.Get()")
	}
	beego.Debug("[WxUnlock]: GetOpenId get:", res.OpenId)
	//先存后付授权开门
	if cabinet_id != 0 {
		cd, err := models.GetFreeDoorByCabinetId(cabinet_id)
		if err == orm.ErrNoRows {
			c.Ctx.Output.SetStatus(404)
			c.Data["xml"] = errors.New("没有空闲的门可分配").Error()
			c.ServeXML()
			return
		}
		if err != nil {
			c.Ctx.Output.SetStatus(500)
			c.Data["xml"] = errors.New("服务器崩溃").Error()
			c.ServeXML()
			return
		}
		//先绑定openid,上传关门信息时才修改为被占用
		err, door := models.BindOpenIdForCabinetDoor(res.OpenId, cd.Id)
		if err != nil {
			c.Ctx.Output.SetStatus(501)
			c.Data["xml"] = errors.New("系统错误").Error()
			c.ServeXML()
			return
		}
		cdid = cd.Id
		cid = cabinet_id
		door_no = door
		goto A
	}
	//根据扫码用户的open_id获取已经支付并正在使用的柜子和门
	cid, door_no, cdid, err = models.GetCabinetAndDoorByUserId(res.OpenId)
	if err == orm.ErrNoRows {
		c.Ctx.Output.SetStatus(404)
		c.Data["xml"] = errors.New("未使用已经支付的柜子").Error()
		c.ServeXML()
		return
	}
	if err != nil {
		beego.Error(err)
		c.Data["xml"] = err.Error()
		c.ServeXML()
		return
	}
A:
//添加日志记录
	m := models.Log{
		CabinetDetailId: cdid,
		User:            res.OpenId,
		Time:            time.Now(),
		Action:          OpenDoor,
	}
	models.AddLog(&m)
	cab, _ := models.GetCabinetById(cid)
	rmm := bean.RabbitMqMessage{
		CabinetId: cab.CabinetID,
		Door:      door_no,
		UserId:    res.OpenId,
		DoorState: OpenDoor,
	}
	bs, _ := json.Marshal(&rmm)
	//下发开门信息
	err = tool.Rabbit.Publish("cabinet_"+strconv.Itoa(cid), bs)
	if err != nil {
		beego.Error("[rabbitmq err:] ", err.Error())
		c.Data["xml"] = err.Error()
		c.ServeXML()
		return
	}
	tool.Queues[strconv.Itoa(cid)] = "cabinet_" + strconv.Itoa(cid)
	c.Data["xml"] = payment.WXPayResultResponse{ReturnCode: "SUCCESS", ReturnMsg: ""}
	c.ServeXML()

}

func init() {
	oauth_pri = `-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEA2JlNpVXe2cFdY2seEBj2Ksm06Yxz9ADIUXCvcbeLZamY+kZ
T9YkfDtqesOrGJ22u+hOqI25PhbZC5KbkNoDw35YGF3NYeCCKLptAVi2Df0TOwd
aIplWQWPeBB1ixT6rCXmun55Whx7IwaubLkJVWl60j5orjJeyWuBIiZ8noc8mUl
+uOR7mcEntAHw0X1QwvUdOgHoWSot3RbyIJ8Yf7eEW1A++FOa7w9bIfPyM09O0c
ueAaqilk433JvlaXIFaM+bCXoGLfoql9RyGpwa5Y888kQBd19U6c93qw2r+Xn5n
oX5l06U0FQNSfehXvjx3fMJh0NyZN9aaRKDB2PqzDlwIDAQABAoIBAGP4HcY5o+
mNPbUtM2rqmnOVNVK16K6tzccI43Dw7f22EU0yOH4TE6qfbK7rLRn1ndT+ToCb4
UgtnyI5hQtC5+nKLHWWXzbSjfSE42TjDNYow+TjR569zynA0mS5otzKS3uY5J4W
idzJeV9dtoa85oKK/w7g+4X9dHLwq8CLiCYoAfw6pjosoeMAVJNJb4Gxh2OWiHn
ZFgqwHCezoF7WGCknshyO0rJSC8ltUxic9EoGZqLPJc3UJX0keoPO61u0/FwYoV
FKJ4auEyioy2BEkb7k74kfU4DIf+UE+CRVdByoJtQvMF+4ujwF0E9hG+ChAUV6y
CnhnynJGDzqhDAu0YECgYEA+6f1xB9MzEmsKJNlociGJhUAUaZI8+mpc4HqzP84
HjBd+kjhZ72KSDjGcHWwtOm7WZIzCCIm4ihSh3Y0ky1kdWg54aNO8v1RY7Idmsz
fEd+S9mFAfB+IUSfKJWtogTYp6SOBnzh7ISDNREF+BWcKN8QaXSiC/Nfh30QL/y
8E6UkCgYEA3FZt4Ysvb04a7Waz9pkTD320zdK/bkfyS3G4WsjSEsaWDBQU9bs0u
MFP6yGZ7ozUl4qdW/ODDt+leQ3ki2iMu/OYuL+SDjix3SFoGFvK3BKfC2yPXtpY
RD/FbCSuBP5J++Fh7222hVzlmgJZH21eDECViV/wW8j/PokNo3Q+Jd8CgYAn6Xm
HA1fQxpZxUP87a2wrOgV07aSAWryvPxmYLZoe35joCwsEwwDdd3OxfljqOG+oQx
Go5pG4KKD+LvcjqH1YSZF0gcwRqa9w2lzrojZ2xTivrrjldrLN/DuJN8G5THfVK
/Zw5CpTFLq5apGsFa1/LrDnuXcc1rhSCp7EeBaVUQKBgCPQ1Mmt00cXfh8K68Pw
+/0vpN00HbPyc/s5gAsZy7QLncZW2VVcWeSSX8hLzPbO45vCh3Oz8KDRT9eOn5D
drMq8fR3C3h37r0XPsVkMSrxdNocn3WJAwcpOR2wdxj+/ig0shLvjrKCfCh9vtE
b8gyYgtW4AL1TsJjlnE9V3BscnAoGAZ/2tDFEYal2Lvri2c/N0nXsK6Svbh5ecw
esPP0Qxy7lqjGHaF914C68w7ZCUJ44ji67eifD8CQq8SA+x5pPW6Eq1rDEsnKHW
wbEa+BTX5Rn/X9YasqBoIkhA//zdSGclwxR3qSJCKw/RNFClL3yM7ba6omB9zvG
7Z3oA/Ddxcw4=
-----END RSA PRIVATE KEY-----
`
	oauth_pub = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA2JlNpVXe2cFdY2seEBj
2Ksm06Yxz9ADIUXCvcbeLZamY+kZT9YkfDtqesOrGJ22u+hOqI25PhbZC5KbkNo
Dw35YGF3NYeCCKLptAVi2Df0TOwdaIplWQWPeBB1ixT6rCXmun55Whx7IwaubLk
JVWl60j5orjJeyWuBIiZ8noc8mUl+uOR7mcEntAHw0X1QwvUdOgHoWSot3RbyIJ
8Yf7eEW1A++FOa7w9bIfPyM09O0cueAaqilk433JvlaXIFaM+bCXoGLfoql9RyG
pwa5Y888kQBd19U6c93qw2r+Xn5noX5l06U0FQNSfehXvjx3fMJh0NyZN9aaRKD
B2PqzDlwIDAQAB
-----END PUBLIC KEY-----`
}
