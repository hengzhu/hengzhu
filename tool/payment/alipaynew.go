package payment

import (
	"github.com/smartwalle/alipay"
)

type AliKeyConfig struct {
	APPID 				string //appid
	PARTNER_ID          string //商户id
	SELLER_EMAIL        string //商户支付email
	SIGN_TYPE           string //签名类型 RSA
	ALIPAY_KEY          string //阿里支付密钥
	PARTNET_PRIVATE_KEY string //商户私钥
	ALIPAY_PUBLIC_KEY   string //阿里支付公钥
}

func InitAlipay(conf AliKeyConfig) *alipay.AliPay {
	client 	:= alipay.New(conf.APPID,conf.PARTNER_ID,[]byte(conf.ALIPAY_PUBLIC_KEY),[]byte(conf.PARTNET_PRIVATE_KEY),true)
	return client
}

