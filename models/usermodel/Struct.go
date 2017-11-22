package usermodel

type UserLogin struct {
	UserName  string `json:"username" valid:"Required;Match(/^\w{6,20}$/)" form:"username"`
	Password  string `json:"password" valid:"Required;Match(/^(?=.*?[a-zA-Z])(?=.*?[0-9])([A-Za-z]|[0-9]|[^\w\s]*|[_]*){6,20}$/)" form:"password"`
	AutoLogin bool        `json:"autologin" form:"autologin"`
}

type UserReg struct {
	UserLogin
	RePassword string `valid:"Required;Match(/^(?=.*?[a-zA-Z])(?=.*?[0-9])([A-Za-z]|[0-9]|[^\w\s]*|[_]*){6,20}$/)" form:"repassword"`
	Mobile     string        ` form:"mobile"`
	SmsCode    string        `valid:"Required;Length(4)" form:"smscode"`
	CaptchaId  string        `valid:"Required;Match(/^(?=.*?[a-zA-Z])(?=.*?[0-9])([A-Za-z]|[0-9]|[^\w\s]*|[_]*){6,20}$/)" form:"captcha_id"`
}


type UserChangePwd struct {
	UserLogin
	RePassword string	`valid:"Required;MinSize(6);MaxSize(20);Match(/^[0-9a-zA-Z_]{6,20}$/)" form:"repassword"` //todo
	NewPassword string	`valid:"Required;MinSize(6);MaxSize(20);Match(/^[0-9a-zA-Z_]{6,20}$/)" form:"newpassword"`
}

type UserEmailBind struct {
	Email string `valid:"Required;Email" form:"bind"`
	Vcode string `valid:"Required;Length(6)" form:"vcode"`
}
type UserDetail struct {
	NickName string `form:"nickname"`
	TrueName string	`form:"truename"`
	Sex string	`form:"sex"`
	BornYear string	`form:"year"`
	BornMonth string	`form:"month"`
	BornDay string	`form:"day"`
	QQnum int	`form:"qqnum"`
	Mobile string	`form:"mobile"`
	WeChat string	`form:"wechat"`
	Add string	`form:"add"`
	ZipCode string	`form:"zipcode"`
}
