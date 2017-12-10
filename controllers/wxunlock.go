package controllers

import (
	"hengzhu/tool/payment"
	"github.com/astaxie/beego"
	"hengzhu/models"
)

type WxUnlockController struct {
	BaseController
}

func (c *WxUnlockController) GetCode() {
	state := ""
	wxauth2 := payment.WXOAuth2Authorize{
		RedirectURI: "",    //这里只用指定这个其它的默认就行
		State:       state, //这个用来标识自己的会话
	}
	redirectUrl := wxauth2.ToURL()
	beego.Debug("[WxUnlock] redirect to:", redirectUrl)
	c.Redirect(redirectUrl, 302)
}

func (c *WxUnlockController) GetOpenId() {
	code := c.Input().Get("code")
	wxastoken := payment.WXOAuth2AccessTokenRequest{
		Code: code,
	}
	res, err := wxastoken.Get()
	if err != nil {
		beego.Error("[WxUnlock] GetOpenId err in wxastoken.Get()")
	}
	beego.Debug("[WxUnlock]: GetOpenId get:", res.OpenId)
	detailId := models.GetWxUserInUseByOpenId(res.OpenId)
	if detailId == 0 {
		//没有机器被该用户占用
	}
	//然后可以根据detailId开锁
}
