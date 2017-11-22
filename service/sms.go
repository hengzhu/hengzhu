package service

import (
	"github.com/astaxie/beego"
	"github.com/Zeniubius/golang_utils/glog"
	"net/url"
	"io/ioutil"
	"strings"
	"net/http"
	"time"
	"hengzhu/cache"
)

const (
	SMS_CACHE_KEY = "sms_vcode_"
	DEFAULT_VCODE_LEN = 6
	DEFAULT_EXP_TIME = 15*(time.Minute)
)

/**
*@Desc 发送手机验证码
*
*@Param mobile 手机号码， content 内容
*
*@Return error, If there is a mistake
*/
func SendSms(moblie string, content string) error {
	tourl := beego.AppConfig.String("sms_url")
	appkey := beego.AppConfig.String("sms_appkey")
	glog.Info("SendSms: tourl %v, appkey:%v",tourl,appkey)
	v := url.Values{}
	v.Set("mobile",moblie)
	v.Set("content",content)
	v.Set("appkey",appkey)
	senddata := v.Encode()
	glog.Info("SendSms: senddata:%v",senddata)
	//body := ioutil.NopCloser(strings.NewReader(senddata))
	body := strings.NewReader(senddata)
	client := &http.Client{}
	glog.Info("SendSms: body: %v",body)
	req, err := http.NewRequest("POST",tourl,body); if err != nil {
		glog.Info("SendSms: NewRequest error: %v",err)
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;param=value")
	resd, err2 := client.Do(req); if err2 != nil{
		glog.Info("SendSms: Client.Do error: %v",err2)
		return err2
	}
	defer resd.Body.Close()
	data, _ := ioutil.ReadAll(resd.Body)
	glog.Info("SendSms: res data:%v",string(data))
	return nil
}

/**
*@Desc 获取验证码内容
*@Param mobile 手机号，len 长度
*@Return content 内容，an error if have
*/
func GetSendVcodeContent(mobile string, len int) (content string, err error) {
	vcode, err := GetVerifyCode(len); if err != nil {
		return
	}
	content = "绑定手机验证码为" + vcode + "【有才互娱】"

	err = cache.Bm.Put(SMS_CACHE_KEY+mobile,vcode,DEFAULT_EXP_TIME)
	return
}







