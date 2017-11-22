package controllers

import (
	m "hengzhu/models/usermodel"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"hengzhu/admin/src/lib"
	"github.com/deepzz0/go-com/log"
	"hengzhu/cache"
	"hengzhu/service"
	"time"
	"hengzhu/ucenter/uc_client"
	"github.com/astaxie/beego"
	"hengzhu/admin/src/rbac"
)

// UserController oprations for User
type UserController struct {
	BaseController
}

// URLMapping ...
func (c *UserController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
	c.Mapping("Logout", c.Logout)
}

func (c *UserController) Login() {
	log.Printf("url is %v", c.Ctx.Request.RequestURI)
	isajax := c.GetString("isajax")
	if isajax == "1" {
		log.Printf("Requ is %v", string(c.Ctx.Input.RequestBody))
		obj := m.UserLogin{}
		err := c.ParseForm(&obj);
		if err != nil {
			log.Printf("ParseForm error: %v", err)
		}
		log.Printf("obj is %v", obj)
		tp, err := m.CheckUserNameType(obj.UserName)
		isuid := "0"
		if tp == "mail" {
			isuid = "2"
		} else if tp == "mobile" {
			mobuser := m.GetUserByMulti(obj.UserName, "mobile")
			obj.UserName = mobuser.UserName
		}
		res, err := uc_client.UcUserLogin(obj, isuid, c.Ctx.Request.UserAgent())
		log.Infof("res %+v", res)
		if res.Result <= 0 {
			//检查Ucenter返回的登录的状态
			switch res.Result {
			case 0:
				c.Response(FAIL, "")
				return
			case -1:
				c.Response(USER_NOT_EXIST, "")
				return
			case -2:
				c.Response(PASSWORD_ERROR, "")
				return
			}
		}

		member, err := CheckLogin(res.UserName, res.Password, "username");
		if err != nil {
			if member.Id == 0 {
				//log.Infof("test:%v",lib.Strtomd5("test123"))
				member, err = m.SyncUser(res.Result, res.UserName, res.Password, res.Email);
				if err != nil {
					c.Response(FAIL, "")
					return
				}
			} else {
				log.Infof("error: %v", err)
				me := m.User{
					Id:      res.Result,
					Password:res.Password,
					Mail:    res.Email,
				}
				id, err := m.UpdateUser(&me);
				if err == nil && id > 0 {
					log.Infof("UpdateUser id: %v", id)
				}
			}
		}
		member = m.CheckSyncUserEmail(res.Email, member)

		sess := c.StartSession()
		defer sess.SessionRelease(c.Ctx.ResponseWriter)
		sess.Set("memberinfo", member)
		sid := sess.SessionID()
		log.Infof("obj.AutoLogin:%v", obj.AutoLogin)

		var sessiontime int
		if obj.AutoLogin == true {
			sessiontime = 1296000
		} else {
			sessiontime = 0
		}

		c.Ctx.Output.Cookie(beego.AppConfig.String("sessionname"), sid, sessiontime)

		log.Printf("member login success: %v", member)
		c.Response(SUCCESS, uc_client.UcUserSynlogin(member))
	}

	memberinfo := c.GetSession("memberinfo")
	if memberinfo != nil {
		uc_client.UcUserSynlogin(memberinfo.(m.User))
		c.Ctx.Redirect(302, "/account")
	}

}

//@Desc: 注册
//
//@Param:
//
//@Return:
//
func (c *UserController) Register() {
	log.Printf("url is %v", c.Ctx.Request.RequestURI)
	isajax := c.GetString("isajax")
	if isajax == "1" {
		log.Print("Requ Isajax")
		log.Printf("Requ is %v", string(c.Ctx.Input.RequestBody))
		//log.Printf("c.Ctx.Request is %v", c.Ctx.Request)
		if !rbac.Cpt.VerifyReq(c.Ctx.Request) {
			log.Infof("captcha error")
			c.Response(CAPTCHA_ERROR, "")
			return
		}
		var obj m.UserReg
		c.ParseForm(&obj)
		log.Infof("obj: %+v", obj)
		c.ValidData(obj.UserLogin)
		if obj.Mobile != "" {
			c.ValidData(obj)
			//if !service.VcodeVerify(obj.Mobile, obj.SmsCode) {
			//	c.Response(SMS_CODE_ERROR, "")
			//	return
			//}
			if m.CheckUserExist("mobile", obj.Mobile) {
				c.Response(MOBILE_ALREADY_REGISTER, "")
				return
			}
		}
		log.Printf("agant:%s", c.Ctx.Request.UserAgent())
		//注册到Ucenter，返回code
		code, err := uc_client.UcUserRegister(obj, c.Ctx.Request.UserAgent());
		if err != nil {
			log.Infof("uc_client.UcUserRegister error:%v", err)
			c.Response(FAIL, "")
			return
		}
		if code < 0 {
			if code == -3 || m.CheckUserExist("user_name", obj.UserName) || m.CheckUserExist("nick_name", obj.UserName) {
				c.Response(USER_ALREADY_EXIST, "")
				return
			} else {
				c.ResponseMsg = "用户名不合法"
				c.Response(FAIL, "");
				return
			}
		} else {
			var usertype, accountmark int8 = 1, 0
			member := m.User{
				Id:        code,
				UserName:  obj.UserName,
				Mobile:    obj.Mobile,
				Password:  lib.Pwdhash(lib.Pwdhash(obj.Password)),
				CreateTime:uint(time.Now().Unix()),
				UserType:  usertype,
				NickName:  obj.UserName,
				//Initialpwd:lib.Pwdhash(lib.Pwdhash(obj.Password)),
				AccountMark:accountmark,
			}
			id, err := m.AddUser(&member);
			if err == nil && id > 0 {
				infoid, err := m.InitUserInfo(int(id), member.NickName, member.Mobile, uint(time.Now().Unix()));
				if err == nil && infoid > 0 {
					c.Response(SUCCESS, "")
					return
				}
			}
			c.Response(FAIL, "")
			return
		}
		return
	} else {
		return
	}
}

