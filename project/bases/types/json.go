package types

import (
	"github.com/nogio/noggo"
	. "github.com/nogio/noggo/base"
	"encoding/json"
)

func init() {





	noggo.Mapping.Type("json", Map{
		"name": "JSON", "text": "JSON",
		"valid": func(value Any, config Map) bool {

			switch v := value.(type) {
			case Map:
				return true
			case map[string]interface{}:
				return true
			case string:

				m := Map{}
				err := json.Unmarshal([]byte(v), &m)
				if err == nil {
					return true
				}
			}

			return false;
		},
		"value": func(value Any, config Map) Any {


			switch vv := value.(type) {
			case Map:
				return vv
			case map[string]interface{}:
				m := Map{}
				for k,v := range vv {
					m[k] = v
				}
				return m
			case string:
				m := Map{}
				err := json.Unmarshal([]byte(vv), &m)
				if err == nil {
					return m
				}
			}
			return Map{}
		},
	})


	noggo.Mapping.Type("[json]", Map{
		"name": "JSON数组", "text": "JSON数组",
		"valid": func(value Any, config Map) bool {

			switch v := value.(type) {
			case []Map:
				return true
			case []map[string]interface{}:
				return true
			case string:
				m := []Map{}
				err := json.Unmarshal([]byte(v), &m)
				if err == nil {
					return true
				}
			}

			return false;
		},
		"value": func(value Any, config Map) Any {

			switch v := value.(type) {
			case []Map:
				return v
			case []map[string]interface{}:
				return v
			case string:
				m := []Map{}
				err := json.Unmarshal([]byte(v), &m)
				if err == nil {
					return m
				}
			}
			return []Map{}
		},
	})



}