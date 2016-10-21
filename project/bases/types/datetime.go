package types

import (
	"github.com/nogio/noggo"
	. "github.com/nogio/noggo/base"
	"time"
)

func init() {


	noggo.Mapping.Type("datetime", Map{
		"name": "日期时间", "text": "日期时间",
		"type": []string{ "datetime" },
		"valid": func(value Any, config Map) bool {

			switch v := value.(type) {
			case time.Time:
				return true
			case *time.Time:
				return true
			case string:
				return noggo.Const.Valid(v, "datetime")
			}

			return false;
		},
		"value": func(value Any, config Map) Any {
			return value
		},
	})

	noggo.Mapping.Type("[datetime]", Map{
		"name": "日期时间数组", "text": "日期时间数组",
		"type": []string{ "[datetime]" },
		"valid": func(value Any, config Map) bool {

			switch v := value.(type) {
			case []time.Time:
				return true
			case *[]time.Time:
				return true
			case string:
				return noggo.Const.Valid(v, "datetime")
			}

			return false;
		},
		"value": func(value Any, config Map) Any {
			switch v := value.(type) {
			case []time.Time:
				return v
			case *[]time.Time:
				return v
			default:
				return v
			}
		},
	})


}