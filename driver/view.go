package driver

import (
	. "github.com/nogio/noggo/base"
	//"time"
)

type (
	//视图驱动
	ViewDriver interface {
		Connect(config Map) (ViewConnect)
	}
	//ViewAcceptFunc func(string,string,time.Duration,Map)

	ViewParse struct {
		Node    string  //当前节点
		Lang    string  //当前语言
		Data    Map     //ctx.Data
		View    string  //视图文件
		Model   Map     //视图模型
		Helpers Map     //工具方法
	}

	//视图连接
	ViewConnect interface {
		//打开连接
		Open() error
		//关闭连接
		Close() error

		//解析
		Parse(*ViewParse) (error,string)
	}
)
