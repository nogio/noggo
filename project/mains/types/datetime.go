package types

import (
	"github.com/nogio/noggo"
	. "github.com/nogio/noggo/base"
	"time"
)

func init() {


	noggo.Mapping.Type("datetime", Map{
		"name": "日期时间", "text": "日期时间",
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
			switch v := value.(type) {
			case time.Time:
				return v
			case *time.Time:
				return v
			case string:
				dt,err := time.Parse("2006-01-02 15:04:05", v)
				if err == nil {
					return dt
				} else {
					return v
				}
			}
			return value
		},
	})

	noggo.Mapping.Type("[datetime]", Map{
		"name": "日期时间数组", "text": "日期时间数组",
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