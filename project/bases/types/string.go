package types

import (
	"github.com/nogio/noggo"
	. "github.com/nogio/noggo/base"
	"fmt"
)

func init() {

	//注册Mapping模块的类型
	noggo.Mapping.Type("string", Map{
		"name": "字符串", "text": "字符串类型",
		"valid": func(value Any, config Map) bool {
			return false
		},
		"value": func(value Any, config Map) Any {
			return fmt.Sprintf("%s", value)
		},
	})
}
