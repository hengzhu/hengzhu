package controllers

import (
	"hengzhu/models/usermodel"
	"strconv"
	"log"
	"github.com/astaxie/beego/orm"
)

const (
	INVALID_REQUEST = "Invalid Request"
	AUTHENTICATION_EXPIRED = "Authentication expired"
	API_RETURN_SUCCEED = "1"
	API_SYNLOGIN = 1
	API_RETURN_FORBIDDEN = 1
)

var UcNote = make(map[string]func(c *ApiController, req map[string]string) (string))

func init() {
	UcNote["test"] = func(c *ApiController, req map[string]string) (string) {
		return API_RETURN_SUCCEED
	}

	UcNote["synlogin"] = func(c *ApiController, req map[string]string) string {
		if API_SYNLOGIN != 1 {
			return INVALID_REQUEST
		}
		//c.Ctx.ResponseWriter.Header().Set("P3P", `CP="CURa ADMa DEVa PSAo PSDo OUR BUS UNI PUR INT DEM STA PRE COM NAV OTC NOI DSP COR"`)
		id, err := strconv.Atoi(req["uid"])
		member, err := usermodel.GetUserById(id); if err != nil {
			if err == orm.ErrNoRows {
				memberp, err := usermodel.SyncUser(id, req["username"], "", ""); if err != nil {
					return err.Error()
				}
				member = &memberp
			} else {
				return err.Error()
			}

		}
		c.SetSession("memberinfo", *member)
		log.Printf("Struct: %v", c.CruSession)
		return member.UserName

	}

	UcNote["synlogout"] = func(c *ApiController, req map[string]string) string {
		if API_SYNLOGIN != 1 {
			return INVALID_REQUEST
		}
		c.Ctx.ResponseWriter.Header().Set("P3P", `CP="CURa ADMa DEVa PSAo PSDo OUR BUS UNI PUR INT DEM STA PRE COM NAV OTC NOI DSP COR"`)
		c.DelSession("memberinfo")
		return ""
	}

	UcNote["updatepw"] = func(c *ApiController, req map[string]string) string {
		return ""
	}

}
