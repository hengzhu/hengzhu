package models

import (
	"testing"
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func TestA(t *testing.T) {
	flag := CheckIfAdd("012345")
	fmt.Printf("flag:%v\n", flag)
}

func TestB(t *testing.T) {
	typ := GetDefaultType()
	fmt.Printf("typ:%v\n", typ)
}

func init() {
	link := fmt.Sprintf("%s:%s@(%s:%s)/%s", "root", "123456", "localhost", "3306", "hengzhu")
	orm.RegisterDataBase("default", "mysql", link)

	orm.Debug = true
}
