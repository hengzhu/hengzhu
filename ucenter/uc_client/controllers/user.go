package controllers

import (
	"hengzhu/models/usermodel"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"time"
	"net/url"
	"fmt"
	"hengzhu/utils"
	"log"
	"hengzhu/ucenter/uc_client/data/cache"
	"hengzhu/ucenter"
)

// UserController oprations for User
type UcUserController struct {
	UcBaseController
}

func init() {

}



// URLMapping ...
func (c *UcUserController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

func (c *UcUserController)Usercontrol() {
	c.cache = cache.Apps
	c.app = c.cache[ucenter.UC_APPID]
}

func (c *UcUserController)OnSynLogin(u usermodel.User) string {
	if c.app.Synlogin != 0 {
		var synstr string
		for _, app := range c.cache {
			if app.Synlogin != 0 && app.Appid != c.app.Appid {
				code := url.QueryEscape(utils.UcAuthcode(fmt.Sprintf("action=synlogin&username=%s&uid=%d&password=%s&time=%d", u.UserName, u.Id, u.Password, time.Now().Unix()), "ENCODE", app.AuthKey, 0))
				synstr = fmt.Sprintf(`<script type="text/javascript" src="%s/api/uc.php?time=%d&code=%s"></script>`, app.Url, time.Now().Unix(), code)
			}
		}
		log.Printf("SynLogin synstr: %v", synstr)
		return synstr
	}
	return ""
}

func (c *UcUserController)OnSynLogout() string {
	c.Usercontrol()
	if c.app.Synlogin != 0 {
		var synstr string
		for _, app := range c.cache {
			if app.Synlogin != 0 && app.Appid != c.app.Appid {
				code := url.QueryEscape(utils.UcAuthcode(fmt.Sprintf("action=synlogout&time=%d", time.Now().Unix()), "ENCODE", app.AuthKey, 0))
				synstr = fmt.Sprintf(`<script type="text/javascript" src="%s/api/uc.php?time=%d&code=%s"></script>`, app.Url, time.Now().Unix(), code)
			}
		}
		log.Printf("Logout Synstr is: %v", synstr)
		return synstr
	}
	return ""
}

func (c *UcUserController)OnRegister() {

}



// Post ...
// @Title Post
// @Description create User
// @Param	body		body 	models.User	true		"body for User content"
// @Success 201 {int} models.User
// @Failure 403 body is empty
// @router / [post]
func (c *UcUserController) Post() {
	var v usermodel.User
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if _, err := usermodel.AddUser(&v); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = v
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get User by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.User
// @Failure 403 :id is empty
// @router /login/:id [get]
func (c *UcUserController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := usermodel.GetUserById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
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
func (c *UcUserController) GetAll() {
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

// Put ...
// @Title Put
// @Description update the User
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.User	true		"body for User content"
// @Success 200 {object} models.User
// @Failure 403 :id is not int
// @router /:id [put]
func (c *UcUserController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := usermodel.User{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := usermodel.UpdateUserById(&v); err == nil {
			c.Data["json"] = "OK"
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the User
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *UcUserController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := usermodel.DeleteUser(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
