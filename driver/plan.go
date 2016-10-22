/*
	计划接口
	2016-10-22  定稿
*/

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

	PlanHandler func(*PlanRequest, PlanResponse)

	//计划连接
	PlanConnect interface {
		//打开连接
		Open() error
		//关闭连接
		Close() error

		//注册回调
		Accept(PlanHandler) error

		//注册计划
		Register(name string, time string) error

		//开始计划
		Start() error
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
