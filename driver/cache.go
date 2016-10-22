/*
	缓存接口
	2016-10-22  未定稿
*/


package driver

import (
	. "github.com/nogio/noggo/base"
)

type (
	//缓存驱动
	CacheDriver interface {
		Connect(config Map) (CacheConnect,error)
	}
	//缓存连接
	CacheConnect interface {
		//打开连接
		Open() error
		//关闭连接
		Close() error

		//获取数据库对象
		Base(string) (CacheBase,error)
	}

	//缓存库
	CacheBase interface {
		//查询缓存，不存在就创建新的返回
		Get(key string) (Any,error)
		//更新缓存数据，不存在则创建，存在就更新
		Set(key string, val Any, exp int64) error
		//删除缓存
		Del(key string) error

		//获取keys
		Keys(args ...string) ([]string,error)

		//清空
		Empty() (error)
	}
)
