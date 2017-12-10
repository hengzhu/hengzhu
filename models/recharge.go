package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type Recharge struct {
	Id             int     `orm:"column(id);auto"`
	RechargeNo     string  `orm:"column(recharge_no);size(20);null" description:"充值单号"`
	UserId         uint    `orm:"column(user_id);null" description:"用户id"`
	Amount         float64 `orm:"column(amount);null;digits(16);decimals(2)" description:"金额"`
	Type           uint8   `orm:"column(type);null" description:"充值方式:1支付宝，2微信, 在支付回调后赋值"`
	Status         uint8   `orm:"column(status);null" description:"充值状态 1:正在支付 2:支付成功 3:超时关闭 4:手动关闭"`
	ThirdNo        string  `orm:"column(third_no);size(90);null" description:"第三方订单号"`
	ExtraInfo      string  `orm:"column(extra_info);null" description:"扩展信息"`
	CreateTime     uint    `orm:"column(create_time);null" description:"提交时间"`
	PayTime        uint    `orm:"column(pay_time);null" description:"支付时间"`
	SuccessTime    uint    `orm:"column(success_time);null" description:"充值成功时间"`
	Ip             string  `orm:"column(ip);size(15);null" description:"用户ip"`
	Wlan           string  `orm:"column(wlan);size(20);null" description:"用户mac地址"`
	Device         string  `orm:"column(device);size(64);null" description:"用户设备"`
	GameId         int     `orm:"column(game_id);null" description:"充值的游戏id"`
	ActivityAmount float32 `orm:"column(activity_amount);null" description:"活动期的虚拟币金额"`
	ActivityId     int     `orm:"column(activity_id);null" description:"关联活动id"`
}

func (t *Recharge) TableName() string {
	return "recharge"
}

func init() {
	orm.RegisterModel(new(Recharge))
}

// AddRecharge insert a new Recharge into database and returns
// last inserted Id on success.
func AddRecharge(m *Recharge) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetRechargeById retrieves Recharge by Id. Returns error if
// Id doesn't exist
func GetRechargeById(id int) (v *Recharge, err error) {
	o := orm.NewOrm()
	v = &Recharge{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllRecharge retrieves all Recharge matches certain condition. Returns empty list if
// no records exist
func GetAllRecharge(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Recharge))
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

	var l []Recharge
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

// UpdateRecharge updates Recharge by Id and returns error if
// the record to be updated doesn't exist
func UpdateRechargeById(m *Recharge) (err error) {
	o := orm.NewOrm()
	v := Recharge{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteRecharge deletes Recharge by Id and returns error if
// the record to be deleted doesn't exist
func DeleteRecharge(id int) (err error) {
	o := orm.NewOrm()
	v := Recharge{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Recharge{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
