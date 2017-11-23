package main

import (
	_ "hengzhu/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

func main() {
	//if beego.BConfig.RunMode == "dev" {
	//	beego.BConfig.WebConfig.DirectoryIndex = true
	//	beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	//}

	beego.Run()
}

func init() {
	//link := fmt.Sprintf("%s:%s@(%s:%s)/%s", beego.AppConfig.String("mysqluser"),
	//	beego.AppConfig.String("mysqlpass"), beego.AppConfig.String("mysqlurls"),
	//	beego.AppConfig.String("mysqlport"), beego.AppConfig.String("mysqldb"))
	//orm.RegisterDataBase("default", "mysql", link)

	orm.Debug = beego.BConfig.RunMode == "dev"
	orm.RunSyncdb("default", false, true)
	//logs.SetLogger(logs.AdapterFile, `{"filename":"info.log","level":6,"maxlines":0,"maxsize":0,"daily":true,"maxdays":1000}`)
}
