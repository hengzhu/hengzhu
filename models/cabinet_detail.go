package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"github.com/astaxie/beego/orm"
	"time"
	"strconv"
	"github.com/astaxie/beego"
	"hengzhu/utils"
)

type CabinetDetail struct {
	Id            int    `orm:"column(id);auto"`
	CabinetId     int    `orm:"column(cabinet_id);null" description:"柜子的id"`
	Door          int    `orm:"column(door);null" description:"门号"`
	OpenState     int    `orm:"column(open_state);null" description:"开关状态，1:关，2:开"`
	Using         int    `orm:"column(using);null" description:"占用状态，1:空闲，2:占用"`
	UserID        string `orm:"column(userID);size(255);null" description:"存物ID"`
	StoreTime     int    `orm:"column(store_time);null" description:"存物时间"`
	UseState      int    `orm:"column(use_state);null" description:"启用状态，1:启用，2:停用"`
	WireConnected int    `orm:"column(wire_connected);null" description:"该门电线连接状态,1:正常连接，2:不正常"`

	ID                string `orm:"-"`
	Logs              []Log  `orm:"-"`
	StoreTimeFormated string `orm:"-"`
}

type Total struct {
	Doors int `json:"-" orm:"column(doors)"`
	OnUse int `json:"-" orm:"column(onUse)"`
	Close int `json:"-" orm:"column(close)"`
}

func (t *CabinetDetail) TableName() string {
	return "cabinet_detail"
}

func init() {
	orm.RegisterModel(new(CabinetDetail))
}

// 根据柜子id，获取该柜子的门的详情
func GetDetailsByCabinetId(cabinetId int) (details []CabinetDetail, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable(new(CabinetDetail)).Filter("CabinetId", cabinetId).All(&details)
	return
}

func GetDetail(cabinetId int, door int) (detail CabinetDetail, err error) {
	o := orm.NewOrm()
	err = o.QueryTable(new(CabinetDetail)).Filter("CabinetId", cabinetId).Filter("door", door).One(&detail)
	return
}

//
func AddAllInfo(details []CabinetDetail) {
	if len(details) == 0 {
		return
	}
	for i, detail := range details {
		if detail.StoreTime != 0 {
			details[i].StoreTimeFormated = time.Unix(int64(detail.StoreTime), 0).Format("2006-01-02 15:04:05")
		} else {
			details[i].StoreTimeFormated = "--"
		}
	}
}

// 根据柜子id，获取该柜子的总门数
func GetTotalDoors(cabinetId int) (total int) {
	tot := Total{}
	sql := "SELECT COUNT(id) AS doors FROM cabinet_detail WHERE cabinet_id=?"
	orm.NewOrm().Raw(sql, cabinetId).QueryRow(&tot)
	return tot.Doors
}

// 根据柜子id，获取该柜子的总使用中的门数
func GetTotalOnUse(cabinetId int) (onUse int) {
	tot := Total{}
	sql := "SELECT COUNT(`using`) AS onUse FROM cabinet_detail WHERE cabinet_id=? AND `using`=2"
	orm.NewOrm().Raw(sql, cabinetId).QueryRow(&tot)
	return tot.OnUse
}

// 根据柜子id，获取该柜子的总关闭状态门数
func GetTotalClose(cabinetId int) (close int) {
	tot := Total{}
	sql := "SELECT COUNT(open_state) AS `close` FROM cabinet_detail WHERE cabinet_id=? AND open_state=1"
	orm.NewOrm().Raw(sql, cabinetId).QueryRow(&tot)
	return tot.Close
}

// 给柜子门附加柜子ID信息
func AddIDInfo(detail *CabinetDetail) {
	if detail == nil {
		return
	}
	cabinet := Cabinet{}
	sql := "SELECT cabinet_ID FROM cabinet WHERE id=?"
	orm.NewOrm().Raw(sql, detail.CabinetId).QueryRow(&cabinet)
	detail.ID = cabinet.CabinetID
	return
}

// 根据起始日期和终止日期，给柜子门附加历史记录
func AddLogInfo(detail *CabinetDetail, beginTime string, endTime string) {
	if detail == nil {
		return
	}
	o := orm.NewOrm()
	logs := make([]Log, 0)

	o.QueryTable(new(Log)).Filter("CabinetDetailId", detail.Id).
		Filter("time__gte", beginTime).
		Filter("time__lte", endTime).OrderBy("-time").All(&logs)
	detail.Logs = logs
}

