package usermodel

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
	"regexp"
	"hengzhu/admin/src/lib"
	"time"
	"qiniupkg.com/x/log.v7"
)

type User struct {
	Id           int    `orm:"column(user_id);auto"`
	Mail         string `orm:"column(mail);size(50);null"`
	Mobile       string `orm:"column(mobile);size(11);null"`
	UserName     string `orm:"column(user_name);size(50);null"`
	NickName     string `orm:"column(nick_name);size(20)"`
	Password     string `orm:"column(password);size(32)"`
	CreateTime   uint   `orm:"column(create_time)"`
	UserType     int8   `orm:"column(user_type)"`
	RegisterWlan string `orm:"column(register_wlan);size(17);null"`
	GameId       int32  `orm:"column(game_id)"`
	AccountMark  int8   `orm:"column(account_mark)"`
	Initialpwd   string `orm:"column(initialpwd);size(12)"`
	Solidify     int8   `orm:"column(solidify)"`
	ChineseName  string `orm:"column(chinese_name);size(20)"`
	IdCardNo     string `orm:"column(id_card_no);size(20)"`
}

func (t *User) TableName() string {
	return "user"
}

func init() {
	//o := orm.NewOrm()
	//o.Using(Dbname)
	orm.RegisterModel(new(User))

}

// AddUser insert a new User into database and returns
// last inserted Id on success.
func AddUser(m *User) (id int64, err error) {
	o := orm.NewOrm()
	o.Using(Dbname)
	id, err = o.Insert(m)
	return
}

// GetUserById retrieves User by Id. Returns error if
// Id doesn't exist
func GetUserById(id int) (v *User, err error) {
	o := orm.NewOrm(); o.Using(Dbname)
	v = &User{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

//@Desc:
//
//@Param:
//
//@Return:
//
func GetUserByMulti(username string, ty string) (user User) {
	user = User{
		UserName:username,
		Mail:username,
		Mobile:username}

	o := orm.NewOrm(); o.Using(Dbname)
	switch ty {
	case "mail":
		o.Read(&user, "mail")
		return
	case "mobile":
		o.Read(&user, "mobile")
		return
	case "username":
		o.Read(&user, "user_name")
		return
	default:
		return
	}
}

func UpdateUser(u *User) (int64, error) {
	o := orm.NewOrm(); o.Using(Dbname)
	user := make(orm.Params)
	if len(u.UserName) > 0 {
		user["user_name"] = u.UserName
	}
	if len(u.NickName) > 0 {
		user["nick_name"] = u.NickName
	}
	if len(u.Mail) > 0 {
		user["mail"] = u.Mail
	}
	if len(u.Mobile) > 0 {
		user["mobile"] = u.Mobile
	}
	if len(u.Password) > 0 {
		user["password"] = lib.Pwdhash(lib.Pwdhash(u.Password))
	}
	if len(user) == 0 {
		return 0, errors.New("update field is empty")
	}
	num, err := o.QueryTable("user").Filter("Id", u.Id).Update(user)
	return num, err
}

func SyncUser(uid int, username string, password string, email string, ) (User, error) {
	var usertype, accountmark int8 = 1, 0
	u := &User{
		Id:uid,
		UserName:username,
		Password:lib.Pwdhash(lib.Pwdhash(password)),
		Mail:email,
		CreateTime:uint(time.Now().Unix()),
		UserType:usertype,
		NickName:username,
		AccountMark:accountmark,
	}
	if strings.Contains(email, "youcai") {
		u.Mail = ""
	}
	id, err := AddUser(u); if err == nil && id > 0 {
		InitUserInfo(uid, u.UserName, u.Mobile, uint(time.Now().Unix()));
		return *u, nil
	}
	return *u, err
}

//@Desc:检查Ucenter用户邮箱绑定信息并同步
//
//@Param:
//
//@Return:
//
func CheckSyncUserEmail(email string,u User) (User) {
	if u.Mail != email {
		m := User{
			Id:u.Id,
			Mail:email,
		}
		u.Mail = email
		if strings.Contains(email, "youcai") {
			m.Mail = ""
			u.Mail = ""
		}
		UpdateUser(&m)

	}
	return u
}

// GetAllUser retrieves all User matches certain condition. Returns empty list if
// no records exist
func GetAllUser(query map[string]string, fields []string, sortby []string, order []string,
offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	o.Using(Dbname)
	qs := o.QueryTable(new(User))
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

	var l []User
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

// UpdateUser updates User by Id and returns error if
// the record to be updated doesn't exist
func UpdateUserById(m *User) (err error) {
	o := orm.NewOrm()
	v := User{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteUser deletes User by Id and returns error if
// the record to be deleted doesn't exist
func DeleteUser(id int) (err error) {
	o := orm.NewOrm()
	v := User{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&User{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}




//@Desc 判断用户名的类型
//
//@Param
//
//@Return tp: 类型; err: an error if not nil
func CheckUserNameType(username string) (tp string, err error) {
	if strings.Contains(username, "@") {
		return "mail", nil
	}
	r2, err := regexp.Compile(`^(13[0-9]|15[0-3,5-9]|18[0-9]|14[57]|17[0-9])\d{8}$`); if err != nil {
		return "", err
	}
	if r2.Match([]byte(username)) {
		return "mobile", nil
	}
	return "username", nil
}

func CheckUserExist(column string, filter ...interface{}) bool {
	o := orm.NewOrm()
	o.Using(Dbname)
	return o.QueryTable("user").Filter(column, filter...).Exist()
}

//@Desc: 获取用户星级和绑定状态
//
//@Param:
//
//@Return:star 星级，mobile 手机绑定状态，email 邮箱绑定状态
//
func GetUserStarAndBindIfo(m User) (star int, mobile bool, email bool) {
	star = 0
	if len(m.Mobile) > 0 {
		star += 2
		mobile = true
	}
	if len(m.Mail) > 0 {
		star += 2
		email = true
	}
	log.Infof("Len(mobile): %v, Len(Mail):%v",len(m.Mobile),len(m.Mail))
	return
}
