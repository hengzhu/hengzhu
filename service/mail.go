package service

import (
	"log"
	"github.com/astaxie/beego"
	"hengzhu/cache"
	"gopkg.in/gomail.v2"
	"crypto/tls"
)

const (
	MAIl_CACHE_KEY = "mail_vcode_"
	DEFAULT_MAIL_SUBJECT = "有才互娱验证"
)

/**
*@Desc 发送邮件验证码
*@Param to 邮件地址
*@Return error if not nil
*/
func SendMail(to string) (err error) {
	user := beego.AppConfig.String("mail_user")
	password := beego.AppConfig.String("mail_pass")
	host := beego.AppConfig.String("mail_host")
	port, err := beego.AppConfig.Int("mail_port"); if err != nil {
		log.Printf("SendMail error get mail port error %v", err)
		return
	}
	content, err := GetSendMailContent(to, DEFAULT_VCODE_LEN); if err != nil {
		return
	}
	m := gomail.NewMessage()
	m.SetHeader("From", user)
	m.SetHeader("To", to)
	m.SetHeader("Subject", DEFAULT_MAIL_SUBJECT)
	m.SetBody("text/plain", content)
	d := gomail.NewDialer(host, port, user, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		log.Printf("SendToMail error: %v", err)
		return err
	}
	return

}

func GetSendMailContent(mail string, len int) (content string, err error) {
	vcode, err := GetVerifyCode(len); if err != nil {
		return
	}
	content = "绑定邮箱验证码为" + vcode + "【有才互娱】"
	err = cache.Bm.Put(MAIl_CACHE_KEY + mail, vcode, DEFAULT_EXP_TIME)
	return
}

//test
func SendToMail(user string, password string, host string, port int, to string, subject, body, mailtype string) (err error) {

	m := gomail.NewMessage()
	m.SetHeader("From", user)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/plain", body)
	d := gomail.NewDialer(host, port, user, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err = d.DialAndSend(m); err != nil {
		log.Printf("error: %v", err)
	}
	return
}
