package controllers

import (
	"github.com/astaxie/beego"
	"hengzhu/admin/src/lib"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/validation"
	"bytes"
	"github.com/Zeniubius/golang_utils/glog"
	"time"
)

const (
	SMS_SENDED_CACHE_KEY = "vocode_sended_"
	SMS_SEND_TIMEOUT = 1*(time.Minute)
)

type BaseController struct {
	beego.Controller
	ResponseMsg string
}

//var cpta *captcha.Captcha
//func init() {
//	store := cache.NewMemoryCache()
//	cpta = captcha.NewWithFilter("/captcha/", store)
//	cpta.ChallengeNums = 4
//	cpta.StdWidth = 100
//	cpta.StdHeight = 40
//}

func (c *BaseController) GetTemplatetype() string {
	templatetype := beego.AppConfig.String("template_type")
	if templatetype == "" {
		templatetype = "website"
	}
	return templatetype
}

func (c *BaseController) RequestMethod() string {
	return c.Ctx.Request.Method
}

func (c *BaseController)Response(code int, data interface{})  {
	if c.ResponseMsg != "" {
		codeDesc[code] = c.ResponseMsg
		c.ResponseMsg = ""
	}
	c.Data["json"] = &map[string]interface{}{
		"status":code,
		"msg":codeDesc[code],
		"data":data,
	}
	c.ServeJSON()
}


func (c *BaseController) ValidData(param interface{}) {
	valid := validation.Validation{}
	b, err := valid.Valid(param)
	if err != nil {
		glog.Info("valid error %v ",err)
		c.Response(FAIL,err.Error())
		c.StopRun()
	}

	if !b {
		var verror bytes.Buffer
		for _, err := range valid.Errors {
			verror.WriteString(err.Key+":"+err.Message+";")
			glog.Info("verror: %v",verror.String())
		}
		c.Response(PARAMS_ERROR,"")
		c.StopRun()
	}
}



func (c *BaseController)DecodeFormData(obj interface{})  {
	err := c.ParseForm(&obj); if err != nil {
		glog.Info("error: %v", err)
	}
}

//@Desc: 判断是否登录
//
//@Param:
//
//@Return:
//
func AccessRegister() {
	var FilterUser = func(ctx *context.Context) {
		memberinfo := ctx.Input.Session("memberinfo")
		if memberinfo==nil && ctx.Request.RequestURI != "/login" {
			ctx.Redirect(302, "/login")
		}
	}
	beego.InsertFilter("/account", beego.BeforeRouter, FilterUser)
}


func init() {
	AccessRegister()
}

const (
	SUCCESS = iota
	FAIL
	ILLEGAL_REQUEST
	PARAMS_ERROR
	LOGIN_FAIL
	USER_NOT_EXIST
	USER_ALREADY_EXIST
	MOBILE_ALREADY_REGISTER
	PASSWORD_ERROR
	SMS_CODE_ERROR
	CAPTCHA_ERROR
	LOGIN_EXPIRE
	DIFFERENCE_NEWPASSWORD
	SAME_NEW_OLD_PASSWORD
	EMAIL_ALREADY_REGISTER
	TOO_MANY_REQUEST

	TEST
)

var codeDesc = map[int]string{
	SUCCESS:      "操作成功",
	FAIL:         "操作失败",
	ILLEGAL_REQUEST:"非法请求",
	PARAMS_ERROR: "参数不符合要求",
	LOGIN_FAIL:   "登录失败",
	USER_NOT_EXIST: "用户不存在",
	USER_ALREADY_EXIST: "用户已存在",
	MOBILE_ALREADY_REGISTER:"手机号码已注册",
	PASSWORD_ERROR: "密码错误",
	SMS_CODE_ERROR: "短信验证码错误",
	CAPTCHA_ERROR:"验证码错误",
	LOGIN_EXPIRE: "登录过期",
	DIFFERENCE_NEWPASSWORD:"两次的密码不同",
	SAME_NEW_OLD_PASSWORD:"新旧密码相同",
	EMAIL_ALREADY_REGISTER:"邮箱已被注册",
	TOO_MANY_REQUEST:"操作太过频繁",

	TEST: lib.Strtomd5("test123"),
}




