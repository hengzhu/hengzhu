package controllers

import (
	"github.com/Zeniubius/golang_utils/glog"
	m "hengzhu/models/usermodel"
	"time"
	"fmt"
	"hengzhu/ucenter/uc_client"
)

type AccountController struct {
	BaseController
}

func (c *AccountController) URLMapping() {
	c.Mapping("Post", c.Post)
}

//@Desc: 账号中心首页
//
//@Param:
//
//@Return:
//
func (c *AccountController)Index() {
	glog.Info("requs:%v", c.Ctx.Request.RequestURI)
	membersess := c.GetSession("memberinfo")
	if membersess == nil {
		c.Response(LOGIN_EXPIRE, "")
		c.Ctx.Redirect(302, "/login")
	}
	c.Data["islogin"] = true
	member := membersess.(m.User)
	//id := member.Id
	glog.Info("Member:%v",member)
	stime, err := c.GetInt("stime")
	etime, err := c.GetInt("etime")
	if stime != 0 && etime != 0 {
		//uint(member.Id) //test: uint(2837)
		//Orders, err := m.QueryOrderByTime(uint(2837), stime, etime); if err != nil {
			//todo 改回用户自己的id
		//	c.Response(FAIL, "")
		//	return
		//}
		//c.Response(SUCCESS, Orders)
		return
	} else {
		//Orders, err := cache.GetOutputG5Order(uint(2837)); if err != nil {
		//	log.Infof("output error: %v", err)
		//	c.Response(FAIL, "")
		//	return
		//}
		//log.Infof("user orders nums %v", len(Orders.([]m.OutputG5Order)))
		//c.Data["orders"] = Orders.([]m.OutputG5Order)
	}

	memberinfo, err := m.GetUserInfoById(member.Id); if err != nil {
		c.Response(FAIL, "")
		return
	}

	c.Data["member"] = member
	c.Data["memberinfo"] = memberinfo
	//log.Infof("memberinfo %+v", memberinfo)
	c.Data["star"], c.Data["IsBindMobile"], c.Data["IsBindEmail"] = m.GetUserStarAndBindIfo(member)
	//log.Infof("star; %v", c.Data["star"])
	//log.Infof("Orders %+v",Orders)


	//c.Ctx.WriteString("lala")
}

//@Desc: 修改密码
//
//@Param:
//
//@Return:
//
func (c *AccountController)ChangePwd() {
	memberinfo := c.GetSession("memberinfo")
	if memberinfo == nil {
		c.Response(LOGIN_EXPIRE, "")
		c.Ctx.Redirect(302, "/login")
	}
	var obj m.UserChangePwd
	c.ParseForm(&obj)
	c.ValidData(obj)
	glog.Info("obj: %+v", obj)
	if obj.NewPassword != obj.RePassword {
		c.Response(DIFFERENCE_NEWPASSWORD, "")
		return
	}
	member := memberinfo.(m.User)
	code, err := uc_client.UcUserEdit(member.UserName, obj.Password, obj.NewPassword, "", "", c.Ctx.Request.UserAgent()); if err != nil {
		c.Response(FAIL, "")
		return
	}
	switch code {
	case -1:
		c.Response(PASSWORD_ERROR, "")
		return
	case 0:
		c.ResponseMsg = "没有做任何修改"
		c.Response(FAIL, "")
		return
	case 8:
		c.ResponseMsg = "该用户受保护无权限更改"
		c.Response(FAIL, "")
		return
	case 1:
		me := m.User{
			Id:member.Id,
			Password:obj.NewPassword,
		}
		id, err := m.UpdateUser(&me); if err == nil && id > 0 {
		c.Response(SUCCESS, "")
		return
	} else {
		c.Response(FAIL, "")
		return
	}
	}

	//log.Infof("username %v", memberinfo.(m.User).UserName)
	//member, err := CheckLogin(memberinfo.(m.User).UserName, obj.Password, "username"); if err != nil {
	//	log.Infof("memberinfo:%+v", member)
	//	log.Infof("error %v", err)
	//	c.Response(PASSWORD_ERROR, "")
	//	return
	//}
	//if lib.Pwdhash(lib.Pwdhash(obj.NewPassword)) == member.Password {
	//	c.Response(SAME_NEW_OLD_PASSWORD,nil)
	//	return
	//}

}

