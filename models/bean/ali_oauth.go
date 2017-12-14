package bean

import (
	"encoding/json"
	"github.com/smartwalle/alipay"
)

type AliOauthClient struct {
	Ali *alipay.AliPay
}

type AliOauth struct {
	GrantType string `json:"grant_type"` // 值为authorization_code时，代表用code换取；值为refresh_token时，代表用refresh_token换取
	Code      string `json:"code"`       // 用户对应用授权后得到，即第二步中开发者获取到的auth_code值
}

type AliOauthResponse struct {
	AlipaySystemOauthTokenResponse struct {
		AccessToken     string `json:"access_token"`
		AlipayUserId    string `json:"alipay_user_id"`
		UserId          string `json:"user_id"`
		ExpiresIn       int    `json:"expires_in"`
		ReExpiresIn     int    `json:"re_expires_in"`
		AppRefreshToken int    `json:"app_refresh_token"`
	} `json:"alipay_system_oauth_token_response"`
	Sign string `json:"sign"`
}

// TradePreCreate https://doc.open.alipay.com/docs/api.htm?spm=a219a.7395905.0.0.EnCSXC&docType=4&apiId=862
func (this *AliOauthClient) Oauth(param AliOauth) (results *AliOauthResponse, err error) {
	err = this.Ali.DoRequest("POST", param, &results)
	return results, err

}
func (this AliOauth) APIName() string {
	return "alipay.system.oauth.token"
}

func (this AliOauth) Params() map[string]string {
	var m = make(map[string]string)
	m["grant_type"] = this.GrantType
	m["code"] = this.Code
	return m
}

func (this AliOauth) ExtJSONParamName() string {
	return "biz_content"
}

func (this AliOauth) ExtJSONParamValue() string {
	var bytes, err = json.Marshal(this)
	if err != nil {
		return ""
	}
	return string(bytes)
}

//func (this AliOauth) GetGrantType() string {
//	return this.GrantType
//}
