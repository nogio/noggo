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



type (
	//驱动
	DefaultPlanDriver struct {
	}
	//连接
	DefaultPlanConnect struct {
		config Map
		cron *cron.Cron
		callback    driver.PlanAccept

		//同步锁
		mutex   sync.Mutex
		//计划数据保存，三方驱动可以持久化
		datas   map[string]PlanData
	}

	//响应对象
	DefaultPlanResponse struct {
		connect *DefaultPlanConnect
	}


	PlanData struct{
		Name    string
		Time    string
		Value   Map
	}
)


//返回驱动
func Driver() (driver.PlanDriver) {
	return &DefaultPlanDriver{}
}




//连接
func (drv *DefaultPlanDriver) Connect(config Map) (error,driver.PlanConnect) {
	return nil,&DefaultPlanConnect{
		config: config, datas: map[string]PlanData{},
	}
}




//打开连接
func (drv *DefaultPlanConnect) Open() error {
	drv.cron = cron.New()
	return nil
}
//关闭连接
func (drv *DefaultPlanConnect) Close() error {
	return drv.Stop()
}



//注册回调
func (connect *DefaultPlanConnect) Accept(callback driver.PlanAccept) error {
	connect.callback = callback
	return nil
}




//创建计划
//直接就在这addfunc了。
//应该在这里保存所有计划
//然后在Start中运行， 这样科学一些
func (connect *DefaultPlanConnect) Create(name, time string) error {
	if connect.cron == nil {
		return errors.New("plan-default.accept: cron is nil")
	}

	connect.cron.AddFunc(time, func() {

		//新建计划
		id := NewMd5Id()
		plan := PlanData{
			Name: name, Time: time, Value: Map{},
		}
		//保存计划
		connect.datas[id] = plan
		//调用计划
		connect.execute(id, plan.Name, plan.Time, plan.Value)

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



//执行统一到这里
func (connect *DefaultPlanConnect) execute(id string, name string, time string, value Map) {
	req := &driver.PlanRequest{ Id: id, Name: name, Time: time, Value: value }
	res := &DefaultPlanResponse{ connect }
	connect.callback(req, res)
}














//完成任务，从列表中清理
func (connect *DefaultPlanConnect) finish(id string) error {
	connect.mutex.Lock()
	delete(connect.datas, id)
	connect.mutex.Unlock()
	return nil
}
//重开任务
func (connect *DefaultPlanConnect) replan(id string, delay time.Duration) error {
	if plan,ok := connect.datas[id]; ok {
		time.AfterFunc(delay, func() {
			connect.execute(id, plan.Name, plan.Time, plan.Value)
		})
	}

	return errors.New("计划不存在")
}











//响应接口对象



//完成任务，从列表中清理
func (res *DefaultPlanResponse) Finish(id string) error {
	return res.connect.finish(id)
}
//重开任务
func (res *DefaultPlanResponse) Replan(id string, delay time.Duration) error {
	return res.connect.replan(id, delay)
}