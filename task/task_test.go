package task

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"testing"
	_ "github.com/go-sql-driver/mysql"
)

// 测试所有游戏停运相关的
func TestFlushLog(t *testing.T) {
	FlushLog()
}


func init() {
	//link := fmt.Sprintf("%s:%s@(%s:%s)/%s", "kuaifa_on",
	//	"kuaifazs", "10.8.230.17",
	//	"3308", "work_together_online")
	link := fmt.Sprintf("%s:%s@(%s:%s)/%s", "root",
		"123456", "localhost",
		"3306", "hengzhu")
	fmt.Printf("link:%v", link)
	orm.RegisterDataBase("default", "mysql", link)

	orm.Debug = true
}
