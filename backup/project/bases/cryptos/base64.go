package types

import (
	"github.com/nogio/noggo"
	. "github.com/nogio/noggo/base"
	"../utils"
)

func init() {

	noggo.Mapping.Crypto("base64", Map{
		"name": "BASE64加解密", "text": "BASE64加解密",
		"encode": func(value Any) Any {
			return utils.Encode64(value.(string))
		},
		"decode": func(value Any) Any {
			return utils.Decode64(value.(string))
		},
	})
}