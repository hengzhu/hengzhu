package tool

import (
	"testing"
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"hengzhu/models"
	"time"
	"reflect"
)

func TestB(t *testing.T) {

	cabinetDetail, _ := models.GetCabinetDetail(1, 5)
	fmt.Printf("detail:%v\n", cabinetDetail)

	fmt.Printf("time64:%v,type:%v\n", time.Now().Unix(), reflect.TypeOf(time.Now().Unix()))
	fmt.Printf("time:%v,type:%v\n", int(time.Now().Unix()), reflect.TypeOf(int(time.Now().Unix())))
}

func init() {
	link := fmt.Sprintf("%s:%s@(%s:%s)/%s", "root", "123456", "localhost", "3306", "hengzhu")
	orm.RegisterDataBase("default", "mysql", link)

	orm.Debug = true
}
