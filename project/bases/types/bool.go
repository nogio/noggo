package types

import (
	"github.com/nogio/noggo"
	. "github.com/nogio/noggo/base"
	"fmt"
)

func init() {

	noggo.Mapping.Type("bool", Map{
		"name": "布尔型", "text": "布尔型",
		"valid": func(value Any, config Map) bool {

			switch v := value.(type) {
			case bool: {
				return true
			}
			case string: {
				if v=="true" || v=="false" || v=="0" || v=="1" || v=="yes" || v=="no" {
					return true
				}
			}
			case int,int8,int16,int32,int64,float32,float64:{
				return true
			}
			default:
			}

			return false;
		},
		"value": func(value Any, config Map) Any {

			switch v := value.(type) {
			case bool: {
				return v
			}
			case string: {
				if v=="true" || v=="1" || v=="yes" {
					return true
				} else {
					return false
				}
			}
			case int,int8,int16,int32,int64,float32,float64:{
				s := fmt.Sprintf("%v", v)
				if s == "0" {
					return false
				} else {
					return true
				}
			}
			default:

			}

			return false;
		},
	})


	noggo.Mapping.Type("[bool]", Map{
		"name": "布尔型数组", "text": "布尔型数组",
		"valid": func(value Any, config Map) bool {

			switch v := value.(type) {
			case bool: {
				return true
			}
			case []bool: {
				return true
			}
			case string: {
				if v=="true" || v=="false" || v=="0" || v=="1" || v=="yes" || v=="no" {
					return true
				}
			}
			case []string: {
				for _,s := range v {
					if !(s=="true" || s=="false" || s=="0" || s=="1" || s=="yes" || s=="no") {
						return false
					}
				}
				return true
			}
			default:

			}

			return false;
		},
		"value": func(value Any, config Map) Any {

			switch v := value.(type) {
			case bool: {
				return []bool { true }
			}
			case []bool: {
				return v
			}
			case string: {
				if v=="true" || v=="1" || v=="yes" {
					return []bool { true }
				} else {
					return []bool { false }
				}
			}
			case []string: {
				vvvvv := []bool { }
				for _,s := range v {
					if s=="true" || s=="1" || s=="yes" {
						vvvvv = append(vvvvv, true)
					} else {
						vvvvv = append(vvvvv, false)
					}
				}
				return vvvvv
			}
			default:

			}

			return false;
		},
	})



}