package uc_client

import (
	"hengzhu/ucenter/uc_client/data/cache"
	"net/url"
	"regexp"
	"hengzhu/utils"
	"hengzhu/ucenter"
	"fmt"
	"github.com/deepzz0/go-com/log"
	ucc "hengzhu/ucenter/uc_client/controllers"
	"hengzhu/models/usermodel"
	"time"
	"strings"
	"net/http"
	"io/ioutil"
	"hengzhu/admin/src/lib"
	"errors"
	"strconv"
	"encoding/xml"
)

const (
	IN_UC = true
	UC_CLIENT_VERSION = "1.6.0"
	UC_CLIENT_RELEASE = "20110501"
	//UC_ROOT
	//UC_DATADIR
	//UC_DATAURL
	//UC_API_FUNC

)

var uc ucc.UcUserController

func init() {
	uc.Usercontrol()
}

func UcUserSynlogin(u usermodel.User) string {
	if len(cache.Apps) > 0 {
		return uc.OnSynLogin(u)
	}
	return ""
}

func UcUserSynlogout() string {
	if len(cache.Apps) > 0 {
		//uc.Usercontrol()
		return uc.OnSynLogout()
	}
	return ""
}

func UcApiRequestData(module string, action string, arg string, extra string) (post string) {
	input := UcApiInput(arg, "")
	post = fmt.Sprintf("m=%s&a=%s&inajax=2&release=%s&input=%s&appid=%s%s",
		module, action, UC_CLIENT_RELEASE, input, ucenter.UC_APPID, extra, )
	return
}

func UcApiInput(data string, useragent string) string {
	ua := &http.Request{}
	log.Printf("useragent:%+v", ua)
	log.Printf("useragent:%+v", useragent)

	//useragent1 := "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2840.99 Safari/537.36"
	params := fmt.Sprintf("%s&agent=%s&time=%d", data, lib.Strtomd5(useragent), time.Now().Unix())        //time.Now().Unix()
	log.Printf("params:%v", params)
	log.Printf("UC_KEY:%v", ucenter.UC_KEY)
	step := utils.UcAuthcode(params, "ENCODE", ucenter.UC_KEY, 0)
	return step
}

func DecodeUrl(u url.Values) (string, error) {
	s := u.Encode()
	reg, err := regexp.Compile(`%5B(.*?)%5D`); if err != nil {
		return "", err
	}
	return reg.ReplaceAllString(s, "[$1]"), nil
}


//@Desc: 注册到Ucenter
//
//@Param: u usermodel.UserReg 用户; useragent string 用户代理
//
//@Return: res int Ucenter返回的结果; err error
//
func UcUserRegister(u usermodel.UserReg, useragent string) (res int, err error) {
	args := map[string]string{
		"username":u.UserName,
		"password":u.Password,
		"email":fmt.Sprintf("%s@youcai.com", u.UserName),
	}
	step, err := UcApiPost("user", "register", args, useragent); if err != nil {
		return -7, err
	}
	res, err = strconv.Atoi(step)
	return
}

//@Desc: Ucenter用户登录
//
//@Param:
//
//@Return:
//
func UcUserLogin(u usermodel.UserLogin, isuid string, useragent string) (res UserLoginRes, err error) {
	args := map[string]string{
		"username":u.UserName,
		"password":u.Password,
		"email":fmt.Sprintf("%s@youcai.com", u.UserName),
		"isuid":isuid,
	}
	stepxml, err := UcApiPost("user", "login", args, useragent)
	step := Root{}
	xml.Unmarshal([]byte(stepxml), &step)
	for key, _ := range step.Items {
		switch step.Items[key].Id {
		case 0:
			result, _ := strconv.Atoi(step.Items[key].Data)
			res.Result = result
		case 1:
			res.UserName = step.Items[key].Data
		case 2:
			res.Password = step.Items[key].Data
		case 3:
			res.Email = step.Items[key].Data
		case 4:
			res.IsRepeatName = step.Items[key].Data
		}
	}
	return
}

//@Desc: 获取用户数据
//
//@Param:
//
//@Return:
//
func UcGetUser(u usermodel.User, gettype string, useragent string) (res string, err error) {
	args := make(map[string]string)
	if gettype == "username" {
		args["username"] = u.UserName
	} else if gettype == "uid" {
		args["username"] = strconv.Itoa(u.Id)
	} else {
		return "", errors.New("")
	}
	return UcApiPost("user", "get_user", args, useragent)

}

//@Desc: 更新用户资料
//
//@Param:
//
//@Return:
//
func UcUserEdit(username string, oldpw string, newpw string, email string, ignoreoldpw string, useragent string) (res int, err error) {
	args := map[string]string{
		"username":username,
		"oldpw":oldpw,
		"newpw":newpw,
		"email":email,
		"ignoreoldpw":ignoreoldpw,
	}
	step, err := UcApiPost("user", "edit", args, useragent); if err != nil {
		return -7, err
	}
	res, err = strconv.Atoi(step)
	return

}



//@Desc: 发送请求到论坛的Ucenter Server
//
//@Param: module string 模块; action string 方法; args map[string]string 参数; useragent string 用户代理
//
//@Return: res string Ucenter 返回的结果; err error 执行过程的错误
//
func UcApiPost(module string, action string, args map[string]string, useragent string) (res string, err error) {
	uv := url.Values{}
	for key, value := range args {
		uv.Set(key, value)
	}
	userdata := uv.Encode()
	log.Infof("userdata %v",userdata)
	input := UcApiInput(userdata, useragent)

	v := url.Values{}
	v.Set("m", module)
	v.Set("a", action)
	v.Set("inajax", "2")
	v.Set("input", input)
	v.Set("release", UC_CLIENT_RELEASE)
	appid := fmt.Sprintf("%v%v", ucenter.UC_APPID, "")
	v.Set("appid", appid)
	senddata := v.Encode()
	tourl := ucenter.UC_API + "/index.php"
	log.Printf("UcApiPost to url: %v", tourl)
	log.Printf("UcApiPost senddata: %v", senddata)

	body := strings.NewReader(senddata)
	client := &http.Client{}
	log.Printf("UcApiPost: body: %v", body)
	req, err := http.NewRequest("POST", tourl, body); if err != nil {
		log.Printf("UcApiPost: NewRequest error: %v", err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;param=value")
	req.Header.Set("User-Agent", useragent)
	resd, err2 := client.Do(req); if err2 != nil {
		log.Printf("UcApiPost: Client.Do error: %v", err2)
		return "", err2
	}
	defer resd.Body.Close()
	data, _ := ioutil.ReadAll(resd.Body)
	res = string(data)
	log.Printf("UcApiPost: res data:%v", res)
	return res, nil
}
