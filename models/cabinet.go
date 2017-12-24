package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Cabinet struct {
	Id         int       `orm:"column(id);auto"`
	CabinetID  string    `orm:"column(cabinet_ID);size(255);null" description:"柜子id"`
	TypeId     int       `orm:"column(type_id);null" description:"柜子计费类型id，初始化时为默认类型"`
	Address    string    `orm:"column(address);null" description:"柜子位置"`
	Number     string    `orm:"column(number);null" description:"编号"`
	Desc       string    `orm:"column(desc);null" description:"备注"`
	CreateTime time.Time `orm:"column(create_time);type(timestamp);null;auto_now_add" description:"创建时间"`
	LastTime   time.Time `orm:"column(last_time);type(timestamp);null" description:"最后一次上报时间"`
	UpdateTime time.Time `orm:"column(update_time);type(timestamp);null" description:"更新时间"`

	TypeName string          `orm:"-" description:"类型名称"`
	IsOnline string          `orm:"-" description:"是否在线"`
	Doors    int             `orm:"-" description:"门数"`
	OnUse    int             `orm:"-" description:"使用中数量"`
	Close    int             `orm:"-" description:"关闭状态门数量"`
	Detail   []CabinetDetail `orm:"-" description:"柜子的门详情"`
}

func (t *Cabinet) TableName() string {
	return "cabinet"
}

func init() {
	orm.RegisterModel(new(Cabinet))
}

func AddOtherInfo(cabinets *[]Cabinet) {
	if cabinets == nil || len(*cabinets) == 0 {
		return
	}

	for i, cabinet := range *cabinets {
		AddInfo(&cabinet)
		(*cabinets)[i] = cabinet
	}
}

func AddInfo(cabinet *Cabinet) {
	typ, err := GetTypeById(cabinet.TypeId)
	if err == nil {
		cabinet.TypeName = typ.Name
	}

	cabinet.Doors = GetTotalDoors(cabinet.Id)
	cabinet.OnUse = GetTotalOnUse(cabinet.Id)
	cabinet.Close = GetTotalClose(cabinet.Id)

	cabinet.IsOnline = "是"
	if time.Now().Unix()-cabinet.LastTime.Unix() > 90 {
		cabinet.IsOnline = "否"
	}
}

// 给柜子附加上门详情
func AddDetails(cabinet *Cabinet) {
	details, err := GetDetailsByCabinetId(cabinet.Id)
	if err != nil {
		return
	}
	cabinet.Detail = details
	return
}

// AddCabinet insert a new Cabinet into database and returns
// last inserted Id on success.
func AddCabinet(m *Cabinet) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetCabinetById retrieves Cabinet by Id. Returns error if
// Id doesn't exist
func GetCabinetById(id int) (v *Cabinet, err error) {
	o := orm.NewOrm()
	v = &Cabinet{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllCabinet retrieves all Cabinet matches certain condition. Returns empty list if
// no records exist
func GetAllCabinet(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Cabinet))
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

	var l []Cabinet
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

// UpdateCabinet updates Cabinet by Id and returns error if
// the record to be updated doesn't exist
func UpdateCabinetById(m *Cabinet) (err error) {
	o := orm.NewOrm()
	v := Cabinet{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteCabinet deletes Cabinet by Id and returns error if
// the record to be deleted doesn't exist
func DeleteCabinet(id int) (err error) {
	o := orm.NewOrm()
	v := Cabinet{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Cabinet{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func GetCabinetAndDoorByUserId(user_id string) (cid int, door_no int, cdid int, err error) {
	o := orm.NewOrm()
	cd := CabinetDetail{}
	err = o.Raw("select cabinet_id,door from cabinet_detail where id = (select cabinet_detail_id from cabinet_order_record where customer_id = '?' and is_pay = 1 limit 1);", user_id).QueryRow(&cd)
	if err != nil {
		return
	}
	cdid = cd.Id
	cid = cd.CabinetId
	door_no = cd.Door
	return
}

func GetCabinetByMac(cabinet_mac string) (v *Cabinet, err error) {
	o := orm.NewOrm()
	v = &Cabinet{CabinetID: cabinet_mac}
	if err = o.Read(v, "cabinet_ID"); err == nil {
		return v, nil
	}
	return nil, err
}

// 检查是否需要初始化柜子
func CheckIfAdd(cabinet_mac string) (bool) {
	o := orm.NewOrm()
	v := &Cabinet{CabinetID: cabinet_mac}
	if err := o.Read(v, "cabinet_ID"); err == nil && v.Id > 0 {
		return false
	}

	return true
}

func GetCabinetQueues() ([]string) {
	o := orm.NewOrm()
	result := []Cabinet{}
	queue := []string{}

	sql := "SELECT cabinet_ID FROM cabinet"
	o.Raw(sql).QueryRows(&result)

	for _, res := range result {
		queue = append(queue, res.CabinetID)
	}

	return queue
}
