package types

import (
	"github.com/nogio/noggo"
	. "github.com/nogio/noggo/base"
	"encoding/json"
	"strings"
	"fmt"
)

func init() {



	noggo.Mapping.Type("enum", Map{
		"name": "枚举", "text": "枚举",
		"valid": func(value Any, config Map) bool {

			sv := fmt.Sprintf("%v", value)

			if e,ok := config["enum"]; ok {
				for k,_ := range e.(Map) {
					if k == sv {
						return true
					}
				}
			}
			return false;
		},
		"value": func(value Any, config Map) Any {
			return fmt.Sprintf("%v", value);
		},
	})


	noggo.Mapping.Type("[enum]", Map{
		"name": "枚举数组", "text": "枚举数组",
		"valid": func(value Any, config Map) bool {

			vals := []string{}

			switch v := value.(type) {
			case string: {

				//如果是 {},  []  包括的字串，就做split
				//postgres中的， {a,b,c} 格式
				if v[0:1] == "{" && v[len(v)-1:len(v)] == "}" {
					v = v[1:len(v) - 1]
					vals = strings.Split(v, ",")
				} else if v[0:1] == "[" && v[len(v)-1:len(v)] == "]" {
					//json数组格式
					json.Unmarshal([]byte(v), &vals)
				} else {
					vals = append(vals, v)
				}
			}
			case []string: {
				vals = v
			}
			default:
				vals = append(vals, fmt.Sprintf("%v", v))
			}


			oks := 0

			//遍历数组， 全部在enum里才行
			if e,ok := config["enum"]; ok {
				for k,_ := range e.(Map) {

					for _,v := range vals {
						if k == v {
							oks++
						}
					}
				}
			}


			if oks >= len(vals) {
				return true
			} else {
				return false
			}

		},
		"value": func(value Any, config Map) Any {
			vals := []string{}

			switch v := value.(type) {
			case string: {

				//如果是 {},  []  包括的字串，就做split
				//postgres中的， {a,b,c} 格式
				if v[0:1] == "{" && v[len(v)-1:len(v)] == "}" {
					v = v[1:len(v) - 1]
					vals = strings.Split(v, ",")
				} else if v[0:1] == "[" && v[len(v)-1:len(v)] == "]" {
					//json数组格式
					json.Unmarshal([]byte(v), &vals)
				} else {
					vals = append(vals, v)
				}
			}
			case []string: {
				vals = v
			}
			default:
				vals = append(vals, fmt.Sprintf("%v", v))
			}
			return vals
		},
	})


}