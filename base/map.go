package base

import "fmt"

type (
	Map         map[string]Any
)


//返回字符串
func (m Map) String(key string) string {
	return fmt.Sprintf("%v", m[key])
}

//返回数字
func (m Map) Number(key string) float64 {
	if m[key] != nil {

		switch vv:=m[key].(type) {
		case int:
			return float64(vv)
		case int8:
			return float64(vv)
		case int16:
			return float64(vv)
		case int32:
			return float64(vv)
		case float32:
			return float64(vv)
		case float64:
			return float64(vv)
		}
		
	}

	return float64(0)
}
