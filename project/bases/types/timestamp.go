package types

import (
	"github.com/nogio/noggo"
	. "github.com/nogio/noggo/base"
	"time"
)

func init() {

	noggo.Mapping.Type("timestamp", Map{
		"name": "时间戳", "text": "时间戳",
		"type": []string{ "timestamp" },
		"valid": func(value Any, config Map) bool {

			switch v := value.(type) {
			case time.Time:
				return true
			case string:
				return noggo.Const.Valid(v, "datetime")
			}

			return false;
		},
		"value": func(value Any, config Map) Any {
			switch v := value.(type) {
			case time.Time:
				return v.Unix()
			case string:
				dt,err := time.Parse("2006-01-02 15:04:05", v)
				if err == nil {
					return dt.Unix()
				} else {
					return v
				}
			}

			return value;
		},
	})

	noggo.Mapping.Type("[timestamp]", Map{
		"name": "时间戳数组", "text": "时间戳数组",
		"type": []string{ "[timestamp]" },
		"valid": func(value Any, config Map) bool {

			switch v := value.(type) {
			case time.Time:
				return true
			case []time.Time:
				return true
			case string:
				return noggo.Const.Valid(v, "datetime")
			}

			return false;
		},
		"value": func(value Any, config Map) Any {
			switch v := value.(type) {
			case time.Time:
				return []int64{ v.Unix() }
			case []time.Time: {
				ts := []int64{}
				for _,dt := range v {
					ts = append(ts, dt.Unix())
				}
				return ts
			}
			case string:
				//应该JSON解析
				dt,err := time.Parse("2006-01-02 15:04:05", v)
				if err == nil {
					return dt.Unix()
				} else {
					return v
				}
			}

			return value;
		},
	})



}