package driver

import (
	. "github.com/nogio/noggo/base"
	"database/sql"
)

type (
	DataDriver interface {
		//连接驱动的时候
		//应该做如下工作：
		//1. 检查config合法性
		//2. 初始化连接相关对象
		Connect(config Map) (error,DataConnect)
	}
	DataConnect interface {
		//打开数据库连接
		Open() (error)
		//关闭数据库连接
		Close() (error)
		//获取数据库对象
		DB(string) (error,DataDB)
	}

	//数据库接口
	DataDB interface {
		Model(name string) (DataModel)
		Close()
		Begin() (error,*sql.Tx)
		Submit() (error)
		Cancel() (error)
	}
	//数据模型接口
	DataModel interface {
		Create(Map) (error,Map)
		Change(Map,Map) (error,Map)
		Remove(Map) (error)
		Entity(Any) (error,Map)
	}


)
