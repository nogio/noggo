/*
	数据接口
	2016-10-22  定稿
*/

package driver

import (
	. "github.com/nogio/noggo/base"
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
	}
	//数据模型接口
	DataModel interface {
		Create(Map) (Map,error)
		Change(Map,Map) (Map,error)
		Remove(Map) (error)
		Entity(Any) (Map,error)

		//批量更新
		//其中args[0]为更新内容
		Update(args ...Map) (int64,error)
		Delete(args ...Map) (int64,error)

		Count(...Map) (int64,error)
		Single(...Map) (Map,error)
		Query(...Map) ([]Map,error)
		Limit(Any,Any,...Map) ([]Map,error)

		Group(string,...Map) ([]Map,error)
	}

)
