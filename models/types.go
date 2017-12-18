package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"github.com/astaxie/beego/orm"
	"time"
)

type Types struct {
	Id         int     `orm:"column(id);auto"`
	Name       string  `orm:"column(name);size(255);null" description:"类型名称"`
	Default    int     `orm:"column(default);null" description:"是否默认，1:默认，2:否"`
	ChargeMode int     `orm:"column(charge_mode);null" description:"计费方式，1:计次，2:计时"`
	TollTime   int     `orm:"column(toll_time);null" description:"收费时间，1:存物时，2:取物时"`
	Price      float64 `orm:"column(price);null" description:"价格，若方式为计次，则价格为每次存取物价格，若方式为计时，则价格为unit时间内价格"`
	Unit       int     `orm:"column(unit);null" description:"计时单位（分钟），当计费方式为计时时有"`
	CreateTime int64   `orm:"column(create_time);null" description:"创建时间"`

	CreateTimeFormated string `orm:"-"`
}

func (t *Types) TableName() string {
	return "type"
}

func init() {
	orm.RegisterModel(new(Types))
}

//
func AddTypesInfo(types []Types) {
	if len(types) == 0 {
		return
	}
	for i, typ := range types {
		types[i].CreateTimeFormated = time.Unix(int64(typ.CreateTime), 0).Format("2006-01-02 15:04:05")
	}
}

// 将id的类型设为默认，其他的改为非默认
func SetDefault(id int) (err error) {
	o := orm.NewOrm()
	_, err = o.Raw("UPDATE `type` SET `default`=2").Exec()
	if err != nil {
		return
	}

	_, err = o.Raw("UPDATE `type` SET `default`=1 WHERE id=?", id).Exec()

	return
}

// AddType insert a new Types into database and returns
// last inserted Id on success.
func AddType(m *Types) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

func GetAllTypes() (types []Types) {
	o := orm.NewOrm()
	o.QueryTable(new(Types)).All(&types)
	return
}

// GetTypeById retrieves Types by Id. Returns error if
// Id doesn't exist
func GetTypeById(id int) (v *Types, err error) {
	o := orm.NewOrm()
	v = &Types{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllType retrieves all Types matches certain condition. Returns empty list if
// no records exist
func GetAllType(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Types))
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

	var l []Types
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

// UpdateType updates Types by Id and returns error if
// the record to be updated doesn't exist
func UpdateTypeById(m *Types) (err error) {
	o := orm.NewOrm()
	v := Types{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteType deletes Types by Id and returns error if
// the record to be deleted doesn't exist
func DeleteType(id int) (err error) {
	o := orm.NewOrm()
	v := Types{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		if v.Default == 1 {
			return errors.New("默认类型不可删除")
		}
		var num int64
		if num, err = o.Delete(&Types{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func GetDefaultType() (t *Types) {
	o := orm.NewOrm()
	typ := Types{}
	o.Raw("select id from type where default = 1 limit 1 ;").QueryRow(&typ)
	t = &typ
	return
}
