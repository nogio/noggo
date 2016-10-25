package types

import (
	"github.com/nogio/noggo"
	. "github.com/nogio/noggo/base"
	"../utils"
	"fmt"
)

func init() {

	noggo.Mapping.Type("password", Map{
		"name": "密码", "text": "密码",
		"valid": func(value Any, config Map) bool {

			if value == nil {
				return false
			}

			switch v := value.(type) {
			case string: {
				if v == "" {
					return false
				}
			}
			}
			return true
		},
		"value": func(value Any, config Map) Any {
			switch v := value.(type) {
			case string:
				if noggo.Const.Valid(v, "password") {
					return v
				} else {
					return utils.Password(v)
				}
			}
			return fmt.Sprintf("%v", value)
		},
	})


}