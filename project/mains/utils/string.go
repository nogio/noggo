package utils

import (
)
import (
	. "github.com/nogio/noggo/base"
	"time"
	"strings"
	"math/rand"
	"encoding/json"
	"fmt"
)



//密码加密格式
func Password(str string) string {
	return Md5(str)
}



func MakeVerifyCode(l int64) string {
	keys := []string{ "0","1","2","3","4","5","6","7","8","9" }
	codes := []string{}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i:=int64(0);i<l;i++ {
		num := r.Intn(len(keys))
		codes = append(codes, keys[num])
	}
	return strings.Join(codes, "")
}



//转成string
func String(val Any) string {

	sv := ""
	switch v:=val.(type) {
	case string:
		sv = v
	case Map,map[string]interface{}:
		d,e := json.Marshal(v)
		if e == nil {
			sv = string(d)
		} else {
			sv = "{}"
		}
	case []Map,[]map[string]interface{}:
		d,e := json.Marshal(v)
		if e == nil {
			sv = string(d)
		} else {
			sv = "[]"
		}
	case []int,[]int8,[]int16,[]int32,[]int64,[]float32,[]float64,[]string,[]bool,[]Any:
		d,e := json.Marshal(v)
		if e == nil {
			sv = string(d)
		} else {
			sv = "[]"
		}
	default:
		sv = fmt.Sprintf("%v", v)
	}

	return sv
}
