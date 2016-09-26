package noggo

import (
	. "github.com/nogio/noggo/base"
)


type (
	Sign struct {
		Value Map
	}
)


//设置登录信息
func (sign *Sign) In(key string, id,name Any, args ...Map) (bool) {
	if sign.Value != nil {
		data := Map{}
		if len(args) > 0 {
			data = args[0]
		}
		m := Map{
			KeySignId: id, KeySignName: name, KeySignData: data,
		}
		sign.Value[KeySign+key] = m

		return true
	}

	return false
}

//设置为登出
func (sign *Sign) Out(key string) (bool) {
	if sign.Value != nil {
		delete(sign.Value, KeySign + key)
		return true
	}
	return false
}



//判断是否已经登入了
func (sign *Sign) Yes(key string) bool {
	if sign.Value != nil{
		if sign.Value[KeySign+key] != nil {
			return true
		}
	}

	return false
}



//判断是否已经登出了
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
func (sign *Sign) Id(key string) Any {
	if sign.Value != nil {
		if m := sign.Info(key); m != nil {
			if v, ok := m[KeySignId]; ok {
				return v
			}
		}
	}
	return ""
}


//返回名称
func (sign *Sign) Name(key string) Any {
	if sign.Value != nil {
		if m := sign.Info(key); m != nil {
			if v, ok := m[KeySignName]; ok {
				return v
			}
		}
	}
	return ""
}



//返回数据
func (sign *Sign) Data(key string) Map {
	if sign.Value != nil {
		if m := sign.Info(key); m != nil {
			if v, ok := m[KeySignData].(Map); ok {
				return v
			}
		}
	}
	return nil
}




