package data_oracle


import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo"
	"fmt"
	"strings"
	"errors"
	"github.com/nogio/noggo/depend/go-oci8"
)

type (
	OracleModel struct {
		OracleView
		/*
		base    *OracleBase
		name    string  //模型名称
		schema  string  //架构名
		object   string  //这里可能是表名，视图名，或是集合名（mongodb)
		key     string  //主键
		fields  Map     //字段定义
		*/
	}
)






//创建对象
func (model *OracleModel) Create(data Map) (Map,error) {

	//按字段生成值
	value := Map{}
	err := noggo.Mapping.Parse([]string{}, model.fields, data, value, false, false);
	noggo.Logger.Debug("create", "mapping", err, model.fields, data)

	if err != nil {
		return nil,err
	} else {

		//对拿到的值进行包装，以适合oracle
		newValue := model.base.packing(value)

		//先拿字段列表
		keys, tags, vals := []string{}, []string{}, make([]interface{},0)
		i := 1
		for k,v := range newValue {
			if k == model.key {
				//如果用seq，主
				if v == nil || model.seq != ""{
					continue
				}
				//id不直接跳过,可以指定ID
				//continue
			}
			keys = append(keys, k)
			vals = append(vals, v)
			tags = append(tags, fmt.Sprintf(":%d", i))
			i++
		}

		tx,err := model.base.begin()
		noggo.Logger.Debug("create", "begin", err)
		if err != nil {
			return nil,err
		} else {



			//oracle的自增有点麻烦，目前计划使用SEQ，
			//但是oracle的字段默认值不支持使用SEQ.nextval，所以，要在添加的时候，手动写在SQL里
			//在model注册的时候， 添加  key ,seq  两个参数 ， 如果 有 seq 对象，表示需要自增，如果没有，则INSERT后不返回自动ID


			sql := fmt.Sprintf(`INSERT INTO %s.%s (%s) VALUES (%s)`, model.schema, model.object, strings.Join(keys, `,`), strings.Join(tags, `,`))

			//如果用了自动编号
			if model.seq != "" {
				sql = fmt.Sprintf(`INSERT INTO %s.%s (%s,%s) VALUES (%s.NEXTVAL,%s)`, model.schema, model.object, model.key, strings.Join(keys, `,`), model.seq, strings.Join(tags, `,`))
			}


			result,err := tx.Exec(sql, vals...)
			noggo.Logger.Debug("create", "exec", err, sql)
			if err != nil {
				return nil,errors.New("数据：插入失败")
			} else {

				//如果有seq才查询最后的ID
				if model.seq != "" {
					lastId,err := result.LastInsertId()
					if err != nil {
						return nil,errors.New("数据：插入获取LastId失败")
					} else {

						rowId := oci8.GetLastInsertId(lastId)

						id := int64(0)

						sql := fmt.Sprintf(`SELECT %s FROM %s.%s WHERE rowid=:1`, model.key, model.schema, model.object)
						row := tx.QueryRow(sql, rowId)
						err := row.Scan(&id)
						noggo.Logger.Debug("create.insert.id", err, sql)
						if err != nil {
							//扫描新ID失败
							return nil,err
						} else {
							//拿到ID了
							value[model.key] = id
						}
					}
				}

					//注意这里，如果手动提交事务， 那这里直接返回，是不需要提交的
					if model.base.manual {

						//这里应该保存触发器
						model.base.trigger(TriggerCreate, Map{ "base": model.base.name, "model": model.name, "entity": value })

						//成功了，但是没有提交事务
						return value,nil

					} else {


						//提交事务
						err := model.base.Submit()
						if err != nil {
							return nil,err
						} else {

							//这里应该有触发器
							noggo.Trigger.Touch(TriggerCreate, Map{ "base": model.base.name, "model": model.name, "entity": value })

							//成功了
							return value,nil

						}

					}

			}
		}
	}
}

//修改对象
func (model *OracleModel) Change(item Map, data Map) (Map,error) {

	//按字段生成值
	value := Map{}
	err := noggo.Mapping.Parse([]string{}, model.fields, data, value, true, false);
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
			sets = append(sets, fmt.Sprintf(`%s=:%d`, k, i))
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
			sql := fmt.Sprintf(`UPDATE %s.%s SET %s WHERE %s=:%d`, model.schema, model.object, strings.Join(sets, `,`), model.key, i)
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


				//注意这里，如果手动提交事务， 那这里直接返回，是不需要提交的
				if model.base.manual {

					//这里应该保存触发器
					model.base.trigger(TriggerChange, Map{ "data": model.base.name, "model": model.name, "before": item, "after": newItem })

					//成功了，但是没有提交事务
					return newItem, nil

				} else {

					//提交事务
					err := model.base.Submit()
					if err != nil {
						return nil, err
					} else {

						//这里应该有触发器
						noggo.Trigger.Touch(TriggerChange, Map{ "data": model.base.name, "model": model.name, "before": item, "after": newItem })

						//成功了
						return newItem, nil

					}
				}


			}
		}
	}
}

