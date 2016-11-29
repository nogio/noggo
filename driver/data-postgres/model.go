package data_postgres


import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo"
	"fmt"
	"strings"
	"errors"
)

type (
	PostgresModel struct {
		PostgresView
		/*
		base    *PostgresBase
		name    string  //模型名称
		schema  string  //架构名
		object   string  //这里可能是表名，视图名，或是集合名（mongodb)
		key     string  //主键
		fields  Map     //字段定义
		*/
	}
)






//创建对象
func (model *PostgresModel) Create(data Map) (Map,error) {

	//按字段生成值
	value := Map{}
	err := noggo.Mapping.Parse([]string{}, model.fields, data, value);
	noggo.Logger.Debug("create", model.fields, data)

	if err != nil {
		return nil,err
	} else {

		//对拿到的值进行包装，以适合postgres
		newValue := model.base.packing(value)

		//先拿字段列表
		keys, tags, vals := []string{}, []string{}, make([]interface{},0)
		i := 1
		for k,v := range newValue {
			if k == model.key {
				if v == nil {
					continue
				}
				//id不直接跳过,可以指定ID
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

			sql := fmt.Sprintf(`INSERT INTO "%s"."%s" ("%s") VALUES (%s) RETURNING "%s";`, model.schema, model.object, strings.Join(keys, `","`), strings.Join(tags, `,`), model.key)
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
					value[model.key] = id


					//注意这里，如果手动提交事务， 那这里直接返回，是不需要提交的
					if model.base.manual {

						//这里应该保存触发器

						//成功了，但是没有提交事务
						return value,nil

					} else {


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
}

//修改对象
func (model *PostgresModel) Change(item Map, data Map) (Map,error) {


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
			sql := fmt.Sprintf(`UPDATE "%s"."%s" SET %s WHERE "%s"=$%d`, model.schema, model.object, strings.Join(sets, `,`), model.key, i)
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

					//成功了，但是没有提交事务
					return newItem, nil

				} else {

					//提交事务
					err := model.base.Submit()
					if err != nil {
						return nil, err
					} else {

						//这里应该有触发器

						//成功了
						return newItem, nil

					}
				}


			}
		}
	}
}

//删除对象
func (model *PostgresModel) Remove(item Map) (error) {

	if id,ok := item[model.key]; ok {

		//开启事务
		tx,err := model.base.begin()
		if err != nil {
			return err
		} else {

			//更新数据库
			sql := fmt.Sprintf(`DELETE FROM "%s"."%s" WHERE "%s"=$1`, model.schema, model.object, model.key)
			_, err := tx.Exec(sql, id)
			noggo.Logger.Debug("change", "remove", err)
			if err != nil {
				return err
			} else {

				//注意这里，如果手动提交事务， 那这里直接返回，是不需要提交的
				if model.base.manual {

					//这里应该保存触发器

					//成功了，但是没有提交事务
					return nil

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
		}

	} else {
		return errors.New("数据：删除失败，主键不存在")
	}
}

//查询唯一对象
func (model *PostgresModel) Entity(id Any) (Map,error) {

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

		sql := fmt.Sprintf(`SELECT "%s" FROM "%s"."%s" WHERE "%s"=$1`, strings.Join(keys, `","`), model.schema, model.object, model.key)
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
func (model *PostgresModel) Delete(args ...Any) (int64,error) {

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

			sql := fmt.Sprintf(`DELETE FROM "%s"."%s" WHERE %s`, model.schema, model.object, where)
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
func (model *PostgresModel) Update(sets Map, args ...Any) (int64,error) {

	//注意，args[0]为更新的内容，之后的为查询条件
	//sets := args[0]
	//args = args[1:]


	//按字段生成值
	value := Map{}
	err := noggo.Mapping.Parse([]string{}, model.fields, sets, value, true);
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
				sql := fmt.Sprintf(`UPDATE "%s"."%s" SET %s WHERE %s`, model.schema, model.object, strings.Join(sets, `,`), where)
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








//统计数量
func (model *PostgresModel) Count(args ...Any) (int64,error) {

	//生成查询条件
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

			sql := fmt.Sprintf(`SELECT COUNT(*) FROM "%s"."%s" WHERE %s`, model.schema, model.object, where)
			row := tx.QueryRow(sql, builds...)
			if row == nil {
				return int64(0),errors.New("数据：查询失败")
			} else {

				count := int64(0)

				err := row.Scan(&count)
				noggo.Logger.Debug("data", "count", err, sql)
				if err != nil {
					return count,errors.New("数据：查询时扫描失败")
				} else {
					return count,nil
				}
			}
		}
	}
}


//查询单条
func (model *PostgresModel) Single(args ...Any) (Map,error) {

	//生成查询条件
	where,builds,orderby,err := model.base.parsing(1,args...)
	if err != nil {
		return nil,err
	} else {

		//开启事务
		tx,err := model.base.begin()
		noggo.Logger.Debug("data", "single", "begin", err)
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

			sql := fmt.Sprintf(`SELECT "%s" FROM "%s"."%s" WHERE %s %s`, strings.Join(keys, `","`), model.schema, model.object, where, orderby)
			row := tx.QueryRow(sql, builds...)
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
				noggo.Logger.Debug("data", "single", err, sql)
				if err != nil {
					return nil,errors.New("数据：查询时扫描失败")
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
					noggo.Logger.Debug("data", "single", "mapping", err)
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
}