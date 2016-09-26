package driver

import (
	. "github.com/nogio/noggo/base"
)

type (

	//路由器结果
	RouterResult struct {
		Name	string
		Uri		string
		Params	Map
	}

	//路由器驱动
	RouterDriver interface {
		Connect(config Map) (RouterConnect)
	}
	//路由器连接
	RouterConnect interface {
		//打开连接
		Open()
		//关闭连接
		Close()
		//注册路由
		Route(name, uri string)
		//解析路由
		Parse(uri string) *RouterResult
	}
)
