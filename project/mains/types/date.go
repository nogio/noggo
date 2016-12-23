package types

import (
	"github.com/nogio/noggo"
	. "github.com/nogio/noggo/base"
	"time"
)

func init() {


	noggo.Mapping.Type("date", Map{
		"name": "日期时间", "text": "日期时间",
		"valid": func(value Any, config Map) bool {

			switch v := value.(type) {
			case time.Time:
				return true
			case *time.Time:
				return true
			case int64:
				return true
			case string:
				return noggo.Const.Valid(v, "date")
			}

			return false;
		},
		"value": func(value Any, config Map) Any {
			switch v := value.(type) {
			case int64:
				return time.Unix(v, 0)
			case string:
				lay := "2006-01-02"
				if len(v) == 8 {
					lay = "20060102"
				} else if len(v) == 10 {
					lay = "2006-01-02"
				} else {
					lay = "2006-01-02"
				}

				dt,err := time.Parse(lay, v)
				if err == nil {
					return dt
				} else {
					return v
				}
			}

			return value
		},
	})

	noggo.Mapping.Type("[date]", Map{
		"name": "日期时间数组", "text": "日期时间数组",
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
			case *[]time.Time:
				return v
			case string:
				lay := "2006-01-02 15:04:05"
				if len(v) == 8 {
					lay = "20060102"
				} else if len(v) == 10 {
					lay = "2006-01-02"
				} else {
					lay = "2006-01-02 15:04:05"
				}

				dt,err := time.Parse(lay, v)
				if err == nil {
					return []time.Time{dt}
				}
			}

			return value
		},
	})


}