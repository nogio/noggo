/*
	缓存接口
	2016-10-22  未定稿
	因为返回的时候，其实无法确定是什么任务，也无法强制转换
	所有只能在使用缓存的地方，手动设置和转换类型
	所以暂时不能确认接口文档
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

	CacheValueFunc func(Any)(Any)


	//缓存库
	CacheBase interface {
		//查询缓存，自带值包装函数
		Get(key string) (Any,error)
		//更新缓存数据，不存在则创建，存在就更新
		Set(key string, val Any, exp int64) error
		//删除缓存
		Del(key string) error

		//获取keys
		Keys(prefixs ...string) ([]string,error)

		//清空
		Empty(prefixs ...string) (error)
	}
)
