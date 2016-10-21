package data_pgsql


import (
	. "github.com/nogio/noggo/base"
	//"github.com/nogio/noggo/driver"
	"github.com/nogio/noggo"
	"fmt"
	"strings"
	"errors"
)

type (
	PgsqlView struct {
		base    *PgsqlBase
		name    string  //模型名称
		schema  string  //架构名
		table   string  //表名
		key     string  //主键
		fields  Map     //字段定义
	}
)


//统计数量
func (view *PgsqlView) Count(args ...Map) (int64,error) {

	//生成查询条件
	where,builds,err := view.base.building(1,args...)
	if err != nil {
		return int64(0),err
	} else {

		//开启事务
		tx,err := view.base.begin()
		noggo.Logger.Debug("data", "count", "begin", err)
		if err != nil {
			return int64(0),err
		} else {

			sql := fmt.Sprintf(`SELECT COUNT(*) FROM "%s"."%s" WHERE %s`, view.schema, view.table, where)
			row := tx.QueryRow(sql, builds...)
			if row == nil {
				return int64(0),errors.New("数据：查询失败")
			} else {

				count := int64(0)

				err := row.Scan(&count)
				noggo.Logger.Debug("data", "count", err, sql)
				if err != nil {
					return count,errors.New("数据：查询时扫描失败 " + err.Error())
				} else {
					return count,nil
				}
			}
		}
	}
}


//查询单条
func (view *PgsqlView) Single(args ...Map) (Map,error) {

	//生成查询条件
	where,builds,err := view.base.building(1,args...)
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

			sql := fmt.Sprintf(`SELECT "%s" FROM "%s"."%s" WHERE %s`, strings.Join(keys, `","`), view.schema, view.table, where)
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
					err := noggo.Mapping.Parse([]string{}, view.fields, m, item)
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
func (view *PgsqlView) Query(args ...Map) ([]Map,error) {
	//生成查询条件
	where,builds,err := view.base.building(1,args...)
	if err != nil {
		return nil,err
	} else {

		//开启事务
		tx,err := view.base.begin()
		noggo.Logger.Debug("data", "query", "begin", err)
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

			sql := fmt.Sprintf(`SELECT "%s" FROM "%s"."%s" WHERE %s`, strings.Join(keys, `","`), view.schema, view.table, where)
			rows,err := tx.Query(sql, builds...)
			noggo.Logger.Debug("data", "query", err, sql)
			if err != nil {
				return nil,err
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
						return nil, errors.New("数据：查询时扫描失败 " + err.Error())
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
						err := noggo.Mapping.Parse([]string{}, view.fields, m, item)
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
func (view *PgsqlView) Limit(offset,limit Any, args ...Map) ([]Map,error) {
	//生成查询条件
	where,builds,err := view.base.building(1,args...)
	if err != nil {
		return nil,err
	} else {

		//开启事务
		tx,err := view.base.begin()
		noggo.Logger.Debug("data", "limit", "begin", err)
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

			sql := fmt.Sprintf(`SELECT "%s" FROM "%s"."%s" WHERE %s OFFSET %d LIMIT %d`, strings.Join(keys, `","`), view.schema, view.table, where, offset, limit)
			rows,err := tx.Query(sql, builds...)
			noggo.Logger.Debug("data", "limit", err, sql)
			if err != nil {
				return nil,err
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
						return nil, errors.New("数据：查询时扫描失败 " + err.Error())
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
						err := noggo.Mapping.Parse([]string{}, view.fields, m, item)
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