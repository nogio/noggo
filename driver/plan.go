package driver

import (
	. "github.com/nogio/noggo/base"
	"time"
)

type (
	//计划驱动
	PlanDriver interface {
		Connect(config Map) (PlanConnect,error)
	}

	PlanAccept func(*PlanRequest, PlanResponse)

	//计划连接
	PlanConnect interface {
		//打开连接
		Open() error
		//关闭连接
		Close() error

		//注册回调
		Accept(PlanAccept) error

		//创建计划
		Create(name string, time string) error
		//删除计划
		Remove(name string) error

		//开始计划
		Start() error
		//停止计划
		Stop() error
	}


	//计划请求实体
	PlanRequest struct {
		Id string
		Name string
		Time string
		Value Map
	}

	//计划响应接口
	PlanResponse interface {
		//完成
		Finish(id string) error
		//重新开始
		Replan(id string, delay time.Duration) error
	}
)
