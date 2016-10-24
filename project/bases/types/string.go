package types

import (
	"github.com/nogio/noggo"
	. "github.com/nogio/noggo/base"
	"fmt"
	"strings"
)

func init() {




	noggo.Mapping.Type("string", Map{
		"name": "字符串", "text": "字符串",
		"type": []string{ "string" },
		"valid": func(value Any, config Map) bool {
			switch v := value.(type) {
			case string:
				if v != "" {
					return true
				}
			case []byte:
				s := fmt.Sprintf("%s", v)
				if s != "" {
					return true
				}
			default:
				if value != nil {
					return true
				}
			}
			return false
		},
		"value": func(value Any, config Map) Any {
			return fmt.Sprintf("%s", value)
		},
	})


	noggo.Mapping.Type("[string]", Map{
		"name": "字符数组", "text": "字符数组",
		"type": []string{ "[string]" },
		"valid": func(value Any, config Map) bool {
			switch v := value.(type) {
			case []string:
				//要不要判断是否为空数组
				return true
			case string:
				left, right := v[0:1], v[len(v)-1:len(v)]
				if left == "[" && right == "]" {
					return true
				} else if left == "{" && right == "}" {
					return true
				} else if strings.Index(v, "\n") != -1 {
					return true
				} else {
					return false
				}
			default:
				return false
			}
		},
		"value": func(value Any, config Map) Any {

			switch v := value.(type) {
			case []string:
				return v
			case string:


				left, right := v[0:1], v[len(v)-1:len(v)]
				if left == "[" && right == "]" {

					s := v[1:len(v)-1]	//去掉[] {}
					return strings.Split(s, ",")

				} else if left == "{" && right == "}" {

					s := v[1:len(v)-1]	//去掉[] {}
					return strings.Split(s, ",")

				} else if strings.Index(v, "\n") != -1 {

					s := strings.Replace(v, "\r", "", -1)
					return strings.Split(s, "\n")

				} else {
					return []string{}
				}


			/*
			s := v[1:len(v)-1]	//去掉[] {}
			if s == "" {
				return []string{}
			} else {
				return strings.Split(s, ",")
			}
			*/
			default:
				return v
			}
		},
	})





}
