package noggo

import (
	. "github.com/nogio/noggo/base"
	"fmt"
	"encoding/json"
	"strings"
	"time"
)

type (
	mappingModule struct {
		types	map[string]Map
		cryptos	map[string]Map
	}
)




//注册类型
func (module *mappingModule) Type(name string, config Map) {
	if module.types == nil {
		module.types = map[string]Map{}
	}
	module.types[name] = config
}

//注册加解密
func (module *mappingModule) Crypto(name string, config Map) {
	if module.cryptos == nil {
		module.cryptos = map[string]Map{}
	}
	module.cryptos[name] = config
}










//类型默认验证
//默认使用string做正则表达式验证
func (module *mappingModule) typeDefaultValid(value Any, config Map) bool {
	if _,ok := config[KeyMapType]; ok {
		//待修改
		//return Check(t.(string), fmt.Sprintf("%s", value))
	}
	return true
}
//类型默认值包
//默认值包装都是返回string类型
func (module *mappingModule) typeDefaultValue(value Any, config Map) Any {
	return fmt.Sprintf("%s", value)
}







//获取类型的验证方法
func (module *mappingModule) TypeValid(name string) (valid func(Any, Map) bool) {
	//配置中的验证方法
	if config, ok := module.types[name]; ok {
		switch method := config[KeyMapValid].(type) {
		case func(Any, Map) bool:
			valid = method
		}
	}

	//没有配置，使用默认验证方法
	if valid == nil {
		valid = module.typeDefaultValid
	}

	return
}
//获取类型的值包方法
func (module *mappingModule) TypeValue(name string) (value func(Any,Map) Any) {
	//配置中的值包方法
	if config, ok := module.types[name]; ok {
		switch method := config[KeyMapValue].(type) {
		case func(Any,Map) Any:
			value = method
		}
	}
	//没有配置，使用默认值包方法
	if value == nil {
		value = module.typeDefaultValue
	}

	return
}


//获取类型方法
func (module *mappingModule) TypeMethod(name string) (func(Any, Map) bool, func(Any,Map) Any) {
	return module.TypeValid(name), module.TypeValue(name)
}
































//默认加密方法
func (module *mappingModule) cryptoDefaultEncode(value Any) Any {
	return value
}
//默认解密方法
func (module *mappingModule) cryptoDefaultDecode(value Any) Any {
	return value
}



//获取加密方法
func (module *mappingModule) CryptoEncode(name string) (encode func(Any) Any) {
	//配置中的加密方法
	if config, ok := module.cryptos[name]; ok {
		switch method := config[KeyMapEncode].(type) {
		case func(Any) Any:
			encode = method
		}
	}
	//没有使用默认加密方法
	if encode == nil {
		encode = module.cryptoDefaultEncode
	}

	return
}

//获取解密方法
func (module *mappingModule) CryptoDecode(name string) (decode func(Any) Any) {
	//默认中的解密方法
	if config, ok := module.cryptos[name]; ok {
		switch method := config[KeyMapDecode].(type) {
		case func(Any) Any:
			decode = method
		}
	}
	//没有使用默认解密方法
	if decode == nil {
		decode = module.cryptoDefaultDecode
	}

	return
}

//获了方法
func (module *mappingModule) CryptoMethod(name string)(func(Any) Any, func(Any) Any) {
	return module.CryptoEncode(name), module.CryptoDecode(name)
}












