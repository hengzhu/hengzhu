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

func TestC(t *testing.T) {
	a := 58
	fmt.Println("value of a before function call is",a)
	b := &a
	change(b)
	fmt.Println("value of a after function call is", a)
}

func init() {
	link := fmt.Sprintf("%s:%s@(%s:%s)/%s", "root", "123456", "localhost", "3306", "hengzhu")
	orm.RegisterDataBase("default", "mysql", link)

	orm.Debug = true
}

func change(val *int) {
	*val = 55
}