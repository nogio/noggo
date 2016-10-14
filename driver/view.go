package driver

import (
	. "github.com/nogio/noggo/base"
	"time"
)

type (
	//视图驱动
	ViewDriver interface {
		Connect(config Map) (ViewConnect)
	}
	ViewAcceptFunc func(string,string,time.Duration,Map)
	//视图连接
	ViewConnect interface {
		//打开连接
		Open() error
		//关闭连接
		Close() error

		//帮助VIEW函数
		Helper(name string, helper Any) error

		//解析
		Parse(view string, model Map, data Map) (error,string)
	}
)
