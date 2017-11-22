package usermodel

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type UserInfo struct {
	Id            int       `orm:"column(user_id);pk"`
	NickName      string    `orm:"column(nick_name);size(20);null"`
	Sex           string    `orm:"column(sex);type(enum)"`
	Birth         time.Time `orm:"column(birth);type(date)"`
	Headpic       string    `orm:"column(headpic);size(64);null"`
	RealName      string    `orm:"column(real_name);size(32);null"`
	IdentityCard  string    `orm:"column(identity_card);size(20);null"`
	UserTel       string    `orm:"column(user_tel);size(20);null"`
	UserMobile    string    `orm:"column(user_mobile);size(11);null"`
	UserAddress   string    `orm:"column(user_address);size(100);null"`
	UserPostcode  string    `orm:"column(user_postcode);size(6);null"`
	UserJob       int8      `orm:"column(user_job)"`
	UserSalary    int8      `orm:"column(user_salary)"`
	UserEducation int8      `orm:"column(user_education)"`
	UserMarital   int8      `orm:"column(user_marital)"`
	Area          string    `orm:"column(area);size(20);null"`
	Exp           uint      `orm:"column(exp)"`
	Level         uint8     `orm:"column(level)"`
	Nosign        string    `orm:"column(nosign)"`
	CreateTime    uint      `orm:"column(create_time)"`
}

func (t *UserInfo) TableName() string {
	return "user_info"
}

func init() {
	orm.RegisterModel(new(UserInfo))
}

// AddUserInfo insert a new UserInfo into database and returns
// last inserted Id on success.
func AddUserInfo(m *UserInfo) (id int64, err error) {
	o := orm.NewOrm(); o.Using(Dbname)
	id, err = o.Insert(m)
	return
}

func InitUserInfo(uid int, nick string, mobile string, create_time uint) (num int64, err error) {
	o := orm.NewOrm(); o.Using(Dbname)
	res, err := o.Raw("INSERT INTO user_info(user_id, nick_name, user_mobile, create_time) VALUES(?,?,?,?)", uid, nick, mobile, create_time).Exec(); if err == nil {
		num, _ := res.RowsAffected()
		return num, nil
	}
	return 0, nil
}



func UpdateUserInfo(u *UserInfo) (int64, error) {
	o := orm.NewOrm();o.Using(Dbname)
	user := make(orm.Params)
	t := time.Time{}
	if len(u.NickName) > 0 {
		user["nick_name"] = u.NickName
	}
	if len(u.Sex) > 0 {
		user["sex"] = u.Sex
	}
	if u.Birth != t  {
		user["birth"] = u.Birth
	}
	if len(u.Headpic) > 0 {
		user["headpic"] = u.Headpic
	}
	if len(u.RealName) > 0 {
		user["real_name"] = u.RealName
	}
	if len(u.UserMobile) > 0 {
		user["user_mobile"] = u.UserMobile
	}
	if len(u.UserAddress) > 0 {
		user["user_address"] = u.UserAddress
	}
	if len(u.UserPostcode) > 0 {
		user["user_postcode"] = u.UserPostcode
	}
	if len(user) == 0 {
		return 0, errors.New("update field is empty")
	}
	num, err := o.QueryTable("user_info").Filter("Id", u.Id).Update(user)
	return num, err
}


// GetUserInfoById retrieves UserInfo by Id. Returns error if
// Id doesn't exist
func GetUserInfoById(id int) (v *UserInfo, err error) {
	o := orm.NewOrm();o.Using(Dbname)
	v = &UserInfo{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllUserInfo retrieves all UserInfo matches certain condition. Returns empty list if
// no records exist
func GetAllUserInfo(query map[string]string, fields []string, sortby []string, order []string,
offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(UserInfo))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Filter(k, v)
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

	var l []UserInfo
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

// UpdateUserInfo updates UserInfo by Id and returns error if
// the record to be updated doesn't exist
func UpdateUserInfoById(m *UserInfo) (err error) {
	o := orm.NewOrm()
	v := UserInfo{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteUserInfo deletes UserInfo by Id and returns error if
// the record to be deleted doesn't exist
func DeleteUserInfo(id int) (err error) {
	o := orm.NewOrm()
	v := UserInfo{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&UserInfo{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
