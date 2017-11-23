package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type Cabinet struct {
	Id        int    `json:"id" orm:"column(id);auto"`
	CabinetID string `json:"cabinet_id" orm:"column(cabinet_ID);size(255);null" description:"柜子id"`
	Isonline  int8   `json:"isonline" orm:"column(isonline);null" description:"是否在线,1:在线，2:不在线"`
	TypeId    int    `json:"type_id" orm:"column(type_id);null" description:"柜子类型id"`
	Address   string `json:"address" orm:"column(address);null" description:"柜子位置"`
	Number    string `json:"number" orm:"column(number);null" description:"编号"`
	Desc      string `json:"desc" orm:"column(desc);null" description:"备注"`
}

func (t *Cabinet) TableName() string {
	return "cabinet"
}

func init() {
	orm.RegisterModel(new(Cabinet))
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
