package driver

import (
	. "github.com/nogio/noggo/base"
	"time"
)

type (
	//任务驱动
	TaskDriver interface {
		Connect(config Map) (error,TaskConnect)
	}

	//回调函数
	TaskAccept func(*TaskRequest, TaskResponse)

	//任务连接
	TaskConnect interface {
		//打开连接
		Open() error
		//关闭连接
		Close() error


		//注册回调
		Accept(TaskAccept) error


		//开始任务
		Start() error
		//停止任务
		Stop() error


		//触发任务
		After(name string, delay time.Duration, value Map) error
	}

	//任务请求实体
	TaskRequest struct {
		Id string
		Name string
		Delay time.Duration
		Value Map
	}

	//任务响应接口
	TaskResponse interface {
		//完成任务
		Finish(id string) error
		//重新开始任务
		Retask(id string, delay time.Duration) error
	}
)
