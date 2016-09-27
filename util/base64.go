package util

import (
	"encoding/base64"
)



//base64加密
func Encode64(value string) string {
	return base64.StdEncoding.EncodeToString([]byte(value))
}
//base64解密
func Decode64(value string) string {
	d, e := base64.StdEncoding.DecodeString(value)
	if e == nil {
		return string(d)
	}
	return value
}


