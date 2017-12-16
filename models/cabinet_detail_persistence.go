package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"hengzhu/tool/payment"
	"time"
	"strconv"
)

func GetIdleDoorByCabinetId(cid int64) CabinetDetail { //通过柜子id得到一个空闲的门
	o := orm.NewOrm()
	cabd := CabinetDetail{}
	err := o.Raw("select * from cabinet_detail where cabinet_id= ? and `using` = 1 limit 1;", cid).QueryRow(&cabd)
	if err != nil {
		beego.Error("[CabinetOrder]: GetIdleDoorByCabinetId err in select:", err)
		return cabd
	}
	return cabd
}

func CreateNewWxOrder(order payment.WXUnifiedorderRequest, cid int) bool {
	o := orm.NewOrm()
	now := time.Now().Unix()
	fee, _ := strconv.Atoi(order.TotalFee)
	_, err := o.Raw(`insert into cabinet_order_record set order_no=?,pay_type=?,cabinet_detail_id=?,fee=?,create_date=?`,
		order.OutTradeNo, 1, cid, fee, now).Exec() //fee虽然是字符串但是可以正确插入
	if err != nil {
		beego.Error("[CabinetOrder]: CreateNewWxOrder err in insert:", err)
		return false
	}
	return true
}

func GetCabDIdByOrderNo(orderNo string) int64 {
	//根据订单得到商品id
	o := orm.NewOrm()
	var did int64
	err := o.Raw(`select cabinet_detail_id from cabinet_order_record where order_no=? limit 1`, orderNo).QueryRow(&did)
	if err != nil {
		beego.Error("[CabinetOrder] VerifyWxOrderNo err in select:", err)
		return 0
	}

	return did
}
func WxPaySuccess(res payment.WXPayResultNotifyArgs, detailId int64) bool {
	/*
	更新订单信息  1修改当前柜子门的状态 2插入用户使用记录 3这应该是一个原子操作 所以要使用事务 注意数据库引擎innodb
	因为1有关跟柜子的通信,我只先写2 3
	*/
	o := orm.NewOrm()
	txerr := o.Begin()
	_, err := o.Raw(`update cabinet_order_record set customer_id=?,third_order_no=?,pay_date=?,is_pay=1`,
		res.OpenId, res.TransactionId, payment.DealStringTime(res.TimeEnd)).Exec()
	if err != nil {
		beego.Error("[CabinetOrder] WxPaySuccess err in update:", err)
		txerr = o.Rollback()
		return false
	}
	now := time.Now().Unix() //start_date 我还是觉得应该从开门算起
	_, err = o.Raw(`insert into cabinet_use_log set customer_id=?,cabinet_detail_id=?,start_date=?`,
		res.OpenId, detailId, now).Exec()
	if err != nil {
		beego.Error("[CabinetOrder] WxPaySuccess err in insert into log:", err)
		txerr = o.Rollback()
		return false
	}
	txerr = o.Commit()
	beego.Debug("[CabinetOrder] run over and txerr:", txerr)
	return true

}

func GetWxUserInUseByOpenId(openId string) int64 {
	o := orm.NewOrm()
	var detailId int64
	err := o.Raw(`select id from cabinet_detail where userID=? and open_state=1 and using=2`, openId).QueryRow(&detailId)
	if err != nil {
		beego.Error("[CabinetOrder]: GetWxUserInUseByOpenId err in select:", err)
		return 0
	}
	return detailId
}