//删除对象
func (model *OracleModel) Remove(item Map) (error) {

	if id,ok := item[model.key]; ok {

		//开启事务
		tx,err := model.base.begin()
		if err != nil {
			return err
		} else {

			//更新数据库
			sql := fmt.Sprintf(`DELETE FROM %s.%s WHERE %s=:1`, model.schema, model.object, model.key)
			_, err := tx.Exec(sql, id)
			noggo.Logger.Debug("remove", err, sql)
			if err != nil {
				return err
			} else {

				//注意这里，如果手动提交事务， 那这里直接返回，是不需要提交的
				if model.base.manual {

					//这里应该保存触发器
					model.base.trigger(TriggerRemove, Map{ "data": model.base.name, "model": model.name, "entity": item })

					//成功了，但是没有提交事务
					return nil

				} else {

					//提交事务
					err := model.base.Submit()
					if err != nil {
						return err
					} else {

						//这里应该有触发器
						noggo.Trigger.Touch(TriggerRemove, Map{ "data": model.base.name, "model": model.name, "entity": item })

						//成功了
						return nil
					}
				}
			}
		}

	} else {
		return errors.New("数据：删除失败，主键不存在")
	}
}


//恢复对象
func (model *OracleModel) Recover(item Map) (error) {
	return errors.New("不支持恢复")
}



//批量删除
func (model *OracleModel) Delete(args ...Any) (int64,error) {

	//生成条件
	where,builds,_,err := model.base.parsing(1,args...)
	if err != nil {
		return int64(0),err
	} else {

		//开启事务
		tx,err := model.base.begin()
		noggo.Logger.Debug("data", "count", "begin", err)
		if err != nil {
			return int64(0),err
		} else {

			sql := fmt.Sprintf(`DELETE FROM %s.%s WHERE %s`, model.schema, model.object, where)
			result,err := tx.Exec(sql, builds...)
			if err != nil {
				return int64(0),err
			} else {


				//注意这里，如果手动提交事务， 那这里直接返回，是不需要提交的
				if model.base.manual {

					//成功了，但是没有提交事务
					return result.RowsAffected()

				} else {

					//提交事务
					err := model.base.Submit()
					if err != nil {
						return int64(0), err
					} else {

						return result.RowsAffected()

					}
				}
			}
		}
	}
}


//批量更新
func (model *OracleModel) Update(sets Map, args ...Any) (int64,error) {

	//注意，args[0]为更新的内容，之后的为查询条件
	//sets := args[0]
	//args = args[1:]


	//按字段生成值
	value := Map{}
	err := noggo.Mapping.Parse([]string{}, model.fields, sets, value, true, false);
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
			sets = append(sets, fmt.Sprintf(`%s=:%d`, k, i))
			i++
		}

		//生成条件
		where,builds,_,err := model.base.parsing(i, args...)
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
				sql := fmt.Sprintf(`UPDATE %s.%s SET %s WHERE %s`, model.schema, model.object, strings.Join(sets, `,`), where)
				result, err := tx.Exec(sql, vals...)
				noggo.Logger.Debug("data", "update", "exec", err)
				if err != nil {
					return int64(0),err
				} else {

					//注意这里，如果手动提交事务， 那这里直接返回，是不需要提交的
					if model.base.manual {

						//成功了，但是没有提交事务
						return result.RowsAffected()

					} else {

						//提交事务
						err := model.base.Submit()
						if err != nil {
							return int64(0), err
						} else {

							//这是真成功了
							return result.RowsAffected()

						}
					}

				}
			}
		}
	}

}



/*


//查询唯一对象
func (model *OracleModel) Entity(id Any) (Map,error) {

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

		sql := fmt.Sprintf(`SELECT %s FROM %s.%s WHERE %s=:1`, strings.Join(keys, `,`), model.schema, model.object, model.key)
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
				return nil,errors.New("数据：查询时扫描失败 ")
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
				err := noggo.Mapping.Parse([]string{}, model.fields, m, item, false, true)
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


*/