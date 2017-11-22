package rbac

import (
	"github.com/astaxie/beego"
	. "hengzhu/admin/src"
	m "hengzhu/admin/src/models"
    "github.com/astaxie/beego/cache"
    "github.com/astaxie/beego/utils/captcha"
	"github.com/Zeniubius/golang_utils/glog"
)

var Cpt *captcha.Captcha
func init() {
    store := cache.NewMemoryCache()
    Cpt = captcha.NewWithFilter("/mcaptcha/", store)
    Cpt.ChallengeNums = 4
    Cpt.StdWidth = 100
    Cpt.StdHeight = 40
}

type MainController struct {
	CommonController
}

type Tree struct {
	Id         int64      `json:"id"`
	Text       string     `json:"text"`
	IconCls    string     `json:"iconCls"`
	Checked    string     `json:"checked"`
	State      string     `json:"state"`
	Children   []Tree     `json:"children"`
	Attributes Attributes `json:"attributes"`
}
type Attributes struct {
	Url   string `json:"url"`
	Price int64  `json:"price"`
}

//首页
func (this *MainController) Index() {
	userinfo := this.GetSession("userinfo")
	if userinfo == nil {
		this.Ctx.Redirect(302, beego.AppConfig.String("rbac_auth_gateway"))
	}
	tree:=this.GetTree()
	if this.IsAjax() {
		this.Data["json"] = &tree
		this.ServeJSON()
		return
	} else {
		groups := m.GroupList()
		this.Data["userinfo"] = userinfo
		this.Data["groups"] = groups
		this.Data["tree"] = &tree
		if this.GetTemplatetype() != "easyui"{
			this.Layout = this.GetTemplatetype() + "/manage/layout.tpl"
		}
		this.TplName = this.GetTemplatetype() + "/manage/index.tpl"
	}
}

//登录
func (this *MainController) Login() {
	needCaptcha,_ := beego.AppConfig.Bool("admin_login_captcha")
	isajax := this.GetString("isajax")
	if isajax == "1" {
		glog.Info("Captcha is %v, name:%v", this.Ctx.Request.Form.Get(Cpt.FieldIDName), this.Ctx.Request.Form.Get(Cpt.FieldCaptchaName))
        if needCaptcha && !Cpt.VerifyReq(this.Ctx.Request) {
            this.Rsp(false, "验证码错误")
            return
        }
		username := this.GetString("username")
		password := this.GetString("password")
		user, err := CheckLogin(username, password)
		if err == nil {
			this.SetSession("userinfo", user)
			accesslist, _ := GetAccessList(user.Id)
			this.SetSession("accesslist", accesslist)
			this.Rsp(true, "登录成功")
			return
		} else {
			this.Rsp(false, err.Error())
			return
		}

	}
	userinfo := this.GetSession("userinfo")
	if userinfo != nil {
		this.Ctx.Redirect(302, "/manage/index")
	}
	this.Data["IsCaptcha"] = needCaptcha
	this.TplName = this.GetTemplatetype() + "/manage/login.tpl"
}

//退出
func (this *MainController) Logout() {
	this.DelSession("userinfo")
	this.Ctx.Redirect(302, "/manage/login")
}

//修改密码
func (this *MainController) Changepwd() {
	userinfo := this.GetSession("userinfo")
	if userinfo == nil {
		this.Ctx.Redirect(302, beego.AppConfig.String("rbac_auth_gateway"))
	}
	oldpassword := this.GetString("oldpassword")
	newpassword := this.GetString("newpassword")
	repeatpassword := this.GetString("repeatpassword")
	if newpassword != repeatpassword {
		this.Rsp(false, "两次输入密码不一致")
	}
	user, err := CheckLogin(userinfo.(m.AdminUser).Username, oldpassword)
	if err == nil {
		var u m.AdminUser
		u.Id = user.Id
		u.Password = newpassword
		id, err := m.UpdateAdminUser(&u)
		if err == nil && id > 0 {
			this.Rsp(true, "密码修改成功")
			return
		} else {
			this.Rsp(false, err.Error())
			return
		}
	}
	this.Rsp(false, "密码有误")

}
