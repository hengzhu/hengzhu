// Copyright © 2016年 kuaifazs.com. All rights reserved.
// 
// @Author: wuchengshuang@kuaifzs.com
// @Date: 2016/12/6
// @Version: -
// @Desc: -

package utils

import (
    "qiniupkg.com/api.v7/kodo"
    "qiniupkg.com/api.v7/conf"
    "qiniupkg.com/api.v7/kodocli"
    "fmt"
)

func QnUpload(accesskey, secret_key, bucket, filepath, filename string) (string, error) {
    //初始化AK，SK
    conf.ACCESS_KEY = accesskey
    conf.SECRET_KEY = secret_key

    //创建一个Client
    c := kodo.New(0, nil)

    key := RandString(6)  + "_" + filename

    //设置上传的策略
    policy := &kodo.PutPolicy{
        Scope:   bucket + ":" + key,
        //设置Token过期时间
        Expires: 3600,
    }
    //生成一个上传token
    token := c.MakeUptoken(policy)

    //构建一个uploader
    zone := 0
    uploader := kodocli.NewUploader(zone, nil)

    var ret kodo.PutRet
    res := uploader.PutFile(nil, &ret, token, key, filepath, nil)
    //打印返回的信息
    //fmt.Println(ret)
    //打印出错信息
    if res != nil {
        fmt.Println("io.Put failed:", res)
        return "", res
    }
    return key, nil
}
