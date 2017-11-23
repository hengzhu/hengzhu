package rbac

import (
	"github.com/astaxie/beego"
	. "hengzhu/admin/src"
	m "hengzhu/admin/src/models"
	"hengzhu/models/bean"
)

type CommonController struct {
	beego.Controller
	Templatetype string //ui template type
}

func (this *CommonController) Rsp(status bool, str string) {
	this.Data["json"] = &map[string]interface{}{"status": status, "info": str}
	this.ServeJSON()
}

func (this *CommonController) GetTemplatetype() string {
	templatetype := beego.AppConfig.String("template_type")
	if templatetype == "" {
		templatetype = "easyui"
	}
	return templatetype
}

func (this *CommonController) GetTree() []Tree {
	nodes, _ := m.GetNodeTree(0, 1)
	tree := make([]Tree, len(nodes))
	for k, v := range nodes {
		tree[k].Id = v["Id"].(int64)
		tree[k].Text = v["Title"].(string)
		children, _ := m.GetNodeTree(v["Id"].(int64), 2)
		tree[k].Children = make([]Tree, len(children))
		for k1, v1 := range children {
			tree[k].Children[k1].Id = v1["Id"].(int64)
			tree[k].Children[k1].Text = v1["Title"].(string)
			tree[k].Children[k1].Attributes.Url = "/" + v["Name"].(string) + "/" + v1["Name"].(string)
		}
	}
	return tree
}

func (c *CommonController) RespJSON(code int, data interface{}) {
	c.AllowCross()
	c.Ctx.Output.SetStatus(code)
	var hasIndent = true
	if beego.BConfig.RunMode == beego.PROD {
		hasIndent = false
	}
	c.Ctx.Output.JSON(data, hasIndent, false)
}

func (c *CommonController) RespDataMsg(code int, msg string) {
	c.RespJSON(code, map[string]interface{}{
		"code":  code,
		"msg": msg,
	})
}

// 只有数据, 返回值默认为200(成功)
func (c *CommonController) RespJSONData(data interface{}) {
	c.AllowCross()
	c.RespJSON(bean.CODE_Success, data)
}

// 只有数据, 返回值默认为200(成功)
func (c *CommonController) RespJSONDataWithTotal(data interface{}, total int64) {
	c.RespJSON(bean.CODE_Success, map[string]interface{}{
		"rows":  data,
		"total": total,
	})
}

func (c *CommonController) RespJSONDataWithSumAndTotal(data interface{}, total int64, sum float64) {
	c.RespJSON(bean.CODE_Success, map[string]interface{}{
		"data":  data,
		"total": total,
		"sum": sum,
	})
}

func (c *CommonController) AllowCross() {
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")                               //允许访问源
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Methods", "POST,DELETE, GET, PUT, OPTIONS") //允许post访问
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")     //header的类型
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Max-Age", "1728000")
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Ctx.ResponseWriter.Header().Set("content-type", "application/json") //返回数据格式是json
}

func init() {

	//验证权限
	AccessRegister()
}
