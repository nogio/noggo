package utils

import (
	"encoding/hex"
	"os"
	"io"
	"fmt"
	"crypto/sha1"
)




//sha1加密
func Sha1(str string) string {
	md5Ctx := sha1.New()
	md5Ctx.Write([]byte(str))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}


//sha1加密文件
func Sha1File(file string) string {
	if f, e := os.Open(file); e == nil {
		defer f.Close()

		h := sha1.New()
		if _, e := io.Copy(h, f); e == nil {
			return fmt.Sprintf("%x", h.Sum(nil))
		}
	}
	return ""
}
