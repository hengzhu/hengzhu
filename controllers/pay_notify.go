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
	"strings"
	"hengzhu/utils"
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
		c.Data["cndata"] = "支付失败"
		c.Data["endata"] = err.Error()
		c.TplName = "resp/resp.html"
		c.Render()
		return
	}
	if noti.TradeStatus != "TRADE_SUCCESS" {
		c.Data["cndata"] = "支付失败"
		c.Data["endata"] = "noti.TradeStatus"
		c.TplName = "resp/resp.html"
		c.Render()
		return
	}

	cd, err := models.UpdateOrderSuccessByNo(noti.TradeNo, noti.OutTradeNo, noti.BuyerId)
	if err != nil {
		if err.Error() == "[重复存物]: "+noti.BuyerId {
			c.Ctx.WriteString("success")
			return
		}
		c.Data["cndata"] = "支付失败"
		c.Data["endata"] = err.Error()
		c.TplName = "resp/resp.html"
		c.Render()
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
	//缓存支付绑定记录(先付)
	if cd.StoreTime == 0 {
		err = utils.Redis.SET(utils.PAY+strconv.Itoa(cd.Id), noti.BuyerId, 0)
		if err != nil {
			beego.Warn("[缓存错误]: ", err.Error())
			c.Data["data"] = "服务器错误"
			c.TplName = "resp/resp.html"
			c.Render()
			return
		}
	}
	cab, _ := models.GetCabinetById(cd.CabinetId)
	rmm := bean.RabbitMqMessage{
		CabinetId: cab.CabinetID,
		Door:      cd.Door,
		UserId:    noti.BuyerId,
		DoorState: OpenDoor,
		Timestamp: CabinetTimeStamp[cab.CabinetID],
	}
	bs, _ := json.Marshal(&rmm)
	err = tool.Rabbit.Publish("cabinet_"+cab.CabinetID, bs)
	if err != nil {
		beego.Error("[rabbitmq err:] ", err.Error())
		c.Data["cndata"] = "支付失败"
		c.Data["endata"] = err.Error()
		c.TplName = "resp/resp.html"
		c.Render()
		return
	}
	c.Ctx.WriteString("success")
}

// @Title 支付宝授权用户信息
// @Description 支付宝授权用户信息
// @router /oauthnotify [post]
func (c *PayNotifyController) OauthNotify() {
	var door_no int
	var cdid int
	var fortime string
	var cabinet_id int
	var prefix string         //缓存前缀
	var syb = false           //取标志
	var free = false          //免费标志
	var isRepeatStore = false //检验先付后存的重复存
	auth_code := c.Ctx.Input.Query("auth_code")
	state := c.Ctx.Input.Query("state")
	if strings.Contains(state, "_") {
		syb = true
	}
	if len(state) >= 1 {
		results := strings.Split(state, "_")
		cabinet_id, _ = strconv.Atoi(results[0])
		if len(results) > 1 {
			fortime = results[1]
		}
	}
	cab, _ := models.GetCabinetById(cabinet_id)
	t, _ := models.GetTypeById(cab.TypeId)
	//免费
	if t.ChargeMode == 3 {
		free = true
	}
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
		c.Data["cndata"] = "授权开门失败"
		c.Data["endata"] = err.Error()
		c.TplName = "resp/resp.html"
		c.Render()
		return
	}
	openid := reults.AlipaySystemOauthTokenResponse.UserId
	//校验是否重复存(计次先付后存)
	if fortime == "1" {
		isRepeatStore = true
		//不是取
		syb = false
		goto C
	}
	//免费(取)
	if free && syb {
		goto C
	}

	//如果为先存后付的取物
	if len(fortime) > 1 {
		_, _, cdid, err = models.GetCabinetAndDoorByUserId(openid, 1)
		if err == orm.ErrNoRows {
			c.Data["cndata"] = "没有找到你的存物记录"
			c.Data["endata"] = "No Record Find For You"
			c.TplName = "resp/resp.html"
			c.Render()
			return
		}
		if err != nil {
			beego.Error(err)
			c.Data["cndata"] = "服务器异常"
			c.Data["endata"] = "Server Exception"
			c.TplName = "resp/resp.html"
			c.Render()
			return
		}
		//计算费用
		cd, _ := models.GetCabinetDetailById(cdid)
		total_fee := ""
		//先存后付非计时
		if fortime == NoForTime {
			total_fee = strconv.FormatFloat(t.Price, 'f', 2, 64)
		} else {
			//计时,不足1分钟按1分钟处理
			diff := int(time.Now().Unix()) - cd.StoreTime
			dis_time := diff / 60
			if (diff)%60 != 0 {
				dis_time = dis_time + 1
			}

			// 不足一个计时单位，按一个计时单位处理
			ratio := dis_time / t.Unit
			if dis_time%t.Unit != 0 {
				ratio = ratio + 1
			}
			fee := float64(ratio) * t.Price
			total_fee = strconv.FormatFloat(fee, 'f', 2, 64)
		}
		//重定向到支付宝付款
		timeStamp := CabinetTimeStamp[cab.CabinetID]
		if timeStamp == 0 {
			timeStamp = int(time.Now().Unix())
		}
		c.redirect("http://cabinet.schengzhu.com/order/reorder?pay_type=2&cabinet_id=" + cab.CabinetID + "&timestamp=" + strconv.Itoa(timeStamp) + "&total_fee=" + total_fee + "&open_id=" + openid)
		return
	}
