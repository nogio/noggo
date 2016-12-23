package types

import (
	"github.com/nogio/noggo"
	. "github.com/nogio/noggo/base"
)

func init() {



	noggo.Mapping.Type("file", Map{
		"name": "file", "text": "file",
		"valid": func(value Any, config Map) bool {

			switch value.(type) {
			case Map:
				return true
			case map[string]interface{}:
				return true
			}

			return false;
		},
		"value": func(value Any, config Map) Any {


			switch vv := value.(type) {
			case Map:
				return vv
			case map[string]interface{}:
				m := Map{}
				for k, v := range vv {
					m[k] = v
				}
				return m
			}
			return Map{}
		},
	})


	noggo.Mapping.Type("[file]", Map{
		"name": "文件数组", "text": "文件数组",
		"valid": func(value Any, config Map) bool {

			switch value.(type) {
			case Map:
				return true
			case map[string]interface{}:
				return true
			case []Map:
				return true
			case []map[string]interface{}:
				return true
			}

			return false;
		},
		"value": func(value Any, config Map) Any {

			switch v := value.(type) {
			case Map:
				return []Map{ v }
			case map[string]interface{}:
				m := Map{}
				for kk,vv := range v {
					m[kk] = vv
				}
				return []Map{ m }

			case []Map:
				return v
			case []map[string]interface{}:
				ms := []Map{}

				for _,i := range v {
					m := Map{}
					for kk,vv := range i {
						m[kk] = vv
					}
					ms = append(ms, m)
				}


				return ms
			}
			return []Map{}
		},
	})



}