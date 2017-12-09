package task

import (
	"github.com/astaxie/beego/toolbox"
	"fmt"
	"github.com/astaxie/beego/orm"
	"time"
	"hengzhu/models"
)

func Run() {
	// 每天凌晨1点清除日志
	FlushLog := toolbox.NewTask("FlushLog", "0 0 1 * * * ", func() error {
		FlushLog()
		return nil
	})

	toolbox.AddTask("FlushLog", FlushLog)

	toolbox.StartTask()
	fmt.Println("task success")
}

// 根据配置文件的日志保存时间，删除之前的日志
func FlushLog() {
	o := orm.NewOrm()

	setting, _ := models.GetSettingById(1)

	o.QueryTable(new(models.Log)).Filter("time__lte", time.Now().AddDate(0, 0, -setting.LogTime)).Delete()
}