//@Desc: 发送验证码接口,限定一个号码一分钟只能发一次
//
//@Param:
//
//@Return:
//
func (c *UserController) SendVCode() {
	log.Printf("Requ is %v", string(c.Ctx.Input.RequestBody))
	log.Printf("Requ is %v", c.Ctx.Request.Body)
	sendtype := c.GetString("type")
	if sendtype == "sms" {
		var obj m.UserBindmobile
		obj.Seccode = "000000"
		c.ParseForm(&obj)
		log.Infof("obj: %+v", obj)
		c.ValidData(obj)
		log.Infof("obj: %+v", obj)
		if cache.Bm.IsExist(SMS_SENDED_CACHE_KEY + obj.Mobile) {
			log.Warnf("phonenum %v too many requests", obj.Mobile)
			c.Response(TOO_MANY_REQUEST, "")
			return
		}
		content, err := service.GetSendVcodeContent(obj.Mobile, service.DEFAULT_VCODE_LEN)
		if err != nil {
			log.Infof("GetSendVcodeContent error: %v", err)
			c.Response(FAIL, "")
			return
		}
		c.Response(SUCCESS, "")
		err = service.SendSms(obj.Mobile, content)
		if err != nil {
			log.Infof("SendSms error: %v", err)
			c.Response(FAIL, "")
			return
		}
		cache.Bm.Put(SMS_SENDED_CACHE_KEY + obj.Mobile, obj.Mobile, SMS_SEND_TIMEOUT)
		c.Response(SUCCESS, "")
		log.Infof("mobile: %v", obj.Mobile)
	} else {
		var obj m.UserEmailBind
		c.ParseForm(&obj)
		log.Infof("obj: %+v", obj)
		obj.Vcode = "000000"
		c.ValidData(obj)
		log.Infof("obj: %+v", obj)
		if cache.Bm.IsExist(SMS_SENDED_CACHE_KEY + obj.Email) {
			log.Warnf("email %v too many requests", obj.Email)
			c.Response(TOO_MANY_REQUEST, "")
			return
		}
		err := service.SendMail(obj.Email);
		if err != nil {
			c.Response(FAIL, "")
			return
		}
		cache.Bm.Put(SMS_SENDED_CACHE_KEY + obj.Email, obj.Email, SMS_SEND_TIMEOUT)
		c.Response(SUCCESS, "")

	}
}

func CheckLogin(username string, password string, ty string) (member m.User, err error) {
	member = m.GetUserByMulti(username, ty)
	log.Infof("member: %+v", member)
	if member.Id == 0 {
		return member, errors.New(codeDesc[USER_NOT_EXIST])
	}
	if member.Password != lib.Pwdhash(lib.Pwdhash(password)) {
		return member, errors.New(codeDesc[PASSWORD_ERROR])
	}
	return member, nil
}

// Logout ...
// @Title Post
// @Description create User
// @Param	body		body 	models.User	true		"body for User content"
// @Success 201 {int} models.User
// @Failure 403 body is empty
// @router /logout [get]
func (c *UserController) Logout() {
	c.DelSession("memberinfo")
	res := uc_client.UcUserSynlogout()
	res += `<script>window.setTimeout(function(){ location.href = "/login"; },500);</script>`
	c.Ctx.WriteString(res)
}

// Post ...
// @Title Post
// @Description create User
// @Param	body		body 	models.User	true		"body for User content"
// @Success 201 {int} models.User
// @Failure 403 body is empty
// @router / [post]
func (c *UserController) Post() {
	var v m.User
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if _, err := m.AddUser(&v); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = v
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get User by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.User
// @Failure 403 :id is empty
// @router /logint/:id [get]
func (c *UserController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := m.GetUserById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get User
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.User
// @Failure 403
// @router / [get]
func (c *UserController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Data["json"] = errors.New("Error: invalid query key/value pair")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := m.GetAllUser(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the User
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.User	true		"body for User content"
// @Success 200 {object} models.User
// @Failure 403 :id is not int
// @router /:id [put]
func (c *UserController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := m.User{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := m.UpdateUserById(&v); err == nil {
			c.Data["json"] = "OK"
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the User
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *UserController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := m.DeleteUser(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
