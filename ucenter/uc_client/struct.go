package uc_client

type UserLoginRes struct {
	Result     int
	UserName     string
	Password	string
	Email	string
	IsRepeatName string
}

type Root struct {
	Items []Item `xml:"item"`
}

type Item struct {
	Id int `xml:"id,attr"`
	Data string `xml:",chardata"`
}
