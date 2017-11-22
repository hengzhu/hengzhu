package cache

import "github.com/astaxie/beego/cache"

var Bm cache.Cache
func init(){
	Bm, _ = cache.NewCache("memory", `{"key":"youcai","conn":":6379","dbNum":"0"`)
}