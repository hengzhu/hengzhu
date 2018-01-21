package tool

import (
	"testing"
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"hengzhu/models"
)

func TestB(t *testing.T) {

	cabinetDetail, _ := models.GetCabinetDetail(1, 5)
	fmt.Printf("detail:%v\n", cabinetDetail)
}

func init() {
	link := fmt.Sprintf("%s:%s@(%s:%s)/%s", "root", "123456", "localhost", "3306", "hengzhu")
	orm.RegisterDataBase("default", "mysql", link)

	orm.Debug = true
}
