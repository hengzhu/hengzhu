package cache

type Item struct {
	Appid        string
	Type         string
	Name         string
	Url          string
	Ip           string
	Viewprourl   string
	Apifilename  string
	Charset      string
	Synlogin     int
	Recvnote     string
	Extra        string
	Tagtemplates string
	Allowips     string
	AuthKey      string
}

var Apps = []Item{
	{

	},
	{
		Appid:"1",
		Type:"OTHER",
		Name:"youcai_web",
		Url:"http://test.youcaibbs.com",
		Ip:"",
		Viewprourl:"",
		Apifilename:"uc",
		Charset:"",
		Synlogin:1,
		Recvnote:"1",
		Extra:"",
		Tagtemplates:"",
		Allowips:"",
		AuthKey:"LcXbK65dY8Kc04beJ2h7P9o9afI6Y2o8o8LfU3p791kfL6F6F2w3scxff4x9Rca5",
	},
	{
		Appid:"2",
		Type:"OTHER",
		Name:"youcai_web",
		Url:"http://test.youcaigw.com:8080",
		Ip:"",
		Viewprourl:"",
		Apifilename:"uc",
		Charset:"",
		Synlogin:1,
		Recvnote:"1",
		Extra:"",
		Tagtemplates:"",
		Allowips:"",
		AuthKey:"LcXbK65dY8Kc04beJ2h7P9o9afI6Y2o8o8LfU3p791kfL6F6F2w3scxff4x9Rca5",
	},
}