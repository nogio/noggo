package native

import (
	"github.com/nogio/noggo"
	. "github.com/nogio/noggo/base"
	"fmt"
	"time"
	"strings"
	"encoding/json"
	"strconv"
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







	noggo.Mapping.Type("date", Map{
		"name": "日期时间", "text": "日期时间",
		"type": []string{ "date" },
		"valid": func(value Any, config Map) bool {

			switch v := value.(type) {
			case time.Time:
				return true
			case *time.Time:
				return true
			case string:
				return noggo.Const.Valid(v, "date")
			}

			return false;
		},
		"value": func(value Any, config Map) Any {
			return value
		},
	})

	noggo.Mapping.Type("[date]", Map{
		"name": "日期时间数组", "text": "日期时间数组",
		"type": []string{ "[date]" },
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
			default:
				return v
			}
		},
	})

	noggo.Mapping.Type("datetime", Map{
		"name": "日期时间", "text": "日期时间",
		"type": []string{ "datetime" },
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
			return value
		},
	})

	noggo.Mapping.Type("[datetime]", Map{
		"name": "日期时间数组", "text": "日期时间数组",
		"type": []string{ "[datetime]" },
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



	noggo.Mapping.Type("enum", Map{
		"name": "枚举", "text": "枚举",
		"type": []string{ "enum" },
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
		"name": "枚举", "text": "枚举",
		"type": []string{ "[enum]" },
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


	noggo.Mapping.Type("float", Map{
		"name": "浮点型", "text": "布尔型",
		"type": []string{ "float", "double" },
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










	noggo.Mapping.Type("json", Map{
		"name": "JSON", "text": "JSON",
		"type": []string{ "json" },
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


	noggo.Mapping.Type("json.array", Map{
		"name": "JSON数组", "text": "JSON数组",
		"type": []string{ "[json]" },
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


	noggo.Mapping.Type("int.array", Map{
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


	noggo.Mapping.Type("string.array", Map{
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

	noggo.Mapping.Type("timestamp.array", Map{
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

