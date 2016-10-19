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
func (model *PgsqlModel) Change(item Map, data Map) (error,Map) {


	//按字段生成值
	value := Map{}
	err := noggo.Mapping.Parse([]string{}, model.fields, data, value, true);
	noggo.Logger.Debug("change", "mapping", err)

	if err != nil {
		return err,nil
	} else {

		//包装值，因为golang本身数据类型和数据库的不一定对版
		//需要预处理一下
		newValue := model.packing(value)

		//先拿字段列表
		sets, vals := []string{}, make([]interface{}, 0)
		i := 1
		for k, v := range newValue {
			//主值不在修改之中
			if k == model.key {
				continue
			}
			//keys = append(keys, k)
			vals = append(vals, v)
			sets = append(sets, fmt.Sprintf(`"%s"=$%d`, k, i))
			i++
		}
		//条件是主键
		vals = append(vals, item[model.key])

		//开启事务
		err, tx := model.base.Begin()
		if err != nil {
			return err,nil
		} else {

			//更新数据库
			sql := fmt.Sprintf(`UPDATE "%s"."%s" SET %s WHERE "id"=$%d`, model.schema, model.table, strings.Join(sets, `,`), i)
			_, err := tx.Exec(sql, vals...)
			noggo.Logger.Debug("change", "exec", err)
			if err != nil {
				return err,nil
			} else {

				// 不改item
				// 先复制item
				newItem := Map{}
				for k, v := range item { newItem[k] = v }
				for k, v := range value { newItem[k] = v }


				//提交事务
				err := model.base.Submit()
				if err != nil {
					return err,nil
				} else {

					//这里应该有触发器

					//成功了
					return nil,newItem

				}


			}
		}
	}

	return errors.New("数据：修改失败"), nil
}

//删除对象
func (model *PgsqlModel) Remove(item Map) (error) {

	if key,ok := item[model.key]; ok {

		//开启事务
		err, tx := model.base.Begin()
		if err != nil {
			return err
		} else {

			//更新数据库
			sql := fmt.Sprintf(`DELETE FROM "%s"."%s" WHERE "id"=$%d`, model.schema, model.table)
			_, err := tx.Exec(sql, key)
			noggo.Logger.Debug("change", "remove", err)
			if err != nil {
				return err
			} else {

				//提交事务
				err := model.base.Submit()
				if err != nil {
					return err
				} else {
					//这里应该有触发器

					//成功了
					return nil
				}
			}
		}

	} else {
		return errors.New("数据：删除失败，主键不存在")
	}
}

//查询唯一对象
func (model *PgsqlModel) Entity(id Any) (error,Map) {

	//开启事务
	err, tx := model.base.Begin()
	noggo.Logger.Debug("data", "entity", "begin", err)
	if err != nil {
		return err,nil
	} else {


		//先拿字段列表
		//不能用*，必须指定字段列表
		//要不然下拉scan的时候，数据库返回的字段和顺序不一定对
		keys := []string{}
		for k,_ := range model.fields {
			keys = append(keys, k)
		}

		sql := fmt.Sprintf(`SELECT "%s" FROM "%s"."%s" WHERE "id"=$1`, strings.Join(keys, `","`), model.schema, model.table)
		row := tx.QueryRow(sql, id)
		if row == nil {
			return errors.New("数据：查询失败"),nil
		} else {

			//扫描数据
			values := make([]interface{}, len(keys))	//真正的值
			pValues := make([]interface{}, len(keys))	//指针，指向值
			for i := range values {
				pValues[i] = &values[i]
			}

			err := row.Scan(pValues...)
			noggo.Logger.Debug("data", "entity", err, sql)
			if err != nil {
				return errors.New("数据：查询时扫描失败 " + err.Error()),nil
			} else {
				m := Map{}
				for i,n := range keys {
					switch v := values[i].(type) {
					case []byte: {
						m[n] = string(v)
					}
					default:
						m[n] = v
					}

					//来调用type来处理值
					//m[n] = model.fielding(m[n], model.field[n].(Map))
				}

				//返回前使用代码生成
				//有必要的, 按模型拿到数据
				item := Map{}
				err := noggo.Mapping.Parse([]string{}, model.fields, m, item)
				noggo.Logger.Debug("data", "entity", "mapping", err)
				if err == nil {
					return nil,item
				} else {
					//如果生成失败,还是返回原始返回值
					//要不然,存在的也显示为不存在
					return nil,m
				}
			}
		}


	}
}