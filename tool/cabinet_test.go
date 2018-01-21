package tool

import (
	"testing"
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func TestB(t *testing.T) {

	fmt.Printf("door:%v---\n", CabinetDoors["asdfasdfadsfsf"])
}

func init() {
	link := fmt.Sprintf("%s:%s@(%s:%s)/%s", "root", "123456", "localhost", "3306", "hengzhu")
	orm.RegisterDataBase("default", "mysql", link)

	orm.Debug = true
}
