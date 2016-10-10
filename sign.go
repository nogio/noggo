/*
	sign 签名功能
	主要是用来请求中签名用的
	一般用在登录验证
*/

package noggo

import (
	. "github.com/nogio/noggo/base"
	"fmt"
)


type (
	Sign struct {
		Value Map
	}
)


//签入
func (sign *Sign) In(key string, id,name Any) (bool) {
	if sign.Value != nil {
		m := Map{
			KeySignId: fmt.Sprintf("%v", id), KeySignName: fmt.Sprintf("%v", name),
		}
		sign.Value[KeySign+key] = m

		return true
	}

	return false
}

//签出
func (sign *Sign) Out(key string) (bool) {
	if sign.Value != nil {
		delete(sign.Value, KeySign + key)
		return true
	}
	return false
}



//判断是否已经签入了
func (sign *Sign) Yes(key string) bool {
	if sign.Value != nil{
		if sign.Value[KeySign+key] != nil {
			return true
		}
	}

	return false
}



//判断是否已经签出了
func (sign *Sign) No(key string) (bool) {
	return sign.Yes(key) == false
}


//得到签名信息
func (sign *Sign) Info(key string) (Map) {
	if sign.Value != nil {
		if sign.Value[KeySign+key] != nil {

			mmm := Map{}

			switch v := sign.Value[KeySign+key].(type) {
			case Map: {
				for k,v := range v {
					mmm[k] = v
				}
			}
			case map[string]interface{}: {
				for k,v := range v {
					mmm[k] = v
				}
			}
			}

			return mmm
		}
	}
	return nil
}


//返回ID
func (sign *Sign) Id(key string) string {
	if sign.Value != nil {
		if m := sign.Info(key); m != nil {
			if v, ok := m[KeySignId].(string); ok {
				return v
			}
		}
	}
	return ""
}


//返回名称
func (sign *Sign) Name(key string) string {
	if sign.Value != nil {
		if m := sign.Info(key); m != nil {
			if v, ok := m[KeySignName].(string); ok {
				return v
			}
		}
	}
	return ""
}

