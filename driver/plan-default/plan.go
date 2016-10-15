package plan_default


import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo/driver"
	"github.com/nogio/noggo/driver/plan-default/cron"
	"errors"
	"time"
	"sync"
)


//有BUG一枚~~~
//Create如果加了同步锁， 会一直被锁着， 要考虑处理一下



//路由器定义
type (
	//驱动
	DefaultPlanDriver struct {
	}
	//连接
	DefaultPlanConnect struct {
		config Map
		cron *cron.Cron
		callback    driver.PlanCallback

		//同步锁
		mutex   sync.Mutex
		//计划数据保存，三方驱动可以持久化
		datas   map[string]driver.PlanData
	}
)


//返回驱动
func Driver() *DefaultPlanDriver {
	return &DefaultPlanDriver{}
}




//打开路由器
func (ddddd *DefaultPlanDriver) Connect(config Map) (driver.PlanConnect) {
	//新建连接
	return &DefaultPlanConnect{
		config: config, datas: map[string]driver.PlanData{},
	}
}




//打开连接
func (driver *DefaultPlanConnect) Open() error {
	driver.cron = cron.New()
	return nil
}


//关闭连接
func (driver *DefaultPlanConnect) Close() error {
	return driver.Stop()
}






//注册回调
func (connect *DefaultPlanConnect) Accept(callback driver.PlanCallback) error {
	connect.callback = callback
	return nil
}






//开始
func (driver *DefaultPlanConnect) Start() error {
	if driver.cron == nil {
		return errors.New("plan-default.accept: cron is nil")
	}
	driver.cron.Start()
	return nil
}

//停止
func (driver *DefaultPlanConnect) Stop() error {
	if driver.cron != nil {
		return errors.New("plan-default.accept: cron is nil")
	}
	driver.cron.Stop()
	return nil
}




//创建计划
func (connect *DefaultPlanConnect) Create(name, time string) error {
	if connect.cron == nil {
		return errors.New("plan-default.accept: cron is nil")
	}

	connect.cron.AddFunc(time, func() {

		//新建计划
		id := NewMd5Id()
		plan := driver.PlanData{
			Name: name, Time: time, Value: Map{},
		}
		//保存计划
		connect.datas[id] = plan
		//调用计划
		connect.callback(id, plan.Name, plan.Time, plan.Value)

	}, name)

	return nil
}





//移除监听
func (connect *DefaultPlanConnect) Remove(name string) error {
	if connect.cron == nil {
		return errors.New("plan-default.accept: cron is nil")
	}
	connect.cron.RemoveJob(name)
	return nil
}







//完成任务，从列表中清理
func (connect *DefaultPlanConnect) Finish(id string) error {


	//从列表中删除
	//三方驱动， 应该做一些持久化的处理
	//更新掉此任务已经完成
	delete(connect.datas, id)

	return nil
}
//重开任务
func (connect *DefaultPlanConnect) Replan(id string, delay time.Duration) error {
	if plan,ok := connect.datas[id]; ok {

		time.AfterFunc(delay, func() {
			connect.callback(id, plan.Name, plan.Time, plan.Value)
		})
	}

	return errors.New("计划不存在")
}