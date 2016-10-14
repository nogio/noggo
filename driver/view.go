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

		//解析
		Parse(node string, helpers Map, data Map, view string, model Map) (error,string)
	}
)