//@Desc: 绑定手机或邮箱
//
//@Param:
//
//@Return:
//
func (c *AccountController)Bind() {
	memberinfo := c.GetSession("memberinfo")
	if memberinfo == nil {
		c.Response(LOGIN_EXPIRE, "")
		c.Ctx.Redirect(302, "/login")
	}
	member := memberinfo.(m.User)
	glog.Info("url is %v", c.Ctx.Request.RequestURI)
	glog.Info("Requ is %v\n", string(c.Ctx.Input.RequestBody))
	bindtype := c.GetString("type")
	nowtime := time.Now().Unix()
	b := m.UserBindmobile{
		UserId:uint(memberinfo.(m.User).Id),
		//UserId:uint(123),
		CreateTime : int(nowtime),
		ExpireTime : int(nowtime) + 15 * 60,
		SendTime:uint(nowtime),
		Seccode: c.GetString("vcode"),
	}
	switch bindtype {
	case "mobile":
		var obj m.UserBindmobile
		c.ParseForm(&obj)
		glog.Info("obj: %+v\n", obj)
		c.ValidData(obj)
		//if !service.VcodeVerify(obj.Mobile, obj.Seccode) {
		//	c.Response(SMS_CODE_ERROR, "")
		//	return
		//}
		if m.CheckUserExist("mobile", obj.Mobile) {
			c.Response(MOBILE_ALREADY_REGISTER, "")
			return
		}
		b.Mark = 0
		b.Mobile = obj.Mobile
		me := m.User{
			Id:memberinfo.(m.User).Id,
			Mobile:obj.Mobile,
		}
		id, err := m.UpdateUser(&me);
		if err != nil || id == 0 {
			c.Response(FAIL, "")
			return
		}
	case "email":
		var obj m.UserEmailBind
		c.ParseForm(&obj)
		glog.Info("obj: %+v\n", obj)
		c.ValidData(obj)

		code, err := uc_client.UcUserEdit(member.UserName, member.Password, "", obj.Email, "", c.Ctx.Request.UserAgent());
		if err != nil {
			c.Response(FAIL, "")
			return
		}

		if code <= 0 {
			switch code {
			case 0,-7:
				c.ResponseMsg = "没有做任何修改"
				c.Response(FAIL, "")
				return
			case -6:
				c.Response(EMAIL_ALREADY_REGISTER, "")
				return
			case -5:
				c.ResponseMsg = "Email 不允许注册"
				c.Response(FAIL,"")
				return
			default:
				c.Response(FAIL,"")
				glog.Info("uc_client.UcUserEdit return: %v",code)
				return
			}
		}
		//if !service.EmailVcodeVerify(obj.Email, obj.Vcode) {
		//	c.Response(CAPTCHA_ERROR, "")
		//	return
		//}
		if m.CheckUserExist("mail", obj.Email) {
			c.Response(EMAIL_ALREADY_REGISTER, "")
			return
		}
		b.Mark = 1
		b.Mobile = obj.Email
		me := m.User{
			Id:member.Id,
			Mail:obj.Email,
		}
		id, err := m.UpdateUser(&me);
		if err != nil || id == 0 {
			c.Response(FAIL, "")
			return
		}
	default:
		c.Response(ILLEGAL_REQUEST, "")
		return
	}
	id, err := m.AddUserBindmobile(&b); if err == nil && id > 0 {
		c.Response(SUCCESS, "")
		return
	}
	c.Response(FAIL, "");
}

//@Desc: 更新用户资料
//
//@Param:
//
//@Return:
//
func (c *AccountController)UpdateDetail() {
	glog.Info("url is %v\n", c.Ctx.Request.RequestURI)
	glog.Info("Requ is %v\n", string(c.Ctx.Input.RequestBody))
	memberinfo := c.GetSession("memberinfo")
	if memberinfo == nil {
		c.Response(LOGIN_EXPIRE, "")
		c.Ctx.Redirect(302, "/login")
	}
	var obj m.UserDetail
	c.ParseForm(&obj)
	c.ValidData(obj)
	if m.CheckUserExist("nick_name", obj.NickName) {
		c.Response(USER_ALREADY_EXIST, "")
		return
	}

	births := fmt.Sprintf("%s-%s-%s", obj.BornYear, obj.BornMonth, obj.BornDay)
	glog.Info("births: %v", births)
	birth, err := time.Parse("2006-1-2", births); if err != nil {
		glog.Info("time.Parse error %v", err)
	}
	glog.Info("birth %v", birth)
	meinfo := m.UserInfo{
		Id:memberinfo.(m.User).Id,
		NickName:obj.NickName,
		Sex:obj.Sex,
		Birth:birth,
		UserMobile:obj.Mobile,
		RealName:obj.TrueName,
		UserAddress:obj.Add,
		UserPostcode:obj.ZipCode,
	}
	glog.Info("memberinfo:%+v", meinfo)
	id, err := m.UpdateUserInfo(&meinfo); if err == nil && id > 0 {
		if len(meinfo.NickName) > 0 {
			me := m.User{
				Id:meinfo.Id,
				NickName:meinfo.NickName,
			}
			id2, err := m.UpdateUser(&me); if err == nil && id2 > 0 {
				c.Response(SUCCESS, "")
				return
			} else {
				glog.Info("error:%v", err)
				c.Response(FAIL, "")
				return
			}

		}
		c.Response(SUCCESS, "")
		return
	}
	glog.Info("obj: %+v\n", obj)
	c.Response(FAIL, "")
}
