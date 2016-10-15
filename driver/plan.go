package driver

import (
	. "github.com/nogio/noggo/base"
	"time"
)

type (
	//计划驱动
	PlanDriver interface {
		Connect(config Map) (PlanConnect)
	}

	PlanData struct {
		Name    string
		Time    string
		Value   Map
	}

	PlanCallback func(string,string,string,Map)

	//计划连接
	PlanConnect interface {
		//打开连接
		Open() error
		//关闭连接
		Close() error

		//注册回调
		Accept(callback PlanCallback) error

		//创建计划
		Create(name string, time string) error
		//删除计划
		Remove(name string) error

		//开始计划
		Start() error
		//停止计划
		Stop() error


		//完成计划
		Finish(id string) error
		//重开计划
		Replan(id string, delay time.Duration) error
	}
)
