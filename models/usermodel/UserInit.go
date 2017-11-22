package usermodel

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"fmt"
)

var o orm.Ormer
var Dbname = beego.AppConfig.String("db_name")

func Connect()  {
	var dsn string
	db_type := beego.AppConfig.String("db_type")
	db_host := beego.AppConfig.String("db_host")
	db_port := beego.AppConfig.String("db_port")
	db_user := beego.AppConfig.String("db_user")
	db_pass := beego.AppConfig.String("db_pass")
	db_name := beego.AppConfig.String("minisdk_db_name")
	orm.RegisterDriver("mysql", orm.DRMySQL)
	dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", db_user, db_pass, db_host, db_port, db_name)
	orm.RegisterDataBase("minisdktest", db_type, dsn)
}