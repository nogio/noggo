package native

import (
	"github.com/nogio/noggo"
	. "github.com/nogio/noggo/base"
)


func init() {

	//注册base64加解密
	noggo.Mapping.Crypto("base64", Map{
		"name": "BASE64加解密", "text": "BASE64加解密",
		"encode": func(value Any) Any {
			return encode64(value.(string))
		},
		"decode": func(value Any) Any {
			return decode64(value.(string))
		},
	})

}

