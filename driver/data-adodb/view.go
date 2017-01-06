package data_adodb


import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo"
	"fmt"
	"strings"
	"errors"
	"time"
)

type (
	AdodbView struct {
		base    *AdodbBase
		name    string  //模型名称
		schema  string  //架构名
		object   string  //这里可能是表名，视图名，或是集合名（mongodb)
		key     string  //主键
		fields  Map     //字段定义
	}
)







//统计数量
func (view *AdodbView) Count(args ...Any) (int64,error) {

	//生成查询条件
	where,builds,_,err := view.base.parsing(args...)
	if err != nil {
		return int64(0),err
	} else {

		//开启事务
		tx,err := view.base.begin()
		noggo.Logger.Debug("data", "count", "begin", err)
		if err != nil {
			return int64(0),err
		} else {

			sql := fmt.Sprintf("SELECT COUNT(*) FROM `%s`.`%s` WHERE %s", view.schema, view.object, where)
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
func (view *AdodbView) Single(args ...Any) (Map,error) {

	//生成查询条件
	where,builds,orderby,err := view.base.parsing(args...)
	if err != nil {
		return nil,err
	} else {

		//开启事务
		tx,err := view.base.begin()
		noggo.Logger.Debug("data", "single", "begin", err)
		if err != nil {
			return nil,err
		} else {

			//先拿字段列表
			//不能用*，必须指定字段列表
			//要不然下拉scan的时候，数据库返回的字段和顺序不一定对
			keys := []string{}
			for k,_ := range view.fields {
				keys = append(keys, k)
			}

			sql := fmt.Sprintf("SELECT `%s` FROM `%s`.`%s` WHERE %s %s", strings.Join(keys, "`,`"), view.schema, view.object, where, orderby)
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
					err := noggo.Mapping.Parse([]string{}, view.fields, m, item, false, true)
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
//查询列表
func (view *AdodbView) Query(args ...Any)  ([]Map,error) {
	//生成查询条件
	where,builds,orderby,err := view.base.parsing(args...)
	if err != nil {
		return []Map{},err
	} else {

		//开启事务
		tx,err := view.base.begin()
		noggo.Logger.Debug("data", "query", "begin", err)
		if err != nil {
			return []Map{},err
		} else {

			//先拿字段列表
			//不能用*，必须指定字段列表
			//要不然下拉scan的时候，数据库返回的字段和顺序不一定对
			keys := []string{}
			for k,_ := range view.fields {
				keys = append(keys, k)
			}

			sql := fmt.Sprintf("SELECT `%s` FROM `%s`.`%s` WHERE %s %s", strings.Join(keys, "`,`"), view.schema, view.object, where, orderby)
			rows,err := tx.Query(sql, builds...)
			noggo.Logger.Debug("data", "query", err, sql)
			if err != nil {
				return []Map{},err
			} else {
				defer rows.Close()

				//返回结果在这
				items := []Map{}

				//遍历结果
				for rows.Next() {
					//扫描数据
					values := make([]interface{}, len(keys))    //真正的值
					pValues := make([]interface{}, len(keys))    //指针，指向值
					for i := range values {
						pValues[i] = &values[i]
					}
					err := rows.Scan(pValues...)

					if err != nil {
						return nil, errors.New("数据：查询时扫描失败")
					} else {
						m := Map{}
						for i, n := range keys {
							switch v := values[i].(type) {
							case []byte: {
								m[n] = string(v)
							}
							case time.Time: {
								m[n] = v
							}
							default:
								m[n] = v
							}
						}

						//返回前使用代码生成
						//有必要的, 按模型拿到数据
						item := Map{}
						err := noggo.Mapping.Parse([]string{}, view.fields, m, item, false, true)
						if err == nil {
							items = append(items, item)
						} else {
							//如果生成失败,还是返回原始返回值
							//要不然,存在的也显示为不存在
							items = append(items, m)
						}
					}
				}

				return items,nil
			}
		}
	}
}


//分页查询
func (view *AdodbView) Limit(offset,limit Any, args ...Any) (int64,[]Map,error) {
	//生成查询条件
	where,builds,orderby,err := view.base.parsing(args...)
	if err != nil {
		return int64(0),[]Map{},err
	} else {

		//开启事务
		tx,err := view.base.begin()
		noggo.Logger.Debug("data", "limit", "begin", err)
		if err != nil {
			return int64(0),[]Map{},err
		} else {

			//先拿字段列表
			//不能用*，必须指定字段列表
			//要不然下拉scan的时候，数据库返回的字段和顺序不一定对
			keys := []string{}
			for k,_ := range view.fields {
				keys = append(keys, k)
			}


			//先统计数据
			sql := fmt.Sprintf("SELECT COUNT(*) FROM `%s`.`%s` WHERE %s", view.schema, view.object, where)
			row := tx.QueryRow(sql, builds...)
			if row == nil {
				return int64(0),[]Map{},errors.New("数据：查询失败")
			} else {

				count := int64(0)

				err := row.Scan(&count)
				noggo.Logger.Debug("data", "count", err, sql)
				if err != nil {
					return count,[]Map{},errors.New("数据：查询时扫描失败")
				} else {


					sql := fmt.Sprintf("SELECT `%s` FROM `%s`.`%s` WHERE %s %s LIMIT %d,%d", strings.Join(keys, "`,`"), view.schema, view.object, where, orderby, offset, limit)
					rows,err := tx.Query(sql, builds...)
					noggo.Logger.Debug("data", "limit", err, sql)
					if err != nil {
						return int64(0),[]Map{},err
					} else {
						defer rows.Close()

						//返回结果在这
						items := []Map{}

						//遍历结果
						for rows.Next() {
							//扫描数据
							values := make([]interface{}, len(keys))    //真正的值
							pValues := make([]interface{}, len(keys))    //指针，指向值
							for i := range values {
								pValues[i] = &values[i]
							}
							err := rows.Scan(pValues...)

							if err != nil {
								return count,[]Map{}, errors.New("数据：查询时扫描失败")
							} else {
								m := Map{}
								for i, n := range keys {
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
								err := noggo.Mapping.Parse([]string{}, view.fields, m, item, false, true)
								if err == nil {
									items = append(items, item)
								} else {
									//如果生成失败,还是返回原始返回值
									//要不然,存在的也显示为不存在
									items = append(items, m)
								}
							}
						}

						return count,items,nil
					}

				}
			}

		}
	}
}



//查询列表
func (view *AdodbView) Group(field string, args ...Any) ([]Map,error) {
	//生成查询条件
	where,builds,orderby,err := view.base.parsing(args...)
	if err != nil {
		return []Map{},err
	} else {

		//开启事务
		tx,err := view.base.begin()
		noggo.Logger.Debug("data", "group", "begin", err)
		if err != nil {
			return []Map{},err
		} else {

			//暂时只支持字段本身， 后续支持 count,sum,avg,max,min啥啥的
			keys := []string{ field }

			sql := fmt.Sprintf("SELECT `%s` FROM `%s`.`%s` WHERE %s GROUP BY `%s` %s", field, view.schema, view.object, where, field, orderby)
			rows,err := tx.Query(sql, builds...)
			noggo.Logger.Debug("data", "group", err, sql)
			if err != nil {
				return []Map{},err
			} else {
				defer rows.Close()

				//返回结果在这
				items := []Map{}

				//遍历结果
				for rows.Next() {

					//扫描数据
					values := make([]interface{}, len(keys))    //真正的值
					pValues := make([]interface{}, len(keys))    //指针，指向值
					for i := range values {
						pValues[i] = &values[i]
					}
					err := rows.Scan(pValues...)

					if err != nil {
						return []Map{},errors.New("数据：查询时扫描失败")
					} else {
						m := Map{}
						for i, n := range keys {
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
						err := noggo.Mapping.Parse([]string{}, view.fields, m, item, false, true)
						if err == nil {
							items = append(items, item)
						} else {
							//如果生成失败,还是返回原始返回值
							//要不然,存在的也显示为不存在
							items = append(items, m)
						}
					}
				}

				return items,nil
			}
		}
	}
}






//查询唯一对象
func (view *AdodbView) Entity(id Any) (Map,error) {

	//开启事务
	tx,err := view.base.begin()
	noggo.Logger.Debug("data", "entity", "begin", err)
	if err != nil {
		return nil,err
	} else {

		//先拿字段列表
		//不能用*，必须指定字段列表
		//要不然下拉scan的时候，数据库返回的字段和顺序不一定对
		keys := []string{}
		for k,_ := range view.fields {
			keys = append(keys, k)
		}

		sql := fmt.Sprintf("SELECT `%s` FROM `%s`.`%s` WHERE `%s`=?", strings.Join(keys, "`,`"), view.schema, view.object, view.key)
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
				err := noggo.Mapping.Parse([]string{}, view.fields, m, item, false, true)
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