//解析
func (module *mappingModule) Parse(tree []string, config Map, data Map, value Map, args ...bool) *Error {

	argn := false
	if len(args) > 0 {
		argn = args[0]
	}

	//遍历配置	begin
	for fieldName,fv := range config {
		fieldConfig := Map{}

		//注意，这里存在2种情况
		//1. Map对象
		//2. map[string]interface{}
		//要分开处理

		switch c := fv.(type) {
		case Map:
			fieldConfig = c
		case map[string]interface{}:
			for k,v := range c {
				fieldConfig[k] = v
			}
			//fieldConfig = c
		default:
			//类型不对，跳过
			continue
		}




		//Map 如果是JSON文件，或是发过来的消息，就不支持
		//而用下面的，就算是MAP也可以支持，OK
		//麻烦来了，web.args用下面这样处理不了
		//if fieldConfig, ok := fv.(map[string]interface{}); ok {

		fieldMust, fieldEmpty := false, false
		fieldValue, fieldExist := data[fieldName]
		fieldAuto, fieldJson := fieldConfig[KeyMapAuto], fieldConfig[KeyMapJson]
		//_, fieldEmpty = data[fieldName]

		/* 处理是否必填和为空 */
		if v, ok := fieldConfig[KeyMapMust]; ok {
			if v == nil {
				fieldEmpty = true
			}
			if vv,ok := v.(bool); ok {
				fieldMust = vv
			}
		}


		trees := append(tree, fieldName)
		//fmt.Printf("t=%s, argn=%v, value=%v\n", strings.Join(trees, "."), argn, fieldValue)
		//fmt.Printf("trees=%v, must=%v, empty=%v, exist=%v, auto=%v, value=%v, config=%v\n\n",
		//	strings.Join(trees, "."), fieldMust, fieldEmpty, fieldExist, fieldAuto, fieldValue, fieldConfig)



		strVal := fmt.Sprintf("%v", fieldValue)

		//等一下。 空的map[]无字段。 需要也表示为空吗?
		//if strVal == "" || strVal == "map[]" || strVal == "{}"{

		//因为go1.6之后。把一个值为nil的map  再写入map之后, 判断 if map[key]==nil 就无效了
		if strVal == "" || data[fieldName] == nil || (fieldMust != true && strVal == "map[]")  {
			fieldValue = nil
		}



		//如果不可为空，但是为空了，
		if (fieldMust && fieldEmpty == false && (fieldValue == nil || strVal == "") && fieldAuto == nil && fieldJson == nil && argn == false) {

			//是否有自定义的状态
			if c,ok := fieldConfig["empty"]; ok {
				if fieldConfig["name"] != nil {
					return NewStateError(c.(string), fmt.Sprintf("%v", fieldConfig["name"]))
				} else {
					return NewStateError(c.(string), strings.Join(trees, "."))
				}

			} else {
				//return errors.New("参数不可为空")
				//return NewStateError("args.empty", fieldName)
				if fieldConfig["name"] != nil {
					return NewStateError("map.empty", fmt.Sprintf("%v", fieldConfig["name"]))
				} else {
					return NewStateError("map.empty", strings.Join(trees, "."))
				}
			}

		} else {

			//如果值为空的时候，看有没有默认值
			if (fieldValue == nil || strVal == "") {

				//如果有默认值
				//可为NULL时，不给默认值
				if (fieldAuto != nil && !argn) {

					//暂时不处理 $id, $date 之类的定义好的默认值，不包装了
					switch autoValue:=fieldAuto.(type) {
					case func() interface{}:
						fieldValue = autoValue()
					case func() time.Time:
						fieldValue = autoValue()
					//case func() bson.ObjectId:	//待修改
					//fieldValue = autoValue()
					case func() string:
						fieldValue = autoValue()
					case func() int:
						fieldValue = int64(autoValue())
					case func() int8:
						fieldValue = int64(autoValue())
					case func() int16:
						fieldValue = int64(autoValue())
					case func() int32:
						fieldValue = int64(autoValue())
					case func() int64:
						fieldValue = autoValue()
					case func() float32:
						fieldValue = float64(autoValue())
					case func() float64:
						fieldValue = autoValue()

					case int: {
						fieldValue = int64(autoValue)
					}
					case int8: {
						fieldValue = int64(autoValue)
					}
					case int16: {
						fieldValue = int64(autoValue)
					}
					case int32: {
						fieldValue = int64(autoValue)
					}
					case float32: {
						fieldValue = float64(autoValue)
					}


					case []int: {
						arr := []int64{}
						for _,nv := range autoValue {
							arr = append(arr, int64(nv))
						}
						fieldValue = arr
					}
					case []int8: {
						arr := []int64{}
						for _,nv := range autoValue {
							arr = append(arr, int64(nv))
						}
						fieldValue = arr
					}
					case []int16: {
						arr := []int64{}
						for _,nv := range autoValue {
							arr = append(arr, int64(nv))
						}
						fieldValue = arr
					}
					case []int32: {
						arr := []int64{}
						for _,nv := range autoValue {
							arr = append(arr, int64(nv))
						}
						fieldValue = arr
					}

					case []float32: {
						arr := []float64{}
						for _,nv := range autoValue {
							arr = append(arr, float64(nv))
						}
						fieldValue = arr
					}

					default:
						fieldValue = autoValue
					}

				} else {	//没有默认值, 且值为空


					//有个问题, POST表单的时候.  表单的值是存在的
					//但是POST的时候如果有argn, 实际上是不想存在此字段的

					//如果字段可以不存在
					if (fieldEmpty || argn) {

						//当empty(argn)=true，并且如果传过值中存在字段的KEY，值就要存在，以更新为null
						if (argn && fieldExist) {
							//不操作，自然往下执行
						} else {	//值可以不存在
							continue
						}

					} else if (argn) {


					} else {
						//到这里不管
						//因为字段必须存在，但是值为空
					}



				}

			} else {	//值不为空，处理值


				//值处理前，是不是需要解密
				//如果解密哦
				if ct,ok := fieldConfig["decode"]; ok {
					//而且要值是string类型
					if sv,ok := fieldValue.(string); ok {

						//得到解密方法
						decode := CryptoDecode(ct.(string))
						fieldValue = decode(sv)
					}
				}



				//按类型来做处理

				//验证方法和值方法
				if fieldType, ok := fieldConfig["type"]; ok {
					fieldValidCall, fieldValueCall := TypeFunc(fieldType.(string))

					//如果配置中有自己的验证函数
					if f,ok := fieldConfig["valid"]; ok {
						if call,ok := f.(func(Any,Map)bool); ok {
							fieldValidCall = call
						}
					}
					//如果配置中有自己的值函数
					if f,ok := fieldConfig["value"]; ok {
						if call, ok := f.(func(Any,Map)Any); ok {
							fieldValueCall = call
						}
					}


					//开始调用验证
					if fieldValidCall != nil {
						//如果验证通过
						if (fieldValidCall(fieldValue, fieldConfig)) {
							//包装值
							if fieldValueCall != nil {
								fieldValue = fieldValueCall(fieldValue, fieldConfig)
							}
						} else {	//验证不通过

							//是否有自定义的状态
							if c,ok := fieldConfig["error"]; ok {
								return NewStateError(c.(string), strings.Join(trees, "."))

								if fieldConfig["name"] != nil {
									return NewStateError(c.(string), fmt.Sprintf("%v", fieldConfig["name"]))
								} else {
									return NewStateError(c.(string), strings.Join(trees, "."))
								}


							} else {
								//return errors.New("valid error")
								//类型不匹配
								//return NewStateError("args.error", fieldName)
								if fieldConfig["name"] != nil {
									return NewStateError("map.error", fmt.Sprintf("%v", fieldConfig["name"]))
								} else {
									return NewStateError("map.error", strings.Join(trees, "."))
								}
							}
						}
					}
				}



			}

		}

		//这后面是总的字段处理
		//如JSON，加密

		//如果是JSON， 或是数组啥的处理
		//注意，当 json 本身可为空，下级有不可为空的，值为空时， 应该跳过子级的检查
		//如果 json 可为空， 就不应该有 默认值， 定义的时候要注意啊啊啊啊
		//理论上，只要JSON可为空～就不处理下一级json
		jsonning := true
		if !fieldMust && fieldValue == nil {
			jsonning = false;
		}

		//还有种情况要处理. 当type=json, must=true的时候,有默认值, 但是没有定义json节点.

		if m,ok := fieldConfig["json"]; ok && jsonning {
			jsonConfig := Map{}
			//注意，这里存在2种情况
			//1. Map对象 //2. map[string]interface{}
			switch c := m.(type) {
			case Map:
				jsonConfig = c
			case map[string]interface{}:
				jsonConfig = c
			}


			//如果是数组
			isArray := false
			//fieldData到这里定义
			var fieldData Any

			switch v := fieldValue.(type) {
			case Map:
				fieldData = v
				fieldValue = Map{}
			case map[string]interface{}: {
				//这里要处理, 因为当json字段有多级的时候, 解析出来是 map[string]interface{}  这样处理子级的时候转成了Map就出问题了
				mm := Map{}
				for kk,vv := range v {
					mm[kk] = vv
				}

				fieldData = mm
				fieldValue = Map{}
			}
			case []Map:
				isArray = true
				fieldData = v
				fieldValue = make([]Map,0)
			case []map[string]interface{}:
				isArray = true
				fieldData = v
				fieldValue = make([]map[string]interface{},0)
			default:
				fieldData = Map{}
				fieldValue = Map{}
			}

			//是数组，处理JSON数组
			if isArray {

				//如果是Map数组
				if ds,ok := fieldData.([]Map); ok {
					values := fieldValue.([]Map)

					for _,d := range ds {
						v := Map{}

						err := Mapping(trees, jsonConfig, d, v, args...);
						if err != nil {
							return err
						} else {
							//fieldValue = append(fieldValue, v)
							values = append(values, v)
						}
					}
					fieldValue = values

				}
				//如果是map[string]数组
				if ds,ok := fieldData.([]map[string]interface{}); ok {
					values := fieldValue.([]map[string]interface{})

					for _,d := range ds {
						v := make(map[string]interface{})

						err := Mapping(trees, jsonConfig, d, v, args...);
						if err != nil {
							return err
						} else {
							//fieldValue = append(fieldValue, v)
							values = append(values, v)
						}
					}

					fieldValue = values

				}



			} else {
				//不是数组，直接是JSON，处理下级

				//data,value有点绕
				//data是传过来的参数，  value是要返回来的值
				//在这里 data 应该是处理好的fieldValue 这个字段
				//而，要返回一个新的value做为， 这个字段的value

				//当json有多层的时候, 这里就不是map了
				fieldValueMap := Map{}
				switch c := fieldValue.(type) {
				case Map:
					fieldValueMap = c
				case map[string]interface{}:
					fieldValueMap = c
				default:
				//类型不对
				}




				//当json有多层的时候, 这里就不是map了
				err := Mapping(trees, jsonConfig, fieldData.(Map), fieldValueMap, args...);
				if err != nil {
					return err
				}
			}
		}



		//最后，值要不要加密什么的
		//如果加密
		if ct,ok := fieldConfig["encode"]; ok {

			//全都转成字串再加密
			sv := ""
			switch v:=fieldValue.(type) {
			case string:
				sv = v
			case Map,map[string]interface{}:
				d,e := json.Marshal(v)
				if e == nil {
					sv = string(d)
				} else {
					sv = "{}"
				}
			case []Map,[]map[string]interface{}:
				d,e := json.Marshal(v)
				if e == nil {
					sv = string(d)
				} else {
					sv = "[]"
				}
			default:
				sv = fmt.Sprintf("%v", v)
			}


			//得到解密方法
			decode := CryptoEncode(ct.(string))
			fieldValue = decode(sv)
		}




		//没有JSON要处理，所以给值
		value[fieldName] = fieldValue


	}
	return nil
	//遍历配置	end
}
