package controllers

import (
	"github.com/skip2/go-qrcode"
	"fmt"
)

func CreateQrcodePng(qrcode_url string) string {
	//err := qrcode.WriteFile("http://blog.csdn.net/wangshubo1989", qrcode.Medium, 256, "qr.png")
	filename := "../static/qrcode.png"
	err := qrcode.WriteFile(qrcode_url, qrcode.Medium, 256, filename)
	if err != nil {
		fmt.Println("write error")
	}
	return filename
}
