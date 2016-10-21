package types

import (
	"github.com/nogio/noggo"
	. "github.com/nogio/noggo/base"
	"strconv"
	"encoding/json"
	"strings"
)

func init() {


	noggo.Mapping.Type("int", Map{
		"name": "整型", "text": "整型",
		"type": []string{ "int", "number" },
		"valid": func(value Any, config Map) bool {

			switch v := value.(type) {
			case int,int32,int64,int8: {
				return true
			}
			case float32, float64:
				return true
			case string: {
				if _, e := strconv.Atoi(v); e == nil {
					return true
				}
			}
			default:
			}

			return false;
		},
		"value": func(value Any, config Map) Any {
			switch v := value.(type) {
			case int: {
				return int64(v)
			}
			case int8: {
				return int64(v)
			}
			case int16: {
				return int64(v)
			}
			case int32: {
				return int64(v)
			}
			case int64: {
				return int64(v)
			}
			case float32:
				return int64(v)
			case float64:
				return int64(v)
			case string: {
				if i, e := strconv.Atoi(v); e == nil {
					return int64(i)
				}
			}
			default:

			}

			return int64(0);
		},
	})


	noggo.Mapping.Type("[int]", Map{
		"name": "整型数组", "text": "整型数组",
		"type": []string{ "[int]", "[number]" },
		"valid": func(value Any, config Map) bool {

			switch v := value.(type) {
			case int,int8,int16,int32,int64: {
				return true
			}
			case []int,[]int8,[]int16,[]int32,[]int64: {
				return true
			}
			case float32, float64:
				return true
			case []float32, []float64:
				return true
			case []string: {
				if len(v) > 0 {
					for _,sv := range v {
						if _, e := strconv.ParseInt(sv, 10, 64); e != nil {
							return false
						}
					}
					return true
				}
			}
			case string: {

				//此为postgresql数组返回的数组格式
				//{1,2,3,4,5,6,7,8,9}
				if strings.HasPrefix(v, "{") && strings.HasSuffix(v, "}") {
					arr := strings.Split(v[1:len(v)-1], ",")

					for _,sv := range arr {
						if _, e := strconv.ParseInt(sv, 10, 64); e != nil {
							return false
						}
					}
					return true
				} else if strings.HasPrefix(v, "[") && strings.HasSuffix(v, "]") {
					jv := []interface{}{}
					e := json.Unmarshal([]byte(v), &jv)
					if e == nil {
						return true
					} else {
						return false
					}
				} else {

					if _, e := strconv.ParseInt(v, 10, 64); e == nil {
						return true
					} else {
						return false
					}
				}




			}
			default:
			}

			return false;
		},
		"value": func(value Any, config Map) Any {

			switch v := value.(type) {
			case int: {
				return []int64{ int64(v) }
			}
			case int8: {
				return []int64{ int64(v) }
			}
			case int16: {
				return []int64{ int64(v) }
			}
			case int32: {
				return []int64{ int64(v) }
			}
			case int64: {
				return []int64{ int64(v) }
			}


			case []int: {
				arr := []int64{}
				for _,iv := range v {
					arr = append(arr, int64(iv))
				}
				return arr
			}
			case []int8: {
				arr := []int64{}
				for _,iv := range v {
					arr = append(arr, int64(iv))
				}
				return arr
			}
			case []int16: {
				arr := []int64{}
				for _,iv := range v {
					arr = append(arr, int64(iv))
				}
				return arr
			}
			case []int32: {
				arr := []int64{}
				for _,iv := range v {
					arr = append(arr, int64(iv))
				}
				return arr
			}
			case []int64: {
				arr := []int64{}
				for _,iv := range v {
					arr = append(arr, int64(iv))
				}
				return arr
			}



			case float32:
				return []int64{ int64(v) }
			case float64:
				return []int64{ int64(v) }




			case []float32:
				arr := []int64{}
				for _,iv := range v {
					arr = append(arr, int64(iv))
				}
				return arr
			case []float64:
				arr := []int64{}
				for _,iv := range v {
					arr = append(arr, int64(iv))
				}
				return arr




			case []string: {
				arr := []int64{}
				for _,sv := range v {
					if iv, e := strconv.ParseInt(sv, 10, 64); e == nil {
						arr = append(arr, int64(iv))
					}
				}
				return arr
			}
			case string: {


				if strings.HasPrefix(v, "{") && strings.HasSuffix(v, "}") {
					arr := []int64{}
					strArr := strings.Split(v[1:len(v)-1], ",")
					for _,sv := range strArr {
						if iv, e := strconv.ParseInt(sv, 10, 64); e == nil {
							arr = append(arr, iv)
						}
					}
					return arr
				} else if strings.HasPrefix(v, "[") && strings.HasSuffix(v, "]") {
					jv := []interface{}{}
					e := json.Unmarshal([]byte(v), &jv)
					if e == nil {

						arr := []int64{}
						//所以符合的类型,才写入数组
						//json回转,所有的数都是float64
						for _,anyVal := range jv {
							if newVal,ok := anyVal.(float64); ok {
								arr = append(arr, int64(newVal))
							}
						}
						return arr
					}
				} else {

					if vvv, e := strconv.ParseInt(v, 10, 64); e == nil {
						return []int64{ vvv }
					}
				}


			}
			default:
			}


			return []int64{}
		},
	})



}