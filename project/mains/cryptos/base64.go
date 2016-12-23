package types

import (
	"github.com/nogio/noggo"
	. "github.com/nogio/noggo/base"
	"../utils"
)

func init() {

	//注册mapping模块的加解密方法
	noggo.Mapping.Crypto("base64", Map{
		"name": "BASE64加解密", "text": "BASE64加解密",
		"encode": func(value Any) Any {
			return utils.Encode64(utils.String(value))
		},
		"decode": func(value Any) Any {
			return utils.Decode64(utils.String(value))
		},
	})
}
