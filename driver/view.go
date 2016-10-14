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
		//node      节点名称
		//view      视图文件（不带扩展名）
		//model     视图模型
		//data      viewdata
		Parse(node string, helpers Map, data Map, view string, model Map) (error,string)
	}
)