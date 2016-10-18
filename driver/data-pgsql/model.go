package data_pgsql


import (
	. "github.com/nogio/noggo/base"
	//"github.com/nogio/noggo/driver"
	"fmt"
	"strings"
	"github.com/nogio/noggo"
	"encoding/json"
	"strconv"
	"errors"
)

type (
	PgsqlModel struct {
		base    *PgsqlBase
		name    string  //模型名称
		schema  string  //架构名
		table   string  //表名
		key     string  //主键
		fields  Map     //字段定义
	}
)



//创建的时候,也需要对值来处理,
//数组要转成{a,b,c}格式,要不然不支持
//json可能要转成字串才支持
func (model *PgsqlModel) packing(value Map) Map {

	newValue := Map{}

	for k,v := range value {
		switch t := v.(type) {
		case []string: {
			newValue[k] = fmt.Sprintf("{%s}", strings.Join(t, ","))
		}
		case []int: {
			arr := []string{}
			for _,v := range t {
				arr = append(arr, strconv.Itoa(v))
			}

			newValue[k] = fmt.Sprintf("{%s}", strings.Join(arr, ","))
		}
		case []int8: {
			arr := []string{}
			for _,v := range t {
				arr = append(arr, fmt.Sprintf("%v", v))
			}

			newValue[k] = fmt.Sprintf("{%s}", strings.Join(arr, ","))
		}
		case []int16: {
			arr := []string{}
			for _,v := range t {
				arr = append(arr, fmt.Sprintf("%v", v))
			}

			newValue[k] = fmt.Sprintf("{%s}", strings.Join(arr, ","))
		}
		case []int32: {
			arr := []string{}
			for _,v := range t {
				arr = append(arr, fmt.Sprintf("%v", v))
			}

			newValue[k] = fmt.Sprintf("{%s}", strings.Join(arr, ","))
		}
		case []int64: {
			arr := []string{}
			for _,v := range t {
				arr = append(arr, fmt.Sprintf("%v", v))
			}

			newValue[k] = fmt.Sprintf("{%s}", strings.Join(arr, ","))
		}
		case Map: {
			b,e := json.Marshal(t);
			if e == nil {
				newValue[k] = string(b)
			} else {
				newValue[k] = "{}"
			}
		}
		case map[string]interface{}: {
			b,e := json.Marshal(t);
			if e == nil {
				newValue[k] = string(b)
			} else {
				newValue[k] = "{}"
			}
		}
		default:
			newValue[k] = t
		}
	}
	return newValue
}




//创建对象
func (model *PgsqlModel) Create(data Map) (error,Map) {

	//按字段生成值
	value := Map{}
	err := noggo.Mapping.Parse([]string{}, model.fields, data, value);
	//Log.Debug("create.map.err=%v, v=%v\n", e, value)
	noggo.Logger.Debug("create", model.fields, data)

	if err != nil {
		return err,nil
	} else {



		//对拿到的值进行包装，以适合pgsql
		newValue := model.packing(value)

		//先拿字段列表
		keys, tags, vals := []string{}, []string{}, make([]interface{},0)
		i := 1
		for k,v := range newValue {
			if k == "id" {
				if v == nil {
					continue
				}
				//id不直接跳过,postgres可以指定ID
				//continue
			}
			keys = append(keys, k)
			vals = append(vals, v)
			tags = append(tags, fmt.Sprintf("$%d", i))
			i++
		}

		err,tx := model.base.Begin()
		if err != nil {
			return err,nil
		} else {

			sql := fmt.Sprintf(`INSERT INTO "%s"."%s" ("%s") VALUES (%s) RETURNING "id";`, model.schema, model.table, strings.Join(keys, `","`), strings.Join(tags, `,`))
			row := tx.QueryRow(sql, vals...)
			if row == nil {
				return errors.New("数据：插入：无返回行"),nil
			} else {

				id := int64(0)
				err := row.Scan(&id)
				if err != nil {
					//扫描新ID失败
					return err,nil
				} else {
					value["id"] = id

					//提交事务
					err := model.base.Submit()
					if err != nil {
						return err,nil
					} else {

						//这里应该有触发器

						//成功了
						return nil,value

					}
				}
			}
		}
	}

	return errors.New("数据：创建失败"),nil
}

//修改对象
func (db *PgsqlModel) Change(item Map, data Map) (error,Map) {
	return nil, nil
}

//删除对象
func (db *PgsqlModel) Remove(Map) (error) {
	return nil
}

//查询对象
func (db *PgsqlModel) Entity(key Any) (error,Map) {
	return nil,Map{}
}