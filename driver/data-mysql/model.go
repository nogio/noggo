package data_mysql


import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo"
	"fmt"
	"strings"
	"errors"
)

type (
	MysqlModel struct {
		MysqlView
		/*
		base    *MysqlBase
		name    string  //模型名称
		schema  string  //架构名
		object   string  //这里可能是表名，视图名，或是集合名（mongodb)
		key     string  //主键
		fields  Map     //字段定义
		*/
	}
)






//创建对象
func (model *MysqlModel) Create(data Map) (Map,error) {

	//按字段生成值
	value := Map{}
	err := noggo.Mapping.Parse([]string{}, model.fields, data, value, false, false);
	noggo.Logger.Debug("create", model.fields, data)

	if err != nil {
		return nil,err
	} else {

		//对拿到的值进行包装，以适合pgsql
		newValue := model.base.packing(value)

		//先拿字段列表
		keys, tags, vals := []string{}, []string{}, make([]interface{},0)
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
			tags = append(tags, fmt.Sprintf("?"))
		}

		tx,err := model.base.begin()
		if err != nil {
			return nil,err
		} else {

			sql := fmt.Sprintf("INSERT INTO `%s`.`%s` (`%s`) VALUES (%s)", model.schema, model.object, strings.Join(keys, "`,`"), strings.Join(tags, `,`))
			result,err := tx.Exec(sql, vals...)
			noggo.Logger.Debug("create", err, sql)
			if err != nil {
				return nil,errors.New("数据：插入失败 " + err.Error())
			} else {

				id,err := result.LastInsertId()
				if err != nil {
					//扫描新ID失败
					return nil,errors.New("数据：获取ID失败 " + err.Error())
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
func (model *MysqlModel) Change(item Map, data Map) (Map,error) {


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
		for k, v := range newValue {
			//主值不在修改之中
			if k == model.key {
				continue
			}
			//keys = append(keys, k)
			vals = append(vals, v)
			sets = append(sets, fmt.Sprintf("`%s`=?", k))
		}
		//条件是主键
		vals = append(vals, item[model.key])

		//开启事务
		tx,err := model.base.begin()
		if err != nil {
			return nil,err
		} else {

			//更新数据库
			sql := fmt.Sprintf("UPDATE `%s`.`%s` SET %s WHERE `%s`=?", model.schema, model.object, strings.Join(sets, `,`), model.key)
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
func (model *MysqlModel) Remove(item Map) (error) {

	if id,ok := item[model.key]; ok {

		//开启事务
		tx,err := model.base.begin()
		if err != nil {
			return err
		} else {

			//更新数据库
			sql := fmt.Sprintf("DELETE FROM `%s`.`%s` WHERE `%s`=?", model.schema, model.object, model.key)
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


//恢复对象
func (model *MysqlModel) Recover(item Map) (error) {
	return errors.New("数据：不支持Recover")
}




//批量删除
func (model *MysqlModel) Delete(args ...Any) (int64,error) {

	//生成条件
	where,builds,_,err := model.base.parsing(args...)
	if err != nil {
		return int64(0),err
	} else {

		//开启事务
		tx,err := model.base.begin()
		noggo.Logger.Debug("data", "count", "begin", err)
		if err != nil {
			return int64(0),err
		} else {

			sql := fmt.Sprintf("DELETE FROM `%s`.`%s` WHERE %s", model.schema, model.object, where)
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
func (model *MysqlModel) Update(sets Map, args ...Any) (int64,error) {

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
		for k, v := range newValue {
			//主值不在修改之中
			if k == model.key {
				continue
			}
			//keys = append(keys, k)
			vals = append(vals, v)
			sets = append(sets, fmt.Sprintf("`%s`=?", k))
		}

		//生成条件
		where,builds,_,err := model.base.parsing(args...)
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
				sql := fmt.Sprintf("UPDATE `%s`.`%s` SET %s WHERE %s", model.schema, model.object, strings.Join(sets, `,`), where)
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