C:
//先存后付授权开门
	if cabinet_id != 0 {
		//免费
		if !syb && free {
			prefix = utils.FREE
		} else if !syb && !free {
			prefix = utils.NOPAY
		}
		//todo one cabinet just store once
		v, err := models.GetCabinetDetailByOpenId(openid, cabinet_id)
		if !syb && err == nil && v != nil {
			beego.Warn(openid + "[ 重复存物 ]")
			c.Data["cndata"] = "重复存物"
			c.Data["endata"] = "Repeat Store"
			c.TplName = "resp/resp.html"
			c.Render()
			return
		} else if syb && v == nil {
			beego.Warn(openid + "[ 没有找到你的存物记录 ]")
			c.Data["cndata"] = "没有找到你的存物记录"
			c.Data["endata"] = "No Record Find For You"
			c.TplName = "resp/resp.html"
			c.Render()
			return
		} else if isRepeatStore {
			//重定向到支付宝付款
			totalFee := strconv.FormatFloat(t.Price, 'f', 2, 64)
			timeStamp := CabinetTimeStamp[cab.CabinetID]
			if timeStamp == 0 {
				timeStamp = int(time.Now().Unix())
			}
			c.redirect("http://cabinet.schengzhu.com/order/reorder?pay_type=2&cabinet_id=" + cab.CabinetID + "&timestamp=" + strconv.Itoa(timeStamp) + "&total_fee=" + totalFee + "&open_id=" + openid)
			return
		}
		//已经使用柜子
		if err == nil && v != nil {
			beego.Warn(openid + "[ have using ]")

		} else {
			cd, err := models.GetFreeDoorByCabinetId(cabinet_id)
			if err == orm.ErrNoRows {
				c.Ctx.Output.SetStatus(404)
				c.Data["cndata"] = "没有空闲的门可分配"
				c.Data["endata"] = "No Free Doors Can Be Allocated"
				c.TplName = "resp/resp.html"
				c.Render()
				return
			}
			if err != nil {
				c.Ctx.Output.SetStatus(500)
				c.Data["cndata"] = "服务器崩溃"
				c.Data["endata"] = "Server Exception"
				c.TplName = "resp/resp.html"
				c.Render()
				return
			}
			//先绑定openid,上传关门信息时才修改存物时间
			err, door := models.BindOpenIdForCabinetDoor(openid, cd.Id)
			if err != nil {
				c.Ctx.Output.SetStatus(501)
				c.Data["json"] = errors.New("系统错误").Error()
				c.ServeJSON()
				return
			}
			//缓存用户绑定(后付)
			err = utils.Redis.SET(prefix+strconv.Itoa(cd.Id), openid, 0)
			beego.Warn("支付缓存: ", prefix+strconv.Itoa(cd.Id))
			if err != nil {
				beego.Warn("[缓存错误]: ", err.Error())
				c.Data["data"] = "服务器错误"
				c.TplName = "resp/resp.html"
				c.Render()
				return
			}
			cdid = cd.Id
			door_no = door
			goto A
		}
	}
	//根据扫码用户的user_id获取已经支付并正在使用的柜子和门
	_, door_no, cdid, err = models.GetCabinetAndDoorByUserId(openid)
	if err == orm.ErrNoRows {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = errors.New("未使用已经支付的柜子").Error()
		c.ServeJSON()
		return
	}
	//取的时候删除缓存
	err = utils.Redis.DEL(utils.LOCKED + strconv.Itoa(cdid))
	beego.Warn("删除缓存: ", utils.LOCKED+strconv.Itoa(cdid))
	if err != nil {
		beego.Error(err)
	}
	if err != nil {
		beego.Error(err)
		c.Data["cndata"] = "支付失败"
		c.Data["endata"] = err.Error()
		c.TplName = "resp/resp.html"
		c.Render()
		return
	}
