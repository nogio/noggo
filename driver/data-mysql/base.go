package data_mysql


import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo/driver"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"encoding/json"
	"strconv"
)

type (
	MysqlBase struct {
		name    string
		conn    *MysqlConnect
		models  map[string]Map

		db      *sql.DB
		tx      *sql.Tx
		cache   driver.CacheBase

		//是否手动提交事务，否则为自动
		//当调用begin时， 自动变成手动提交事务
		manual    bool
	}
)


//关闭数据库
func (base *MysqlBase) Close() {

	//好像目前不需要关闭什么东西
	if base.tx != nil {
		//关闭时候,一定要提交一次事务
		//如果手动提交了, 这里会失败, 问题不大
		//如果没有提交的话, 连接不会交回连接池. 会一直占用
		base.Cancel()
	}
}

//获取模型对象
func (base *MysqlBase) Model(name string) (driver.DataModel) {
	if config,ok := base.models[name]; ok {

		//模式，表名
		schema, object, key, fields := "public", name, "id", Map{}
		if n,ok := config["schema"].(string); ok {
			schema = n
		}
		if n,ok := config["object"].(string); ok {
			object = n
		}
		if n,ok := config["key"].(string); ok {
			key = n
		}
		if n,ok := config["fields"].(Map); ok {
			fields = n
		}

		return &MysqlModel{
			base, name, schema, object, key, fields,
		}
	} else {
		panic("数据：模型不存在")
	}
}



//开启手动模式
func (base *MysqlBase) Begin() (driver.DataBase) {
	base.manual = true
	return base
}



//注意，此方法为实际开始事务
func (base *MysqlBase) begin() (*sql.Tx,error) {

	if base.tx == nil {
		tx,err := base.db.Begin()
		if err != nil {
			return nil,err
		}
		base.tx = tx
	}
	return base.tx,nil
}


//提交事务
func (base *MysqlBase) Submit() (error) {

	if base.tx == nil {
		return errors.New("数据：tx未开始")
	}

	err := base.tx.Commit()
	if err != nil {
		return err
	}

	//提交后,要清掉事务
	base.tx = nil

	return nil
}

//取消事务
func (base *MysqlBase) Cancel() (error) {

	if base.tx == nil {
		return errors.New("数据：tx未开始")
	}

	err := base.tx.Rollback()
	if err != nil {
		return err
	}

	//提交后,要清掉事务
	base.tx = nil

	return nil
}







//Exec
func (base *MysqlBase) Exec(query string, args ...interface{}) (sql.Result,error) {
	if base.tx == nil {
		return nil,errors.New("数据：tx未开始")
	}
	return base.tx.Exec(query, args...)
}

//Prepare
func (base *MysqlBase) Prepare(query string) (*sql.Stmt, error) {
	if base.tx == nil {
		return nil,errors.New("数据：tx未开始")
	}
	return base.tx.Prepare(query)
}

//Query
func (base *MysqlBase) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if base.tx == nil {
		return nil,errors.New("数据：tx未开始")
	}
	return base.tx.Query(query, args...)
}

//QueryRow
func (base *MysqlBase) QueryRow(query string, args ...interface{}) (*sql.Row) {
	if base.tx == nil {
		return nil
	}
	return base.tx.QueryRow(query, args...)
}

//QueryRow
func (base *MysqlBase) Stmt(stmt *sql.Stmt) (*sql.Stmt) {
	if base.tx == nil {
		return nil
	}
	return base.tx.Stmt(stmt)
}




















//创建的时候,也需要对值来处理,
//数组要转成{a,b,c}格式,要不然不支持
//json可能要转成字串才支持
func (base *MysqlBase) packing(value Map) (Map) {

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







//把MAP编译成sql查询条件
//加入排序
//where,args,order,error
func (base *MysqlBase) building(args ...Map) (string,[]interface{},string,error) {

	if len(args) > 0 {
		querys := []string{}
		values := make([]interface{}, 0)
		orders := []string{}

		//否则是多个map,单个为 与, 多个为 或
		for _,m := range args {
			for k,v := range m {


				//如果值是ASC,DESC，表示是排序
				if ov,ok := v.(string); ok && (ov==driver.ASC || ov==driver.DESC) {

					if ov == driver.ASC {
						orders = append(orders, fmt.Sprintf("`%s` ASC", k))
					} else {
						orders = append(orders, fmt.Sprintf("`%s` DESC", k))
					}

				} else {


					ands := []string{}

					//v要处理一下如果是map要特别处理
					//key做为操作符，比如 > < >= 等
					//而且多个条件是and，比如 views > 1 AND views < 100
					if opMap, opOK := v.(Map); opOK {

						opAnds := []string{}
						for opKey,opVal := range opMap {
							opAnds = append(opAnds, fmt.Sprintf("`%s` %s ?", k, opKey))
							values = append(values, opVal)
						}
						ands = append(ands, fmt.Sprintf("(%s)", strings.Join(opAnds, " AND ")))

					} else {

						if v == nil {
							ands = append(ands, fmt.Sprintf("`%s` IS NULL", k))
						} else {
							ands = append(ands, fmt.Sprintf("`%s` = ?", k))
							values = append(values, v)
						}
					}

					querys = append(querys, fmt.Sprintf("(%s)", strings.Join(ands, " AND ")))
				}


			}
		}

		orderStr := ""
		if len(orders) > 0 {
			orderStr = fmt.Sprintf("ORDER BY %s", strings.Join(orders, ","))
		}

		return strings.Join(querys, " OR "), values, orderStr, nil

	} else {
		return "1=1",[]interface{}{}, "",nil
	}
}
