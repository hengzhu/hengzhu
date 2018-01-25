package models

import (
	"fmt"
	"strings"

	"github.com/astaxie/beego/orm"
	"strconv"
	"time"
	"math/rand"
	"github.com/astaxie/beego"
	"errors"
)

type CabinetOrderRecord struct {
	Id              int     `orm:"column(id);auto"`
	OrderNo         string  `orm:"column(order_no)" description:"内部生成的订单号"`
	CustomerId      string  `orm:"column(customer_id);size(255);null" description:"顾客id 微信 openid 支付宝？"`
	PayType         int8    `orm:"column(pay_type)" description:"1 微信 2支付宝 3？"`
	ThirdOrderNo    string  `orm:"column(third_order_no);size(255);null" description:"第三方支付id"`
	CabinetDetailId int     `orm:"column(cabinet_detail_id)"`
	Fee             float64 `orm:"column(fee)" description:"钱数"`
	CreateDate      int     `orm:"column(create_date)"`
	PayDate         int     `orm:"column(pay_date);null"`
	IsPay           int8    `orm:"column(is_pay)" description:"是否支付 0 未支付 1已经支付"`
	ActionType      int8    `orm:"column(action_type)" description:"1.重复支付存物"`
	PastFlag        int8    `orm:"column(past_flag)" description:"1.过去的支付记录"`
}

func (t *CabinetOrderRecord) TableName() string {
	return "cabinet_order_record"
}

func init() {
	orm.RegisterModel(new(CabinetOrderRecord))
}

// AddCabinetOrderRecord insert a new CabinetOrderRecord into database and returns
// last inserted Id on success.
func AddCabinetOrderRecord(m *CabinetOrderRecord) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetCabinetOrderRecordById retrieves CabinetOrderRecord by Id. Returns error if
// Id doesn't exist
func GetCabinetOrderRecordById(id int) (v *CabinetOrderRecord, err error) {
	o := orm.NewOrm()
	v = &CabinetOrderRecord{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

//生成订单号
func CreateOrderNo() (no string, err error) {
	no = RandString(10)
	return
}

//生成随机字符串
func RandString(length int) string {
	rand.Seed(time.Now().UnixNano())
	rs := make([]string, length)
	for start := 0; start < length; start++ {
		t := rand.Intn(3)
		if t == 0 {
			rs = append(rs, strconv.Itoa(rand.Intn(10)))
		} else if t == 1 {
			rs = append(rs, string(rand.Intn(26)+65))
		} else {
			rs = append(rs, string(rand.Intn(26)+97))
		}
	}
	return strings.Join(rs, "")
}

//修改订单状态为支付完成
func UpdateOrderSuccessByNo(third_order_no string, order_no string, openid string) (cd *CabinetDetail, err error) {
	o := orm.NewOrm()
	v := CabinetOrderRecord{OrderNo: order_no}
	cdd := &CabinetDetail{}
	cor := &CabinetOrderRecord{}
	if err = o.Read(&v, "order_no"); err == nil {
		//更新以前重新校验柜子是否可用
		cor, err = GetCabinetOrderByDetailIdAndOpenId(v.CabinetDetailId)
		if err == nil && cor.IsPay == 1 {
			cdd, err = GetFreeDoorByCabinetId(cdd.CabinetId)
			if err != nil {
				beego.Warn("分配柜子失败")
				return
			}
			v.CabinetDetailId = cdd.Id
		} else if err != nil && err != orm.ErrNoRows {
			beego.Warn("系统错误")
			return
		} else {
			var num int64
			v.IsPay = 1
			v.CustomerId = openid
			v.PayDate = int(time.Now().Unix())
			v.ThirdOrderNo = third_order_no
			if num, err = o.Update(&v, "customer_id", "is_pay", "pay_date", "third_order_no"); err == nil {
				fmt.Println("Number of records updated in database:", num)
			}
		}
	}
	//先存后支付下单形式
	//已经查到该用户在用
	c := CabinetDetail{UserID: openid, Using: 2, UseState: 1}
	if err = o.Read(&c, "userID", "using", "use_state"); err == nil {
		cor := CabinetOrderRecord{}
		if err = o.Raw("select id from cabinet_order_record where customer_id = ? and cabinet_detail_id = ? and (past_flag is null or past_flag = 0) limit 1 ;", openid, v.CabinetDetailId).QueryRow(&cor); err == nil {
			err = errors.New("[重复存物]: " + openid)
			beego.Error(err)
			o.Raw("update cabinet_order_record set past_flag = 1 ,action_type = 1 where order_no = ? and third_order_no = ?;", order_no, third_order_no).Exec()
			return nil, err
		}
		cd = &c
		return
	}
	cd, err = GetCabinetDetailById(v.CabinetDetailId)
	if err != nil {
		beego.Error(err)
	}
	return

	//更新该柜子的门为使用中>>关门时才绑定

	//cd = cdd
	//err, _ = BindOpenIdForCabinetDoor(openid, cd.Id)
	//if err != nil {
	//	beego.Error(err)
	//}
}

func GetCabinetOrderByDetailIdAndOpenId(detailId int) (cor *CabinetOrderRecord, err error) {
	o := orm.NewOrm()
	v := CabinetOrderRecord{}
	err = o.Raw("select id, is_pay from cabinet_order_record where cabinet_detail_id = ? and (past_flag = 0 or past_flag is null) limit 1;", detailId).QueryRow(&v)
	cor = &v
	return
}