A:
//免费(存)
	if free && !syb {
		//创建免费存订单
		order_no, _ := models.CreateOrderNo()
		cor := models.CabinetOrderRecord{
			CustomerId:      openid,
			CabinetDetailId: cdid,
			PayType:         0,
			Fee:             0.00,
			CreateDate:      int(time.Now().Unix()),
			OrderNo:         order_no,
			ThirdOrderNo:    "-1",
			IsPay:           1,
		}
		_, err = models.AddCabinetOrderRecord(&cor)
		if err != nil {
			beego.Error(err)
			c.Data["cndata"] = "存物失败"
			c.Data["endata"] = "Store Fail"
			c.TplName = "resp/resp.html"
			c.Render()
			return
		}
	}
	//添加日志记录
	m := models.Log{
		CabinetDetailId: cdid,
		User:            openid,
		Time:            time.Now(),
		Action:          OpenDoor,
	}
	models.AddLog(&m)
	rmm := bean.RabbitMqMessage{
		CabinetId: cab.CabinetID,
		Door:      door_no,
		UserId:    openid,
		DoorState: OpenDoor,
		Timestamp: CabinetTimeStamp[cab.CabinetID],
	}
	bs, _ := json.Marshal(&rmm)
	//下发开门信息
	err = tool.Rabbit.Publish("cabinet_"+cab.CabinetID, bs)
	if err != nil {
		beego.Error("[rabbitmq err:] ", err.Error())
		c.Data["cndata"] = "开门失败，请联系管理员"
		c.Data["endata"] = err.Error()
		c.TplName = "resp/resp.html"
		c.Render()
		return
	}
	qs := "&free=" + strconv.Itoa(0)
	if free && syb {
		qs = "&free=" + strconv.Itoa(1)
	}
	if fortime == "2" && syb {
		qs = "&free=" + strconv.Itoa(1)
	}
	c.redirect(beego.AppConfig.String("domain") + "/middle/aliOperation.html?door_no=" + strconv.Itoa(door_no) + qs)
}

