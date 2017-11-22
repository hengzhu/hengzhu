package utils

import (
	"fmt"
	"strconv"
	"strings"
	"log"
	"time"
	"encoding/base64"
)

/**
*@Desc 加密通信code
*@Param data 需要加密的数据，operation 操作，key 通信秘钥
*@Return string
*/
func UcAuthcode(data string, operation string, key string, expiry int64) string {
	ckeyl := 4
	if key == "" {
		return "key can not empty"
	}
	key = EncodeMd5(key)

	keya := EncodeMd5(Substr(key, 0, 16))
	keyb := EncodeMd5(Substr(key, 16, 16))
	keyc := ""
	//log.Printf("key is: %v\n", key)
	//log.Printf("keya is: %v\n", keya)
	//log.Printf("keyb is: %v\n", keyb)
	if operation == "DECODE" {
		keyc = Substr(data, 0, ckeyl)
	} else {
		keyc = Substr(EncodeMd5(strconv.FormatInt(time.Now().UnixNano() / 1e3, 10)), 28, 4)        //strconv.FormatInt(time.Now().UnixNano() / 1e3, 10)
	}
	cryptkey := keya + EncodeMd5((keya + keyc))
	keylen := len(cryptkey)

	//log.Printf("keyc is: %v\n", keyc)
	//log.Printf("cryptkey is: %v, len is %v\n", cryptkey, len(cryptkey))
	//log.Printf("keylen is: %v\n", keylen)
	//return ""

	datab := []byte(data)

	//log.Printf("DECODE before data is: %v", data)
	if operation == "DECODE" {
		log.Printf("DECODE Substrdata is: %v, len is %v", Substr(data, ckeyl, len(data)), len(Substr(data, ckeyl, len(data))))

		subdata := Substr(data, ckeyl, len(data))

		apd := strings.Repeat("=",4-(len([]byte(subdata)) % 4))
		subdata += apd
		//log.Printf("DECODE before base64: %v", subdata)
		//data = Base64DecodeString(Substr(data, ckeyl, len(data)))
		databb, _ := base64.StdEncoding.DecodeString(subdata)
		data = string(databb)

		//log.Printf("DECODE data is: %v", data)
		//log.Printf("DECODEMd5 data is: %v", EncodeMd5(data))
		//log.Printf("string lenth is: %v", len(data))
	} else {
		if expiry != 0 {
			expiry = expiry + time.Now().Unix()
		}
		data = fmt.Sprintf("%010d%s%s", expiry, Substr(EncodeMd5(data + keyb), 0, 16), data)
	}

	datalen := len(data)

	var box [256]int
	for key, _ := range box {
		box[key] = key
	}
	//log.Printf("box[100] is %v", box[100])

	var rndkey [256]int
	cryptkeyb := []byte(cryptkey)
	for i := 0; i < 256; i++ {
		rndkey[i] = int(cryptkeyb[i % keylen])        //返回cryptkey中的每个字符的ACII码并转为int型方便运算
	}
	//log.Printf("rndkey[255] is %v, len is %v", rndkey[255], len(rndkey))
	//log.Printf("rndkey is %v", rndkey)

	var j, tmp int
	for i := 0; i < 256; i++ {
		j = (j + box[i] + rndkey[i]) % 256;
		tmp = box[i]
		box[i] = box[j]
		box[j] = tmp
	}
	//log.Printf("box[120] is %v", box[120])
	//log.Printf("rndkey[120] is %v", rndkey[120])

	tmp = 0
	var datao, step int
	//log.Printf("datalen is %v", datalen)
	//log.Printf("box is %v, len is %v", box, len(box))
	//log.Printf("string is %v,len is %v , md5 is %v", data, len(data), EncodeMd5(data))
	//return ""
	var result = make([]byte, datalen)
	datab = []byte(data)
	for i, a, j := 0, 0, 0; i < datalen; i++ {
		a = (a + 1) % 256
		j = (j + box[a]) % 256
		tmp = box[a]
		box[a] = box[j]
		box[j] = tmp
		datao = int(datab[i])
		step = (datao ^ (box[(box[a] + box[j]) % 256 ]))
		result[i] = byte(step)
	}
	results := string(result)
	//log.Printf("result is: %v, len is: %v", results, len(results))
	//log.Printf("resultMd5 is: %v", EncodeMd5(results))
	//return ""
	if operation == "DECODE" {
		res, _ := strconv.ParseInt(Substr(results, 0, 10), 10, 64)
		//log.Printf("time is %v", res)
		res2 := Substr(results, 10, 16)
		res3 := Substr(EncodeMd5(Substr(results, 26, len(results)) + keyb), 0, 16)
		if (res == 0 || (res - time.Now().Unix()) > 0) && res2 == res3 {
			return Substr(results, 26, len(results));
		} else {
			log.Println("Error: can not DECODE")
			return ""
		}
	} else {
		return fmt.Sprintf("%s%s", keyc, strings.Replace(Base64EncodeString(results), "=", "", -1))
	}

}

