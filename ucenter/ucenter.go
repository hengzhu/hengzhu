package ucenter

import (
	"mime"
	"hengzhu/models/usermodel"
	"github.com/astaxie/beego"
)

var (
	//UC_DBHOST = beego.AppConfig.String("uc_dbhost")
	//UC_CONNECT = beego.AppConfig.String("uc_connect")
	//UC_DBUSER = beego.AppConfig.String("uc_dbuser")
	//UC_DBPW = beego.AppConfig.String("uc_dbpw")
	//UC_DBNAME = beego.AppConfig.String("uc_dbname")
	//UC_DBCONNECT = beego.AppConfig.String("uc_db")
	//UC_DBCHARSET = beego.AppConfig.String("uc_dbcharset")
	//UC_DBTABLEPRE = beego.AppConfig.String("uc_dbtablepre")
	UC_KEY = beego.AppConfig.String("uc_key")
	//UC_API = beego.AppConfig.String("uc_api")
	//UC_CHARSET = beego.AppConfig.String("uc_charset")
	//UC_IP = beego.AppConfig.String("uc_ip")
	UC_APPID, _ = beego.AppConfig.Int("uc_appid")
	UC_PPP = beego.AppConfig.String("uc_ppp")
)

const (
	UC_API = "http://test.youcaibbs.com/uc_server"
	//UC_APPID = "2"
	//UC_KEY = "123456"
)
/**
*@Desc
*@Param
*@Return int64
*/
func Run() {
	initialize()
}

func initialize() {
	mime.AddExtensionType(".css", "text/css")
	usermodel.Connect()
	router()
	//beego.AddFuncMap("stringsToJson", StringsToJson)
}

func initArgs() {

}





