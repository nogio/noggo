package util

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"io"
	"fmt"
)





//md5加密
func Md5(str string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(str))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

//md5加密文件
func Md5File(file string) string {
	if f, e := os.Open(file); e == nil {
		defer f.Close()

		h := md5.New()
		if _, e := io.Copy(h, f); e == nil {
			return fmt.Sprintf("%x", h.Sum(nil))
		}
	}
	return ""
}