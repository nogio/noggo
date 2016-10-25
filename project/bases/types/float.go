package types

import (
	"github.com/nogio/noggo"
	. "github.com/nogio/noggo/base"
	"strconv"
	"encoding/json"
	"strings"
)

func init() {


	noggo.Mapping.Type("float", Map{
		"name": "浮点型", "text": "布尔型",
		"valid": func(value Any, config Map) bool {

			switch v := value.(type) {
			case int,int32,int64,int8: {
				return true
			}
			case float32,float64: {
				return true
			}
			case string: {
				if _,e := strconv.ParseFloat(v, 64); e == nil {
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
				return float64(v)
			}
			case int32: {
				return float64(v)
			}
			case int64: {
				return float64(v)
			}
			case int8: {
				return float64(v)
			}
			case float32: {
				return float64(v)
			}
			case float64: {
				return v
			}
			case string: {
				if v,e := strconv.ParseFloat(v, 64); e == nil {
					return v
				}
			}
			default:
			}

			return float64(0.0);
		},
	})



	noggo.Mapping.Type("[float]", Map{
		"name": "浮点数组", "text": "浮点数组",
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
						if _, e := strconv.ParseFloat(sv, 64); e != nil {
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

					if _, e := strconv.ParseFloat(v, 64); e == nil {
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
				return []float64{ float64(v) }
			}
			case int8: {
				return []float64{ float64(v) }
			}
			case int16: {
				return []float64{ float64(v) }
			}
			case int32: {
				return []float64{ float64(v) }
			}
			case int64: {
				return []float64{ float64(v) }
			}


			case []int: {
				arr := []float64{}
				for _,iv := range v {
					arr = append(arr, float64(iv))
				}
				return arr
			}
			case []int8: {
				arr := []float64{}
				for _,iv := range v {
					arr = append(arr, float64(iv))
				}
				return arr
			}
			case []int16: {
				arr := []float64{}
				for _,iv := range v {
					arr = append(arr, float64(iv))
				}
				return arr
			}
			case []int32: {
				arr := []float64{}
				for _,iv := range v {
					arr = append(arr, float64(iv))
				}
				return arr
			}
			case []int64: {
				arr := []float64{}
				for _,iv := range v {
					arr = append(arr, float64(iv))
				}
				return arr
			}



			case float32:
				return []float64{ float64(v) }
			case float64:
				return []float64{ float64(v) }




			case []float32:
				arr := []float64{}
				for _,iv := range v {
					arr = append(arr, float64(iv))
				}
				return arr
			case []float64:
				return v




			case []string: {
				arr := []float64{}
				for _,sv := range v {
					if iv, e := strconv.ParseFloat(sv, 64); e == nil {
						arr = append(arr, float64(iv))
					}
				}
				return arr
			}
			case string: {


				if strings.HasPrefix(v, "{") && strings.HasSuffix(v, "}") {
					arr := []float64{}
					strArr := strings.Split(v[1:len(v)-1], ",")
					for _,sv := range strArr {
						if iv, e := strconv.ParseFloat(sv, 64); e == nil {
							arr = append(arr, iv)
						}
					}
					return arr
				} else if strings.HasPrefix(v, "[") && strings.HasSuffix(v, "]") {
					jv := []interface{}{}
					e := json.Unmarshal([]byte(v), &jv)
					if e == nil {

						arr := []float64{}
						//所以符合的类型,才写入数组
						//json回转,所有的数都是float64
						for _,anyVal := range jv {
							if newVal,ok := anyVal.(float64); ok {
								arr = append(arr, newVal)
							}
						}
						return arr
					}
				} else {

					if vvv, e := strconv.ParseFloat(v, 64); e == nil {
						return []float64{ vvv }
					}
				}


			}
			default:
			}


			return []float64{}
		},
	})





}