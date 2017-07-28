package data_postgres


import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo"
	"fmt"
	"strings"
	"errors"
	"time"
)

type (
	PostgresTable struct {
		PostgresView
	}
)






//创建对象
func (table *PostgresTable) Create(data Map) (Map,error) {

	//按字段生成值
	value := Map{}
	err := noggo.Mapping.Parse([]string{}, table.fields, data, value, false, false)
	noggo.Logger.Debug("create", err, table.fields, data)

	if err != nil {
		return nil,err
	} else {

		//对拿到的值进行包装，以适合postgres
		newValue := table.base.packing(value)

		//先拿字段列表
		keys, tags, vals := []string{}, []string{}, make([]interface{},0)
		i := 1
		for k,v := range newValue {
			if k == table.key {
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

		tx,err := table.base.begin()
		if err != nil {
			return nil,err
		} else {

			sql := fmt.Sprintf(`INSERT INTO "%s"."%s" ("%s") VALUES (%s) RETURNING "%s";`, table.schema, table.object, strings.Join(keys, `","`), strings.Join(tags, `,`), table.key)
			row := tx.QueryRow(sql, vals...)
			if row == nil {
				return nil,errors.New("数据：插入：无返回行")
			} else {

				id := int64(0)
				err := row.Scan(&id)
				noggo.Logger.Debug("create.insert", err, sql)
				if err != nil {
					//扫描新ID失败
					return nil,err
				} else {
					value[table.key] = id


					//这里处理关联
					//遍历所有父表,加统计
					parents := table.base.parents(table.name)
					for _,v := range parents {
						parent := v.(Map)
						relate := parent["relate"].(Map)
						field := parent["field"].(string)	//如果是mongodb,得是 a.b.c 这样多级,要用一个方法去取值

						//要不要考虑一下,值是数组的情况

						//状态连带
						if relate["count"] != nil && value[field] != nil {

							//加入类型处理
							//非本类型父表, 统计不加1, 因为data_id一个字段,会有3个父表的引用, 每条只加其中之一
							if relate["type"] != nil && value["type"] != relate["type"] {
								continue
							}

							//要处理一下schema
							if parent["schema"]=="public" {
								parent["schema"] = table.schema
							}

							sql = fmt.Sprintf(`UPDATE "%s"."%s" SET "%v"="%v"+1 WHERE "id"=$1`, parent["schema"], parent["table"], relate["count"], relate["count"])

							_, err := tx.Exec(sql, value[field])
							noggo.Logger.Debug("create.cc.err=%v\n", err, sql)
							if err != nil {
								//失败, 如果自动提交 才 回滚事务   非自动提交时, 不回滚, 失败就失败
								//回滚这时候就到逻辑层,手动回了
								//此处尴尬之地
								if table.base.manual == false {
									table.base.Cancel()
								}
								break
							}
						}
					}






					//注意这里，如果手动提交事务， 那这里直接返回，是不需要提交的
					if table.base.manual {

						//这里应该保存触发器
						table.base.trigger(TriggerCreate, Map{ "base": table.base.name, "table": table.name, "entity": value })

						//成功了，但是没有提交事务
						return value,nil

					} else {


						//提交事务
						err := table.base.Submit()
						if err != nil {
							return nil,err
						} else {

							//这里应该有触发器
							noggo.Trigger.Touch(TriggerCreate, Map{ "base": table.base.name, "table": table.name, "entity": value })

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
func (table *PostgresTable) Change(item Map, data Map) (Map,error) {

	//记录修改时间
	if data[FieldChanged] == nil {
		data[FieldChanged] = time.Now()
	}

	//按字段生成值
	value := Map{}
	err := noggo.Mapping.Parse([]string{}, table.fields, data, value, true, false);
	noggo.Logger.Debug("change", "mapping", err)

	if err != nil {
		return nil,err
	} else {

		//包装值，因为golang本身数据类型和数据库的不一定对版
		//需要预处理一下
		newValue := table.base.packing(value)

		if inc,ok := data[INC]; ok {
			newValue[INC] = inc
		}

		//先拿字段列表
		sets, vals := []string{}, make([]interface{}, 0)
		i := 1
		for k, v := range newValue {
			//主值不在修改之中
			if k == table.key {
				continue
			} else if k == INC {
				if vm,ok := v.(Map); ok {
					for kk,vv := range vm {
						vals = append(vals, vv)
						sets = append(sets, fmt.Sprintf(`"%s"="%s"+$%d`, kk, kk, i))
						i++
					}
				}
			} else {
				//keys = append(keys, k)
				vals = append(vals, v)
				sets = append(sets, fmt.Sprintf(`"%s"=$%d`, k, i))
				i++
			}
		}
		//条件是主键
		vals = append(vals, item[table.key])

		//开启事务
		tx,err := table.base.begin()
		if err != nil {
			return nil,err
		} else {

			//更新数据库
			sql := fmt.Sprintf(`UPDATE "%s"."%s" SET %s WHERE "%s"=$%d`, table.schema, table.object, strings.Join(sets, `,`), table.key, i)
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


				//处理父表的统计

				//遍历所有父表,加统计
				parents := table.base.parents(table.name)
				for _,v := range parents {
					parent := v.(Map)
					relate := parent["relate"].(Map)
					field := parent["field"].(string)	//如果是mongodb,得是 a.b.c 这样多级,要用一个方法去取值

					//要处理一下schema
					if parent["schema"]=="public" {
						parent["schema"] = table.schema
					}

					//状态连带
					if relate["count"] != nil {


						//加入类型处理
						//非本类型父表, 统计不加1, 因为data_id一个字段,会有3个父表的引用, 每条只加其中之一
						if relate["type"] != nil && relate["type"] != item["type"] {
							continue
						}


						//要不要考虑一下,值是数组的情况


						//如果老值和新值不一样, 就更新相关的2个父记录
						if item[field] != newItem[field] {

							//如果老值不为空,才减1
							if item[field] != nil {
								//先把老记录统计减1, -1
								sql = fmt.Sprintf(`UPDATE "%s"."%s" SET "%v"="%v"-1 WHERE "id"=$1`, parent["schema"], parent["table"], relate["count"], relate["count"])
								_, oldErr := tx.Exec(sql, item[field])
								noggo.Logger.Debug("change.old.err=%v\n", err, sql)
								if oldErr != nil {
									//失败, 如果自动提交 才 回滚事务   非自动提交时, 不回滚, 失败就失败
									//回滚这时候就到逻辑层,手动回了
									//此处尴尬之地
									if table.base.manual == false {
										table.base.Cancel()
									}
									break
								}
							}

							//如果新值不为空,才加1
							if newItem[field] != nil {
								//再把析记录统计加1, +1
								sql = fmt.Sprintf(`UPDATE "%s"."%s" SET "%v"="%v"+1 WHERE "id"=$1`, parent["schema"], parent["table"], relate["count"], relate["count"])
								_, newErr := tx.Exec(sql, newItem[field])
								noggo.Logger.Debug("change.new.err=%v\n", err, sql)
								if newErr != nil {
									//失败, 如果自动提交 才 回滚事务   非自动提交时, 不回滚, 失败就失败
									//回滚这时候就到逻辑层,手动回了
									//此处尴尬之地
									if table.base.manual == false {
										table.base.Cancel()
									}
									break
								}
							}

						}
					}
				}






				//注意这里，如果手动提交事务， 那这里直接返回，是不需要提交的
				if table.base.manual {

					//这里应该保存触发器
					table.base.trigger(TriggerChange, Map{ "base": table.base.name, "table": table.name, "before": item, "after": newItem })

					//成功了，但是没有提交事务
					return newItem, nil

				} else {

					//提交事务
					err := table.base.Submit()
					if err != nil {
						return nil, err
					} else {

						//这里应该有触发器
						noggo.Trigger.Touch(TriggerChange, Map{ "base": table.base.name, "table": table.name, "before": item, "after": newItem })

						//成功了
						return newItem, nil

					}
				}


			}
		}
	}
}

//删除对象
func (table *PostgresTable) Remove(item Map) (error) {

	if item[StatusField] != nil {
		return errors.New("状态异常")
	} else {

		if id,ok := item[table.key]; ok {

			//开启事务
			tx,err := table.base.begin()
			if err != nil {
				return err
			} else {

				//更新数据库
				//sql := fmt.Sprintf(`DELETE FROM "%s"."%s" WHERE "%s"=$1`, table.schema, table.object, table.key)
				sql := fmt.Sprintf(`UPDATE "%s"."%s" SET "status"=$2 WHERE "%s"=$1`, table.schema, table.object, table.key)
				_, err := tx.Exec(sql, id, StatusRemoved)
				noggo.Logger.Debug("remove", err, sql)
				if err != nil {
					return err
				} else {

					item[StatusField] = StatusRemoved


					//遍历所有父表,减统计
					parents := table.base.parents(table.name)
					for _,v := range parents {
						parent := v.(Map)
						relate := parent["relate"].(Map)
						field := parent["field"].(string)	//如果是mongodb,得是 a.b.c 这样多级,要用一个方法去取值

						//要处理一下schema
						if parent["schema"]=="public" {
							parent["schema"] = table.schema
						}


						//要不要考虑数组外键的情况

						//状态连带
						if relate["count"] != nil && item[field] != nil {

							//加入类型处理
							//非本类型父表, 统计不加1, 因为data_id一个字段,会有3个父表的引用, 每条只加其中之一
							if relate["type"] != nil && relate["type"] != item["type"] {
								continue
							}

							sql = fmt.Sprintf(`UPDATE "%s"."%s" SET "%v"="%v"-1 WHERE "id"=$1`, parent["schema"], parent["table"], relate["count"], relate["count"])
							_, err := tx.Exec(sql, item[field])
							noggo.Logger.Debug("remove", "parent", err, sql)
							if err != nil {
								//失败, 如果自动提交 才 回滚事务   非自动提交时, 不回滚, 失败就失败
								//回滚这时候就到逻辑层,手动回了
								//此处尴尬之地
								if table.base.manual == false {
									table.base.Cancel()
								}
								break
							}
						}
					}


					//再一个个删除子表
					//所有被删除的子表, 他们另外关联的统计,应该当全部重新统计
					//但是如果那样,一条条的连带删除, 太浪费生命,, 当然,也是全部写在事务里的
					//删除和恢复, 都全部重新统计相关的count,  暂时没写
					//每一个子表的每一条记录, 还有可能有额外的外键, 他们这些父表,都应该COUNT-1
					//好乱啊. 楼上这一堆, 要记录的东西太多了
					//所以这个只把子表相关的status变为父表删除状态，其它的不动
					//真实的count信息应该走视图去查询，以保证count的实时性

					//要解决上面的问题， 就要冗余整个父表链的ID
					//比如，user(id), tab(id,user_id), tab_image(id,user_id,tab_id,file_id)
					//理解上也这样设计比较好， 就是数据库设计的时候， 不要关联太多层，要优化设计
					//比如，省/市/区县/街道   4级，哈哈哈，而且这样设计，迁移某中间级，下面所有的冗余外键都需要更新
					//直接单表，id,parent_id,name 最实在

					childs := table.base.childs(table.name)
					for _,v := range childs {
						child := v.(Map)
						relate := child["relate"].(Map)


						//要处理一下schema
						if child["schema"]=="public" {
							child["schema"] = table.schema
						}

						//状态连带
						if relate["status"] != nil {

							//加入类型处理
							if relate["type"] != nil {
								sql = fmt.Sprintf(`UPDATE "%s"."%s" SET "status"=$2 WHERE "type"=$3 AND "%v"=$1 AND "status" IS NULL`, child["schema"], child["table"], child["field"])
								_, err := tx.Exec(sql, item["id"], relate["status"], relate["type"])
								noggo.Logger.Debug("remove", "child", err, sql)
								if err != nil {
									//失败, 如果自动提交 才 回滚事务   非自动提交时, 不回滚, 失败就失败
									//回滚这时候就到逻辑层,手动回了
									//此处尴尬之地
									if table.base.manual == false {
										table.base.Cancel()
									}
									break
								}
							} else {
								sql = fmt.Sprintf(`UPDATE "%s"."%s" SET "status"=$2 WHERE "%v"=$1 AND "status" IS NULL`, child["schema"], child["table"], child["field"])
								_, err := tx.Exec(sql, item["id"], relate["status"])
								noggo.Logger.Debug("remove", "child", err, sql)
								if err != nil {
									//失败, 如果自动提交 才 回滚事务   非自动提交时, 不回滚, 失败就失败
									//回滚这时候就到逻辑层,手动回了
									//此处尴尬之地
									if table.base.manual == false {
										table.base.Cancel()
									}
									break
								}
							}

						}
					}







					//注意这里，如果手动提交事务， 那这里直接返回，是不需要提交的
					if table.base.manual {

						//这里应该保存触发器
						table.base.trigger(TriggerRemove, Map{ "base": table.base.name, "table": table.name, "entity": item })

						//成功了，但是没有提交事务
						return nil

					} else {

						//提交事务
						err := table.base.Submit()
						if err != nil {
							return err
						} else {

							//这里应该有触发器
							noggo.Trigger.Touch(TriggerRemove, Map{ "base": table.base.name, "table": table.name, "entity": item })

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
}


//恢复对象
func (table *PostgresTable) Recover(item Map) (error) {

	if item[StatusField] != StatusRemoved {
		return errors.New("状态异常")
	} else {

		if id,ok := item[table.key]; ok {

			//开启事务
			tx,err := table.base.begin()
			if err != nil {
				return err
			} else {

				//更新数据库
				sql := fmt.Sprintf(`UPDATE "%s"."%s" SET "status"=NULL WHERE "%s"=$1`, table.schema, table.object, table.key)
				_, err := tx.Exec(sql, id)
				noggo.Logger.Debug("remove", err, sql)
				if err != nil {
					return err
				} else {

					item[StatusField] = nil


					//遍历所有父表,加统计
					parents := table.base.parents(table.name)
					for _,v := range parents {
						parent := v.(Map)
						relate := parent["relate"].(Map)
						field := parent["field"].(string)	//如果是mongodb,得是 a.b.c 这样多级,要用一个方法去取值


						//处理schema
						if parent["schema"]=="public" {
							parent["schema"] = table.schema
						}

						//要不要考虑数组外键的情况

						//状态连带
						if relate["count"] != nil && item[field] != nil {


							//加入类型处理
							//非本类型父表, 统计不加1, 因为data_id一个字段,会有3个父表的引用, 每条只加其中之一
							if relate["type"] != nil && relate["type"] != item["type"] {
								continue
							}


							sql = fmt.Sprintf(`UPDATE "%s"."%s" SET "%v"="%v"+1 WHERE "id"=$1`, parent["schema"], parent["table"], relate["count"], relate["count"])
							_, err := tx.Exec(sql, item[field])
							noggo.Logger.Debug("recover", "parent", err, sql)
							if err != nil {
								//失败, 如果自动提交 才 回滚事务   非自动提交时, 不回滚, 失败就失败
								//回滚这时候就到逻辑层,手动回了
								//此处尴尬之地
								if table.base.manual == false {
									table.base.Cancel()
								}
								break
							}
						}
					}

					//再一个个恢复子表
					childs := table.base.childs(table.name)
					for _,v := range childs {
						child := v.(Map)
						relate := child["relate"].(Map)

						//处理schema
						if child["schema"]=="public" {
							child["schema"] = table.schema
						}

						//状态连带
						if relate["status"] != nil {

							if relate["type"] != nil {
								sql = fmt.Sprintf(`UPDATE "%s"."%s" SET "status"=null WHERE "type"=$3 AND "%v"=$1 AND "status"=$2`, child["schema"], child["table"], child["field"])
								_, err := tx.Exec(sql, item["id"], relate["status"], relate["type"])
								noggo.Logger.Debug("recover", "child", err, sql)
								if err != nil {
									//失败, 如果自动提交 才 回滚事务   非自动提交时, 不回滚, 失败就失败
									//回滚这时候就到逻辑层,手动回了
									//此处尴尬之地
									if table.base.manual == false {
										table.base.Cancel()
									}
									break
								}
							} else {
								sql = fmt.Sprintf(`UPDATE "%s"."%s" SET "status"=null WHERE "%v"=$1 AND "status"=$2`, child["schema"], child["table"], child["field"])
								_, err := tx.Exec(sql, item["id"], relate["status"])
								noggo.Logger.Debug("recover", "child", err, sql)
								if err != nil {
									//失败, 如果自动提交 才 回滚事务   非自动提交时, 不回滚, 失败就失败
									//回滚这时候就到逻辑层,手动回了
									//此处尴尬之地
									if table.base.manual == false {
										table.base.Cancel()
									}
									break
								}
							}

						}
					}


					//注意这里，如果手动提交事务， 那这里直接返回，是不需要提交的
					if table.base.manual {

						//这里应该保存触发器
						table.base.trigger(TriggerRecover, Map{ "base": table.base.name, "table": table.name, "entity": item })

						//成功了，但是没有提交事务
						return nil

					} else {

						//提交事务
						err := table.base.Submit()
						if err != nil {
							return err
						} else {

							//这里应该有触发器
							noggo.Trigger.Touch(TriggerRecover, Map{ "base": table.base.name, "table": table.name, "entity": item })

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
}




//批量删除，这可是真删
func (table *PostgresTable) Delete(args ...Any) (int64,error) {

	//生成条件
	where,builds,_,err := table.base.parsing(1,args...)
	if err != nil {
		return int64(0),err
	} else {

		//开启事务
		tx,err := table.base.begin()
		noggo.Logger.Debug("data", "count", "begin", err)
		if err != nil {
			return int64(0),err
		} else {

			sql := fmt.Sprintf(`DELETE FROM "%s"."%s" WHERE %s`, table.schema, table.object, where)
			result,err := tx.Exec(sql, builds...)
			if err != nil {
				return int64(0),err
			} else {


				//注意这里，如果手动提交事务， 那这里直接返回，是不需要提交的
				if table.base.manual {

					//成功了，但是没有提交事务
					return result.RowsAffected()

				} else {

					//提交事务
					err := table.base.Submit()
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


//批量更新，直接更了， 没有任何relate相关处理的
func (table *PostgresTable) Update(sets Map, args ...Any) (int64,error) {

	//注意，args[0]为更新的内容，之后的为查询条件
	//sets := args[0]
	//args = args[1:]


	//按字段生成值
	value := Map{}
	err := noggo.Mapping.Parse([]string{}, table.fields, sets, value, true, false)
	noggo.Logger.Debug("data", "update", "mapping", err)

	if err != nil {
		return int64(0),err
	} else {

		//包装值，因为golang本身数据类型和数据库的不一定对版
		//需要预处理一下
		newValue := table.base.packing(value)

		if inc,ok := sets[INC]; ok {
			newValue[INC] = inc
		}

		//先拿字段列表
		sets, vals := []string{}, make([]interface{}, 0)
		i := 1
		for k, v := range newValue {
			//主值不在修改之中
			if k == table.key {
				continue
			} else if k == INC {
				if vm,ok := v.(Map); ok {
					for kk,vv := range vm {
						vals = append(vals, vv)
						sets = append(sets, fmt.Sprintf(`"%s"="%s"+$%d`, kk, kk, i))
						i++
					}
				}
			} else {
				//keys = append(keys, k)
				vals = append(vals, v)
				sets = append(sets, fmt.Sprintf(`"%s"=$%d`, k, i))
				i++
			}
		}

		//生成条件
		where,builds,_,err := table.base.parsing(i, args...)
		if err != nil {
			return int64(0),err
		} else {

			//把builds的args加到vals中
			for _,v := range builds {
				vals = append(vals, v)
			}

			//开启事务
			tx, err := table.base.begin()
			if err != nil {
				return int64(0),err
			} else {

				//更新数据库
				sql := fmt.Sprintf(`UPDATE "%s"."%s" SET %s WHERE %s`, table.schema, table.object, strings.Join(sets, `,`), where)
				result, err := tx.Exec(sql, vals...)
				noggo.Logger.Debug("data", "update", "exec", err)
				if err != nil {
					return int64(0),err
				} else {

					//注意这里，如果手动提交事务， 那这里直接返回，是不需要提交的
					if table.base.manual {

						//成功了，但是没有提交事务
						return result.RowsAffected()

					} else {

						//提交事务
						err := table.base.Submit()
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





