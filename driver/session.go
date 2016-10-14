package driver

import (
	. "github.com/nogio/noggo/base"
)

type (
	//会话驱动
	SessionDriver interface {
		Connect(config Map) (SessionConnect)
	}
	//会话连接
	SessionConnect interface {
		//打开连接
		Open() error
		//关闭连接
		Close() error


		//查询会话，不存在就创建新的返回
		Query(id string, expiry int64) (error,Map)
		//更新会话数据，不存在则创建，存在就更新
		Update(id string, value Map, expiry int64) error
		//删除会话
		Remove(id string) error
		//回收会话，系统会每一段时间自动调用此方法
		Recycle(expiry int64) error
	}
)
