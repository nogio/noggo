package driver

import (
	. "github.com/nogio/noggo/base"
)

type (
	//会话驱动
	SessionDriver interface {
		Connect(config Map) (SessionConnect,error)
	}
	//会话连接
	SessionConnect interface {
		//打开连接
		Open() error
		//关闭连接
		Close() error


		//查询会话，不存在就创建新的返回
		Entity(id string, expiry int64) (Map,error)
		//更新会话数据，不存在则创建，存在就更新
		Update(id string, value Map, expiry int64) error
		//删除会话
		Remove(id string) error
	}
)
