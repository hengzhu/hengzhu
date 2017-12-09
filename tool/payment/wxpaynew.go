package payment

import (
	"gopkg.in/chanxuehong/wechat.v2/mch/core"
	"net/http"
)

const (
	WxAppID = ""
	MchId   = ""
	ApiKey  = ""
)

var client = http.DefaultClient //这里的client可以自己定义

func NewWxClient() *core.Client {
	return core.NewClient(WxAppID, MchId, ApiKey, client)
}