// eg: transdata=%7B%22transtype%22%3A0%2C%22cporderid%22%3A%22re_4ba3YbGUo1%22%2C%22transid%22%3A%220001191495174433775563781837%22%2C%22pcuserid%22%3A%22263%22%2C%22appid%22%3A%221032017051111958%22%2C%22goodsid%22%3A%22153%22%2C%22feetype%22%3A1%2C%22money%22%3A1%2C%22fact_money%22%3A1%2C%22currency%22%3A%22CHY%22%2C%22result%22%3A1%2C%22transtime%22%3A%2220170519141414%22%2C%22pc_priv_info%22%3A%22%22%2C%22paytype%22%3A%221%22%7D&sign=4047a3826502b339b7f2a55145b99291&signtype=MD5
// @Title 微信支付回调
// @Description 微信支付回调
// @router /wxnotify [post]
func (c *PayNotifyController) WxNotify() {
	notify := payment.WXPayResultNotifyArgs{}
	err := xml.Unmarshal(c.Ctx.Input.RequestBody, &notify)
	if err != nil {
		beego.Error("[WxPay]: PayBack err in Unmarshal:", err)
		//c.Data["xml"] = payment.WXPayResultResponse{ReturnCode: "FAIL", ReturnMsg: "参数格式校验错误"}
		//c.ServeXML()
		c.Data["cndata"] = "参数格式校验错误"
		c.Data["endata"] = "[WxPay]: PayBack err in Unmarshal"
		c.TplName = "resp/resp.html"
		c.Render()
		return
	}
	ok := notify.SignValid()
	if !ok {
		//c.Data["xml"] = payment.WXPayResultResponse{ReturnCode: "FAIL", ReturnMsg: "签名失败"}
		//c.ServeXML()
		c.Data["cndata"] = "签名失败"
		c.Data["endata"] = "Failure Of Signature"
		c.TplName = "resp/resp.html"
		c.Render()
		return
	}
	//go func(){}()这里最好异步处理 需要同步给微信返回结果

	detailId := models.GetCabDIdByOrderNo(notify.OutTradeNo)
	if detailId == 0 {
		//说明这个订单有问题
		beego.Error("[WxPay]: get cabinet by out_order_no fail")
	}
	cd, err := models.UpdateOrderSuccessByNo(notify.TransactionId, notify.OutTradeNo, notify.OpenId)
	if err != nil {
		if err.Error() == "[重复存物]: "+notify.OpenId {
			c.Data["xml"] = payment.WXPayResultResponse{ReturnCode: "SUCCESS", ReturnMsg: ""}
			c.ServeXML()
			return
		}
		c.Data["cndata"] = "支付失败"
		c.Data["endata"] = err.Error()
		c.TplName = "resp/resp.html"
		c.Render()
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
	//缓存支付绑定记录(先付)
	if cd.StoreTime == 0 {
		err = utils.Redis.SET(utils.PAY+strconv.Itoa(cd.Id), notify.OpenId, 0)
		if err != nil {
			beego.Warn("[缓存错误]: ", err.Error())
			c.Data["data"] = "服务器错误"
			c.TplName = "resp/resp.html"
			c.Render()
			return
		}
	}
	cab, _ := models.GetCabinetById(cd.CabinetId)
	rmm := bean.RabbitMqMessage{
		CabinetId: cab.CabinetID,
		Door:      cd.Door,
		UserId:    notify.OpenId,
		DoorState: OpenDoor,
		Timestamp: CabinetTimeStamp[cab.CabinetID],
	}
	bs, _ := json.Marshal(&rmm)
	err = tool.Rabbit.Publish("cabinet_"+cab.CabinetID, bs)
	if err != nil {
		beego.Error("[rabbitmq err:] ", err.Error())
		c.Data["cndata"] = "开门失败，请联系管理员"
		c.Data["endata"] = err.Error()
		c.TplName = "resp/resp.html"
		c.Render()
		//c.Ctx.WriteString(err.Error())
		return
	}
	//c.Data["cndata"] = "支付成功"
	//c.Data["endata"] = "Success"
	//c.TplName = "resp/resp.html"
	//c.Render()
	c.Data["xml"] = payment.WXPayResultResponse{ReturnCode: "SUCCESS", ReturnMsg: ""}
	c.ServeXML()
	return
}

// @Title 微信授权用户信息
// @Description 微信授权用户信息
// @router /wx [post]
func (c *PayNotifyController) Wx() {
	var door_no int
	var cdid int
	var fortime string
	var cabinet_id int
	var prefix string         //缓存前缀
	var syb = false           //取标志
	var free = false          //免费标志
	var isRepeatStore = false //检验先付后存的重复存
	code := c.Input().Get("code")
	state := c.Ctx.Input.Query("state")
	if strings.Contains(state, "_") {
		syb = true
	}
	if len(state) >= 1 {
		results := strings.Split(state, "_")
		cabinet_id, _ = strconv.Atoi(results[0])
		if len(results) > 1 {
			fortime = results[1]
		}
	}
	cab, _ := models.GetCabinetById(cabinet_id)
	t, _ := models.GetTypeById(cab.TypeId)
	//免费
	if t.ChargeMode == 3 {
		free = true
	}
	wxastoken := payment.WXOAuth2AccessTokenRequest{
		Code:      code,
		GrantType: GrantType,
	}
	res, err := wxastoken.Get()
	beego.Warn("[WxUnlock] get code: ", code)
	if err != nil {
		beego.Error("[WxUnlock] GetOpenId err: ", err)
	}
	//校验是否重复存(计次先付后存)
	if fortime == "1" {
		isRepeatStore = true
		//不是取
		syb = false
		goto C
	}
	//免费(取)
	if free && syb {
		goto C
	}
	//如果为先存后付的取物
	if len(fortime) > 1 {
		_, _, cdid, err = models.GetCabinetAndDoorByUserId(res.OpenId, 1)
		if err == orm.ErrNoRows {
			c.Data["cndata"] = "没有找到你的存物记录"
			c.Data["endata"] = "No Record Find For You"
			c.TplName = "resp/resp.html"
			c.Render()
			return
		}
		if err != nil {
			beego.Error(err)
			c.Data["cndata"] = "服务器异常"
			c.Data["endata"] = err.Error()
			c.TplName = "resp/resp.html"
			c.Render()
			return
		}
		//计算费用
		cd, _ := models.GetCabinetDetailById(cdid)
		total_fee := ""
		//先存后付非计时
		if fortime == NoForTime {
			total_fee = strconv.FormatFloat(t.Price, 'f', 2, 64)
		} else {
			//计时,不足1分钟按1分钟处理
			diff := int(time.Now().Unix()) - cd.StoreTime
			dis_time := diff / 60
			if (diff)%60 != 0 {
				dis_time = dis_time + 1
			}

			// 不足一个计时单位，按一个计时单位处理
			ratio := dis_time / t.Unit
			if dis_time%t.Unit != 0 {
				ratio = ratio + 1
			}
			fee := float64(ratio) * t.Price
			total_fee = strconv.FormatFloat(fee, 'f', 2, 64)
		}
		//重定向到微信付款
		timeStamp := CabinetTimeStamp[cab.CabinetID]
		if timeStamp == 0 {
			timeStamp = int(time.Now().Unix())
		}
		c.redirect("http://cabinet.schengzhu.com/order/reorder?pay_type=1&cabinet_id=" + cab.CabinetID + "&timestamp=" + strconv.Itoa(timeStamp) + "&total_fee=" + total_fee + "&open_id=" + res.OpenId)
		return
	}
C:
//先存后付授权开门 || 校验是否重复存(计次先付后存)
	if cabinet_id != 0 {
		//免费
		if !syb && free {
			prefix = utils.FREE
		} else if !syb && !free {
			prefix = utils.NOPAY
		}
		v, err := models.GetCabinetDetailByOpenId(res.OpenId, cabinet_id)
		if !syb && err == nil && v != nil {
			beego.Warn(res.OpenId + "[ 重复存物 ]")
			c.Data["cndata"] = "重复存物"
			c.Data["endata"] = "Repeat Store"
			c.TplName = "resp/resp.html"
			c.Render()
			return
		} else if syb && v == nil {
			beego.Warn(res.OpenId + "[ 没有找到你的存物记录 ]")
			c.Data["cndata"] = "没有找到你的存物记录"
			c.Data["endata"] = "No Record Find For You"
			c.TplName = "resp/resp.html"
			c.Render()
			return

		} else if isRepeatStore {
			//重定向到支付宝付款
			totalFee := strconv.FormatFloat(t.Price, 'f', 2, 64)
			timeStamp := CabinetTimeStamp[cab.CabinetID]
			if timeStamp == 0 {
				timeStamp = int(time.Now().Unix())
			}
			c.redirect("http://cabinet.schengzhu.com/order/reorder?pay_type=1&cabinet_id=" + cab.CabinetID + "&timestamp=" + strconv.Itoa(timeStamp) + "&total_fee=" + totalFee + "&open_id=" + res.OpenId)
			return
		}
		//已经使用柜子
		//v, err := models.GetCabinetDetailByOpenId(res.OpenId)
		if err == nil && v != nil {
			beego.Warn(res.OpenId + "[ have using ]")
		} else {
			cd, err := models.GetFreeDoorByCabinetId(cabinet_id)
			if err == orm.ErrNoRows {
				c.Data["cndata"] = "没有空闲的门可分配"
				c.Data["endata"] = "No Free Doors Can Be Allocated"
				c.TplName = "resp/resp.html"
				c.Render()
				return
			}
			if err != nil {
				c.Ctx.Output.SetStatus(500)
				c.Data["cndata"] = "服务器崩溃"
				c.Data["endata"] = "Server Exception"
				c.TplName = "resp/resp.html"
				c.Render()
				return
			}
			//先绑定openid,上传关门信息时才修改为被占用
			err, door := models.BindOpenIdForCabinetDoor(res.OpenId, cd.Id)
			if err != nil {
				c.Ctx.Output.SetStatus(501)
				c.Data["cndata"] = "服务器崩溃"
				c.Data["endata"] = "Server Exception"
				c.TplName = "resp/resp.html"
				c.Render()
				return
			}
			//缓存用户绑定(后付)
			err = utils.Redis.SET(prefix+strconv.Itoa(cd.Id), res.OpenId, 0)
			beego.Warn("支付缓存: ", prefix+strconv.Itoa(cd.Id))
			if err != nil {
				beego.Warn("[缓存错误]: ", err.Error())
				c.Data["data"] = "服务器错误"
				c.TplName = "resp/resp.html"
				c.Render()
				return
			}
			cdid = cd.Id
			door_no = door
			goto A
		}

	}
	//根据扫码用户的open_id获取已经支付并正在使用的柜子和门
	_, door_no, cdid, err = models.GetCabinetAndDoorByUserId(res.OpenId)
	if err == orm.ErrNoRows {
		c.Ctx.Output.SetStatus(404)
		c.Data["cndata"] = "没有找到你的存物记录"
		c.Data["endata"] = "No Record Find For You"
		c.TplName = "resp/resp.html"
		c.Render()
		return
	}
	//取的时候删除缓存
	err = utils.Redis.DEL(utils.LOCKED + strconv.Itoa(cdid))
	beego.Warn("删除缓存: ", utils.LOCKED+strconv.Itoa(cdid))
	if err != nil {
		beego.Error(err)
	}
	if err != nil {
		beego.Error(err)
		c.Data["cndata"] = "服务器崩溃"
		c.Data["endata"] = "Server Exception"
		c.TplName = "resp/resp.html"
		c.Render()
		return
	}
A:
//免费(存)
	if free && !syb {
		//创建免费存订单
		order_no, _ := models.CreateOrderNo()
		cor := models.CabinetOrderRecord{
			CustomerId:      res.OpenId,
			CabinetDetailId: cdid,
			PayType:         0,
			Fee:             0.00,
			CreateDate:      int(time.Now().Unix()),
			OrderNo:         order_no,
			ThirdOrderNo:    "-1",
			IsPay:           1,
		}
		_, err = models.AddCabinetOrderRecord(&cor)
		if err != nil {
			beego.Error(err)
			c.Data["cndata"] = "存物失败"
			c.Data["endata"] = "Store Fail"
			c.TplName = "resp/resp.html"
			c.Render()
			return
		}
	}
	//添加日志记录
	m := models.Log{
		CabinetDetailId: cdid,
		User:            res.OpenId,
		Time:            time.Now(),
		Action:          OpenDoor,
	}
	models.AddLog(&m)
	rmm := bean.RabbitMqMessage{
		CabinetId: cab.CabinetID,
		Door:      door_no,
		UserId:    res.OpenId,
		DoorState: OpenDoor,
		Timestamp: CabinetTimeStamp[cab.CabinetID],
	}
	bs, _ := json.Marshal(&rmm)
	//下发开门信息
	err = tool.Rabbit.Publish("cabinet_"+cab.CabinetID, bs)
	if err != nil {
		beego.Error("[rabbitmq err:] ", err.Error())
		c.Data["cndata"] = "开门失败，请联系管理员"
		c.Data["endata"] = err.Error()
		c.TplName = "resp/resp.html"
		c.Render()
		return
	}
	qs := "&free=" + strconv.Itoa(0)
	if free && syb {
		qs = "&free=" + strconv.Itoa(1)
	}
	if fortime == "2" && syb {
		qs = "&free=" + strconv.Itoa(1)
	}
	c.redirect(beego.AppConfig.String("domain") + "/middle/wxOperation.html?door_no=" + strconv.Itoa(door_no) + qs)

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
