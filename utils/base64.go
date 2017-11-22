package utils

import (
"encoding/base64"
)

func Base64Encode(src []byte) []byte {
	buf := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(buf, src)

	return buf
}

func Base64Decode(src []byte) (out []byte) {
	buf := make([]byte, base64.StdEncoding.DecodedLen(len(src)))
	base64.StdEncoding.Decode(buf, src)

	return buf
}

func Base64DecodeString(src string) (out string) {
	if src == "" {
		return
	}
	return string(Base64Decode([]byte(src)))
}

func Base64EncodeString(src string) (out string) {
	if src == "" {
		return
	}
	return string(Base64Encode([]byte(src)))
}