// AddCabinetDetail insert a new CabinetDetail into database and returns
// last inserted Id on success.
func AddCabinetDetail(m *CabinetDetail) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetCabinetDetailById retrieves CabinetDetail by Id. Returns error if
// Id doesn't exist
func GetCabinetDetailById(id int) (v *CabinetDetail, err error) {
	o := orm.NewOrm()
	v = &CabinetDetail{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllCabinetDetail retrieves all CabinetDetail matches certain condition. Returns empty list if
// no records exist
func GetAllCabinetDetail(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(CabinetDetail))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []CabinetDetail
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateCabinetDetail updates CabinetDetail by Id and returns error if
// the record to be updated doesn't exist
func UpdateCabinetDetailById(m *CabinetDetail) (err error) {
	o := orm.NewOrm()
	v := CabinetDetail{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// UpdateCabinetDetail with NULL
func UpdateCabinetDetailWithNUll(m *CabinetDetail) (err error) {
	o := orm.NewOrm()
	v := CabinetDetail{Id: m.Id}
	if err = o.Read(&v); err == nil {
		_, err = o.Raw("update cabinet_detail set `using` = ?, userID = null, store_time = ? where id = ? ;", m.Using, m.StoreTime, v.Id).Exec()
	}
	return
}

// DeleteCabinetDetail deletes CabinetDetail by Id and returns error if
// the record to be deleted doesn't exist
func DeleteCabinetDetail(id int) (err error) {
	o := orm.NewOrm()
	v := CabinetDetail{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&CabinetDetail{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

//获取空闲可用的柜子门
func GetFreeDoorByCabinetId(cabinet_id int) (v *CabinetDetail, err error) {
	o := orm.NewOrm()
	cd := CabinetDetail{}
	err = o.QueryTable("cabinet_detail").Filter("open_state", 1).Filter("use_state", 1).
		Filter("using", 1).Filter("cabinet_id", cabinet_id).Filter("wire_connected", 1).Limit(1).One(&cd)
	if err != nil {
		return
	}
	v = &cd
	return
}

////更新柜子状态
//func UpdateCabinetDoorStatusById(m *CabinetDetail, openid string) (err error) {
//	o := orm.NewOrm()
//	v := CabinetDetail{Id: m.Id}
//	// ascertain id exists in the database
//	if err = o.Read(&v); err == nil {
//		m.UserID = openid
//		if _, err = o.Update(m, "using", "userID"); err == nil {
//			fmt.Println("用户" + openid + "使用的门id为:" + strconv.Itoa(m.Id))
//		}
//	}
//	return
//}

func GetCabinetDetail(cid int, door_no int) (v *CabinetDetail, err error) {
	o := orm.NewOrm()
	v = &CabinetDetail{CabinetId: cid, Door: door_no}
	if err = o.Read(v, "cabinet_id", "door"); err == nil {
		return v, nil
	}
	return nil, err
}

func GetCabinetByOutOrderNo(out_order_no string) (v *CabinetDetail, err error) {
	o := orm.NewOrm()
	c := CabinetDetail{}
	err = o.Raw("select cabinet_id from cabinet_detail where id = (select cabinet_detail_id from cabinet_order_record where order_no = ?) limit 1;", out_order_no).QueryRow(&c)
	if err != nil {
		return
	}
	v = &c
	return
}

//绑定openid到柜子门号
func BindOpenIdForCabinetDoor(openid string, cdid int) (err error, door_no int) {
	o := orm.NewOrm()
	v := CabinetDetail{Id: cdid}
	if err = o.Read(&v); err == nil {
		v.UserID = openid
		v.Using = 2
		if _, err = o.Update(&v, "userID", "using"); err == nil {
			fmt.Println("用户" + openid + "使用的门id为：" + strconv.Itoa(cdid))
			door_no = v.Door
		}
	}
	return
}

func GetCabinetDetailByOpenId(open_id string) (v *CabinetDetail, err error) {
	o := orm.NewOrm()
	v = &CabinetDetail{UserID: open_id}
	if err := o.Read(v, "userID"); err == nil {
		return v, nil
	}
	return nil, err
}

//func UpdateCabinetDetail(m *CabinetDetail) (err error) {
//	o := orm.NewOrm()
//	v := CabinetDetail{Id: m.Id}
//	if err = o.Read(&v); err == nil {
//		var num int64
//		if num, err = o.Update(m); err == nil {
//			fmt.Println("Number of records updated in database:", num)
//		}
//	}
//	result1, _ := utils.Redis.GET(utils.PAY + strconv.Itoa(m.Id))
//	result2, _ := utils.Redis.GET(utils.NOPAY + strconv.Itoa(m.Id))
//
//	cd := CabinetDetail{}
//	//先查是否被占用
//	err = o.Raw("select userID,`using`,store_time from cabinet_detail where id = ? limit 1;", m.Id).QueryRow(&cd)
//	if err != nil {
//		return
//	}
//	if cd.UserID != "" && cd.Using == 2 && m.OpenState == 1 {
//		cor := CabinetOrderRecord{}
//		//如果同一用户又用了同一个门?
//		//是否已经支付过
//		err = o.Raw("select * from cabinet_order_record where customer_id = ? and cabinet_detail_id = ? and is_pay = 1 and (past_flag is null or past_flag = 0) limit 1;", cd.UserID, m.Id).QueryRow(&cor)
//		//当前使用但未支付
//		if err == nil && cd.StoreTime == 0 {
//			//第一次关门
//			_, err = o.Raw("update cabinet_detail set `using` = 2, store_time = ?, userID = ? where id = ? limit 1;", int(time.Now().Unix()), m.UserID, cid).Exec()
//			//添加日志记录
//			m := Log{
//				CabinetDetailId: m.Id,
//				User:            cd.UserID,
//				Time:            time.Now(),
//				Action:          "存",
//			}
//			AddLog(&m)
//			//删除缓存
//			err = utils.Redis.DEL(key)
//			if err != nil {
//				beego.Error(err)
//			}
//			return
//		} else if err == nil && cd.StoreTime != 0 {
//			//柜子置为空闲
//			_, err = o.Raw("update cabinet_order_record set past_flag = 1 where id = ? ;", cor.Id).Exec()
//			if err != nil {
//				beego.Error(err)
//				return
//			}
//			_, err = o.Raw("update cabinet_detail set userID = null,`using` = ?,store_time = ? where id = ? ;", 1, 0, m.Id).Exec()
//			//添加日志记录
//			m := Log{
//				CabinetDetailId: m.Id,
//				User:            cd.UserID,
//				Time:            time.Now(),
//				Action:          "取",
//			}
//			AddLog(&m)
//		}
//		if err == orm.ErrNoRows {
//			//先存后付第一次关门
//			_, err = o.Raw("update cabinet_detail set `using` = 2, store_time = ? where userID = ? limit 1;", int(time.Now().Unix()), m.UserID).Exec()
//			//添加日志记录
//			m := Log{
//				CabinetDetailId: m.Id,
//				User:            cd.UserID,
//				Time:            time.Now(),
//				Action:          "存",
//			}
//			AddLog(&m)
//			return
//		}
//		if err != nil && err != orm.ErrNoRows {
//			err = errors.New("系统异常")
//			return
//		}
//	}
//	return
//}

func UpdateCabinetDetail(m *CabinetDetail) (err error) {
	o := orm.NewOrm()
	v := CabinetDetail{Id: m.Id}
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	//管理员操作
	managerResult, _ := utils.Redis.GET(utils.MANAGER + strconv.Itoa(m.Id))
	if managerResult != "" {
		if m.OpenState == 1 {
			err2 := utils.Redis.DEL(utils.MANAGER + strconv.Itoa(m.Id))
			if err2 != nil {
				beego.Error(err2)
				err = err2
			}
		}
		return
	}
	if m.OpenState == 1 {
		var key string
		var result string
		var result1, _ = utils.Redis.GET(utils.PAY + strconv.Itoa(m.Id))
		if result1 != "" {
			result = result1
			key = utils.PAY + strconv.Itoa(m.Id)
		}
		var result2, _ = utils.Redis.GET(utils.NOPAY + strconv.Itoa(m.Id))
		if result2 != "" {
			result = result2
			key = utils.NOPAY + strconv.Itoa(m.Id)
		}
		if result != "" {
			userId := result
			//第一次关门
			_, err = o.Raw("update cabinet_detail set `using` = 2, store_time = ?, userID = ? where id = ? limit 1;", int(time.Now().Unix()), userId, m.Id).Exec()
			//添加日志记录
			m := Log{
				CabinetDetailId: m.Id,
				User:            userId,
				Time:            time.Now(),
				Action:          "存",
			}
			AddLog(&m)
			//删除缓存
			err = utils.Redis.DEL(key)
			beego.Warn("删除缓存:", key)
			if err != nil {
				beego.Error(err)
			}
			return
		} else {
			cd := CabinetDetail{}
			err = o.Raw("select id,userID,`using`,store_time from cabinet_detail where id = ? limit 1;", m.Id).QueryRow(&cd)
			if err != nil {
				beego.Error(err)
				return
			}
			//非管理员操作
			//第二次关门
			if cd.UserID != "" && cd.Using == 2 {
				_, err = o.Raw("update cabinet_order_record set past_flag = 1 where customer_id = ? and cabinet_detail_id = ? and (past_flag is null or past_flag = 0) ;", cd.UserID, cd.Id).Exec()
				if err != nil {
					beego.Error(err)
					return
				}
				_, err = o.Raw("update cabinet_detail set userID = null,`using` = ?,store_time = ? where id = ? ;", 1, 0, cd.Id).Exec()
				if err == orm.ErrNoRows {
					return
				}
				if err != nil {
					beego.Error(err)
					return
				}
				//添加日志记录
				m := Log{
					CabinetDetailId: m.Id,
					User:            cd.UserID,
					Time:            time.Now(),
					Action:          "取",
				}
				AddLog(&m)
			}
		}
	}
	return
}
