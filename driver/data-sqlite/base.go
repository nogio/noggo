package data_sqlite


import (
	"github.com/nogio/noggo"
	. "github.com/nogio/noggo/base"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"encoding/json"
	"strconv"
)

type (
	SqliteBase struct {
		name    string
		conn    *SqliteConnect
		models  map[string]Map
		views   map[string]Map

		db      *sql.DB
		tx      *sql.Tx
		cache   noggo.CacheBase
		caching bool

		//是否手动提交事务，否则为自动
		//当调用begin时， 自动变成手动提交事务
		manual    bool
	}
)


//关闭数据库
func (base *SqliteBase) Close() {

	//好像目前不需要关闭什么东西
	if base.tx != nil {
		//关闭时候,一定要提交一次事务
		//如果手动提交了, 这里会失败, 问题不大
		//如果没有提交的话, 连接不会交回连接池. 会一直占用
		base.Cancel()
	}
}

//获取模型对象
func (base *SqliteBase) Model(name string) (noggo.DataModel) {
	if config,ok := base.models[name]; ok {

		//模式，表名
		//模式默认等于库名，而不是public，万恶的Sqlite
		schema, object, key, fields := base.name, name, "id", Map{}
		if n,ok := config["schema"].(string); ok {
			schema = n
		}

		if n,ok := config["object"].(string); ok {
			object = n
		}
		if n,ok := config["table"].(string); ok {
			object = n
		}

		if n,ok := config["key"].(string); ok {
			key = n
		}
		if n,ok := config["fields"].(Map); ok {
			fields = n
		} else {
			panic("数据：未定义字段fields")
		}

		return &SqliteModel{
			SqliteView{base, name, schema, object, key, fields},
		}
	} else {
		panic("数据：模型不存在")
	}
}


//获取模型对象
func (base *SqliteBase) View(name string) (noggo.DataView) {
	if config,ok := base.views[name]; ok {

		//模式，表名
		//模式默认等于库名，而不是public，万恶的Sqlite
		schema, object, key, fields := base.name, name, "id", Map{}
		if n,ok := config["schema"].(string); ok {
			schema = n
		}

		if n,ok := config["object"].(string); ok {
			object = n
		}
		if n,ok := config["view"].(string); ok {
			object = n
		}

		if n,ok := config["key"].(string); ok {
			key = n
		}
		if n,ok := config["fields"].(Map); ok {
			fields = n
		} else {
			panic("数据：未定义字段fields")
		}

		return &SqliteView{
			base, name, schema, object, key, fields,
		}
	} else {
		panic("数据：视图不存在")
	}
}

//是否开启缓存
func (base *SqliteBase) Cache(use bool) (noggo.DataBase) {
	base.caching = use
	return base
}


//开启手动模式
func (base *SqliteBase) Begin() (noggo.DataBase) {
	base.manual = true
	return base
}



//注意，此方法为实际开始事务
func (base *SqliteBase) begin() (*sql.Tx,error) {

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
func (base *SqliteBase) Submit() (error) {

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
func (base *SqliteBase) Cancel() (error) {

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
func (base *SqliteBase) Exec(query string, args ...interface{}) (sql.Result,error) {
	if base.tx == nil {
		return nil,errors.New("数据：tx未开始")
	}
	return base.tx.Exec(query, args...)
}

//Prepare
func (base *SqliteBase) Prepare(query string) (*sql.Stmt, error) {
	if base.tx == nil {
		return nil,errors.New("数据：tx未开始")
	}
	return base.tx.Prepare(query)
}

//Query
func (base *SqliteBase) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if base.tx == nil {
		return nil,errors.New("数据：tx未开始")
	}
	return base.tx.Query(query, args...)
}

//QueryRow
func (base *SqliteBase) QueryRow(query string, args ...interface{}) (*sql.Row) {
	if base.tx == nil {
		return nil
	}
	return base.tx.QueryRow(query, args...)
}

//QueryRow
func (base *SqliteBase) Stmt(stmt *sql.Stmt) (*sql.Stmt) {
	if base.tx == nil {
		return nil
	}
	return base.tx.Stmt(stmt)
}




















//创建的时候,也需要对值来处理,
//数组要转成{a,b,c}格式,要不然不支持
//json可能要转成字串才支持
func (base *SqliteBase) packing(value Map) (Map) {

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
		case []Map: {
			b,e := json.Marshal(t);
			if e == nil {
				newValue[k] = string(b)
			} else {
				newValue[k] = "[]"
			}
		}
		case MapList: {
			b,e := json.Marshal(t);
			if e == nil {
				newValue[k] = string(b)
			} else {
				newValue[k] = "[]"
			}
		}
		case []map[string]interface{}: {
			b,e := json.Marshal(t);
			if e == nil {
				newValue[k] = string(b)
			} else {
				newValue[k] = "[]"
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
func (base *SqliteBase) parsing(args ...Any) (string,[]interface{},string,error) {

	sql,val,odr,err := noggo.Data.Parsing(args...)
	if err == nil {
		//结果要处理一下，字段包裹、参数处理
		sql = strings.Replace(sql, noggo.DataFieldDelims, "`", -1)
		odr = strings.Replace(odr, noggo.DataFieldDelims, "`", -1)
	}

	return sql,val,odr,err
}
