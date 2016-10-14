package driver

import (
	. "github.com/nogio/noggo/base"
	"time"
)

type (
	//任务驱动
	TaskDriver interface {
		Connect(config Map) (TaskConnect)
	}
	TaskAcceptFunc func(string,string,time.Duration,Map)
	//任务连接
	TaskConnect interface {
		//打开连接
		Open() error
		//关闭连接
		Close() error

		//注册任务
		Accept(name string, call TaskAcceptFunc) error


		//打开连接
		Start() error
		//关闭连接
		Stop() error


		//触发任务
		Touch(id string, name string, delay time.Duration, value Map) error

		//完成任务
		Finish(id string) error
	}
)
