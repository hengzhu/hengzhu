package controllers

import (
	"hengzhu/tool/payment"
	"github.com/astaxie/beego"
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
}
