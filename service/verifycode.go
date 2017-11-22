package service

import (
	"strconv"
	"time"
	"github.com/Zeniubius/golang_utils/glog"
	"math/rand"
	"hengzhu/cache"
)



/**
* 根据位数生成验证码
* @param len 验证码长度
* @return
 */
func GetVerifyCode(len int) (string, error) {
	v, err := RandomNegative(len); if err != nil {
		glog.Info("gensmsVcode: %v", err)
		return "", err
	}
	return strconv.Itoa(v), nil
}

func VcodeVerify(mobile string, vcode string) bool {
	cacheCode := cache.Bm.Get(SMS_CACHE_KEY+mobile)
	if cacheCode == ""{
		glog.Warn("vcodeVerify<|>验证码为空<|>cacheCode:")
		return false
	}
	if cacheCode != vcode {
		glog.Warn("vcodeVerify<|>验证码错误<|>cacheCode:")
		return false
	}
	return true
}

func EmailVcodeVerify(mail string, vcode string) bool {
	cacheCode := cache.Bm.Get(MAIl_CACHE_KEY+mail)
	if cacheCode == ""{
		glog.Warn("vcodeVerify<|>邮件验证码为空<|>cacheCode:")
		return false
	}
	if cacheCode != vcode {
		glog.Warn("vcodeVerify<|>邮件验证码错误<|>cacheCode:")
		return false
	}
	return true
}

/**
*@Desc 生成随机数
*@Param
*@Return int64
*/
type RandomSizeError int

func (k RandomSizeError) Error() string {
	return "youcai_cryp/crypHelper: invalid key size " + strconv.Itoa(int(k)) + " | max size is 13"
}

func RandomNegative(size int) (int, error) {
	if size > 13 {
		return 0, RandomSizeError(size)
	}
	time_int := time.Now().UnixNano()
	r := rand.New(rand.NewSource(time_int))
	max_num := 1
	for i := 1; i < size + 1; i++ {
		max_num = max_num * 10
	}
	max_num--
	res_int := r.Intn(max_num)
	return res_int, nil
}
