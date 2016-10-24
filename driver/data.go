/*
	数据接口
	2016-10-22  定稿
*/

package driver

import (
	. "github.com/nogio/noggo/base"
	"database/sql"
)

const (
	ASC     = "$$$ASC$$$"
	DESC    = "$$$DESC$$$"
	COUNT   = "$$$COUNT$$$"
	AVG     = "$$$AVG$$$"
	SUM     = "$$$SUM$$$"
	MAX     = "$$$MAX$$$"
	MIN     = "$$$MIN$$$"
)

type (
	DataDriver interface {
		//连接驱动的时候
		//应该做如下工作：
		//1. 检查config合法性
		//2. 初始化连接相关对象
		Connect(config Map) (DataConnect,error)
	}
	DataConnect interface {
		//打开数据库连接
		Open() (error)
		//关闭数据库连接
		Close() (error)
		//注册模型
		Model(string,Map)
		//获取数据库对象
		Base(string,CacheBase) (DataBase,error)
	}

	//数据库接口
	DataBase interface {
		Close()
		Model(name string) (DataModel)


		//开启手动提交事务模式
		Begin() (DataBase)
		Submit() (error)
		Cancel() (error)

		//原生SQL的方法，接口可以执行原生查询，以支持Model不能完成的工作
		//但是，必须 调用Begin之后，才能使用下例方法，然后 Submit 或 Cancel
		//因为全部使用事务。
		Exec(query string, args ...interface{}) (sql.Result, error)
		Prepare(query string) (*sql.Stmt, error)
		Query(query string, args ...interface{}) (*sql.Rows, error)
		QueryRow(query string, args ...interface{}) (*sql.Row)
		Stmt(stmt *sql.Stmt) (*sql.Stmt)
}
	//数据模型接口
	DataModel interface {
		Create(Map) (Map,error)
		Change(Map,Map) (Map,error)
		Remove(Map) (error)
		Entity(Any) (Map,error)

		//批量更新
		Update(sets Map,args ...Map) (int64,error)
		//批量删除
		Delete(args ...Map) (int64,error)

		Count(...Map) (int64,error)
		Single(...Map) (Map,error)
		Query(...Map) ([]Map,error)
		Limit(Any,Any,...Map) ([]Map,error)

		Group(string,...Map) ([]Map,error)
	}

)
