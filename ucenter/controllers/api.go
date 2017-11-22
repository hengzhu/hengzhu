package controllers

import (
	"hengzhu/models/usermodel"
	"errors"
	"strings"

	"github.com/astaxie/beego"
	"log"
	"hengzhu/utils"
	"time"
	"strconv"
)

// ApiController Ucenter 接口
type ApiController struct {
	beego.Controller
}


/**
*@Desc Ucenter 接口主函数，判断连接状态，分配相应方法
*@Param
*@Return
*/
func (c *ApiController) Index() {
	code := c.GetString("code")
	log.Printf("url is %v", c.Ctx.Request.RequestURI)
	log.Printf("code is: %v", code)
	//code, _ = url.QueryUnescape(code)
	log.Printf("UC_KEY is: %v", beego.AppConfig.String("uc_key"))
	res := utils.UcAuthcode(code, "DECODE", beego.AppConfig.String("uc_key"), 0)
	if res == "" {
		c.Ctx.ResponseWriter.Write([]byte(INVALID_REQUEST))
		return
	}
	get := utils.QueryString2Map(res)
	tm, _ := strconv.Atoi(get["time"])
	if (int(time.Now().Unix()) - tm) > 3600 {
		c.Ctx.ResponseWriter.Write([]byte(AUTHENTICATION_EXPIRED))
		return
	}
	log.Printf("api res is %v", res)
	log.Printf("map res is %v", get)
	callfunc, ok := UcNote[get["action"]]
	if !ok {
		return
	}
	ress := callfunc(c, get)
	c.Ctx.ResponseWriter.Write([]byte(ress))
	log.Printf("callfunction's return: %v", ress)
}

// URLMapping ...
func (c *ApiController) URLMapping() {
	c.Mapping("GetAll", c.GetAll)
}


// GetAll ...
// @Title Get All
// @Description get User
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.User
// @Failure 403
// @router / [get]
func (c *ApiController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Data["json"] = errors.New("Error: invalid query key/value pair")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := usermodel.GetAllUser(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}

