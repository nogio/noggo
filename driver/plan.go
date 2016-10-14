package driver

import (
	. "github.com/nogio/noggo/base"
)

type (
	//计划驱动
	PlanDriver interface {
		Connect(config Map) (PlanConnect)
	}
	//计划连接
	PlanConnect interface {
		//打开连接
		Open() error
		//关闭连接
		Close() error

		//注册计划
		Accept(name,time string, call func()) error
		//删除计划
		Remove(id string) error
		//清空计划
		Clear() error

		//开始计划
		Start() error
		//停止计划
		Stop() error
	}
)
