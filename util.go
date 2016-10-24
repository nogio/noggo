package noggo

import (
	"encoding/base64"
	"encoding/hex"
	"crypto/md5"
)

func encode64(value string) string {
	return base64.StdEncoding.EncodeToString([]byte(value))
}
func decode64(value string) string {
	d, e := base64.StdEncoding.DecodeString(value)
	if e == nil {
		return string(d)
	}
	return value
}


func encode5(str string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(str))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}