package main

import (
	_ "hengzhu/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"net/url"
	"hengzhu/models/admin"
	_ "github.com/go-sql-driver/mysql"
	"hengzhu/task"
	"hengzhu/tool"
)

func main() {
	//if beego.BConfig.RunMode == "dev" {
	//	beego.BConfig.WebConfig.DirectoryIndex = true
	//	beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	//}
	go tool.GetMessageFromHardWare(tool.Queues)
	go tool.GetNewCabinet(tool.NewCabinet)
	task.Run()
	beego.Run()
}

func init() {
	//link := fmt.Sprintf("%s:%s@(%s:%s)/%s", beego.AppConfig.String("mysqluser"),
	//	beego.AppConfig.String("mysqlpass"), beego.AppConfig.String("mysqlurls"),
	//	beego.AppConfig.String("mysqlport"), beego.AppConfig.String("mysqldb"))
	//orm.RegisterDataBase("default", "mysql", link)

	dbhost := beego.AppConfig.String("db_host")
	dbport := beego.AppConfig.String("db_port")
	dbuser := beego.AppConfig.String("db_user")
	dbpassword := beego.AppConfig.String("db_pass")
	dbname := beego.AppConfig.String("db_name")
	timezone := beego.AppConfig.String("db_timezone")
	if dbport == "" {
		dbport = "3306"
	}
	dsn := dbuser + ":" + dbpassword + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?charset=utf8"
	// fmt.Println(dsn)

	if timezone != "" {
		dsn = dsn + "&loc=" + url.QueryEscape(timezone)
	}
	orm.RegisterDataBase("default", "mysql", dsn)
	orm.RegisterModel(new(admin.Auth), new(admin.Role), new(admin.RoleAuth), new(admin.Admin))

	orm.Debug = beego.BConfig.RunMode == "dev"
	orm.RunSyncdb("default", false, true)

	//logs.SetLogger(logs.AdapterFile, `{"filename":"info.log","level":6,"maxlines":0,"maxsize":0,"daily":true,"maxdays":1000}`)
}
