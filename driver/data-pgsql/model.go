package data_pgsql


import (
	. "github.com/nogio/noggo/base"
	//"github.com/nogio/noggo/driver"
	"fmt"
	"strings"
	"github.com/nogio/noggo"
	"errors"
)

type (
	PgsqlModel struct {
		PgsqlView
	}
)



//创建对象
func (model *PgsqlModel) Create(data Map) (Map,error) {

	//按字段生成值
	value := Map{}
	err := noggo.Mapping.Parse([]string{}, model.fields, data, value);
	noggo.Logger.Debug("create", model.fields, data)

	if err != nil {
		return nil,err
	} else {

		//对拿到的值进行包装，以适合pgsql
		newValue := model.base.packing(value)

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

		tx,err := model.base.begin()
		if err != nil {
			return nil,err
		} else {

			sql := fmt.Sprintf(`INSERT INTO "%s"."%s" ("%s") VALUES (%s) RETURNING "id";`, model.schema, model.table, strings.Join(keys, `","`), strings.Join(tags, `,`))
			row := tx.QueryRow(sql, vals...)
			if row == nil {
				return nil,errors.New("数据：插入：无返回行")
			} else {

				id := int64(0)
				err := row.Scan(&id)
				if err != nil {
					//扫描新ID失败
					return nil,err
				} else {
					value["id"] = id

					//提交事务
					err := model.base.Submit()
					if err != nil {
						return nil,err
					} else {

						//这里应该有触发器

						//成功了
						return value,nil

					}
				}
			}
		}
	}
}

//修改对象
func (model *PgsqlModel) Change(item Map, data Map) (Map,error) {


	//按字段生成值
	value := Map{}
	err := noggo.Mapping.Parse([]string{}, model.fields, data, value, true);
	noggo.Logger.Debug("change", "mapping", err)

	if err != nil {
		return nil,err
	} else {

		//包装值，因为golang本身数据类型和数据库的不一定对版
		//需要预处理一下
		newValue := model.base.packing(value)

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
		tx,err := model.base.begin()
		if err != nil {
			return nil,err
		} else {

			//更新数据库
			sql := fmt.Sprintf(`UPDATE "%s"."%s" SET %s WHERE "id"=$%d`, model.schema, model.table, strings.Join(sets, `,`), i)
			_, err := tx.Exec(sql, vals...)
			noggo.Logger.Debug("change", "exec", err)
			if err != nil {
				return nil,err
			} else {

				// 不改item
				// 先复制item
				newItem := Map{}
				for k, v := range item { newItem[k] = v }
				for k, v := range value { newItem[k] = v }


				//提交事务
				err := model.base.Submit()
				if err != nil {
					return nil,err
				} else {

					//这里应该有触发器

					//成功了
					return newItem,nil

				}


			}
		}
	}
}

//删除对象
func (model *PgsqlModel) Remove(item Map) (error) {

	if key,ok := item[model.key]; ok {

		//开启事务
		tx,err := model.base.begin()
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
func (model *PgsqlModel) Entity(id Any) (Map,error) {

	//开启事务
	tx,err := model.base.begin()
	noggo.Logger.Debug("data", "entity", "begin", err)
	if err != nil {
		return nil,err
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
			return nil,errors.New("数据：查询失败")
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
				return nil,errors.New("数据：查询时扫描失败 " + err.Error())
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
				}

				//返回前使用代码生成
				//有必要的, 按模型拿到数据
				item := Map{}
				err := noggo.Mapping.Parse([]string{}, model.fields, m, item)
				noggo.Logger.Debug("data", "entity", "mapping", err)
				if err == nil {
					return item,nil
				} else {
					//如果生成失败,还是返回原始返回值
					//要不然,存在的也显示为不存在
					return m,nil
				}
			}
		}


	}
}






//批量删除
func (model *PgsqlModel) Delete(args ...Map) (int64,error) {

	//生成条件
	where,builds,err := model.base.building(1,args...)
	if err != nil {
		return int64(0),err
	} else {

		//开启事务
		tx,err := model.base.begin()
		noggo.Logger.Debug("data", "count", "begin", err)
		if err != nil {
			return int64(0),err
		} else {

			sql := fmt.Sprintf(`DELETE FROM "%s"."%s" WHERE %s`, model.schema, model.table, where)
			result,err := tx.Exec(sql, builds...)
			if err != nil {
				return int64(0),err
			} else {

				//提交事务
				err := model.base.Submit()
				if err != nil {
					return int64(0),err
				} else {

					return result.RowsAffected()

				}
			}
		}
	}
}


//批量更新
func (model *PgsqlModel) Update(args ...Map) (int64,error) {

	//注意，args[0]为更新的内容，之后的为查询条件
	data := args[0]
	args = args[1:]


	//按字段生成值
	value := Map{}
	err := noggo.Mapping.Parse([]string{}, model.fields, data, value, true);
	noggo.Logger.Debug("data", "update", "mapping", err)

	if err != nil {
		return int64(0),err
	} else {

		//包装值，因为golang本身数据类型和数据库的不一定对版
		//需要预处理一下
		newValue := model.base.packing(value)

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

		//生成条件
		where,builds,err := model.base.building(i, args...)
		if err != nil {
			return int64(0),err
		} else {

			//把builds的args加到vals中
			for _,v := range builds {
				vals = append(vals, v)
			}

			//开启事务
			tx, err := model.base.begin()
			if err != nil {
				return int64(0),err
			} else {

				//更新数据库
				sql := fmt.Sprintf(`UPDATE "%s"."%s" SET %s WHERE %s`, model.schema, model.table, strings.Join(sets, `,`), where)
				result, err := tx.Exec(sql, vals...)
				noggo.Logger.Debug("data", "update", "exec", err)
				if err != nil {
					return int64(0),err
				} else {


					//提交事务
					err := model.base.Submit()
					if err != nil {
						return int64(0),err
					} else {

						return result.RowsAffected()

					}

				}
			}
		}
	}

}
