package types

import (
	"github.com/nogio/noggo"
	. "github.com/nogio/noggo/base"
	"time"
)

func init() {


	noggo.Mapping.Type("date", Map{
		"name": "日期时间", "text": "日期时间",
		"type": []string{ "date" },
		"valid": func(value Any, config Map) bool {

			switch v := value.(type) {
			case time.Time:
				return true
			case *time.Time:
				return true
			case string:
				return noggo.Const.Valid(v, "date")
			}

			return false;
		},
		"value": func(value Any, config Map) Any {
			return value
		},
	})

	noggo.Mapping.Type("[date]", Map{
		"name": "日期时间数组", "text": "日期时间数组",
		"type": []string{ "[date]" },
		"valid": func(value Any, config Map) bool {

			switch v := value.(type) {
			case []time.Time:
				return true
			case *[]time.Time:
				return true
			case string:
				return noggo.Const.Valid(v, "date")
			}

			return false;
		},
		"value": func(value Any, config Map) Any {
			switch v := value.(type) {
			case []time.Time:
				return v
			default:
				return v
			}
		},
	})


}