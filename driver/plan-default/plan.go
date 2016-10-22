package plan_default


import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo/driver"
	"github.com/nogio/noggo/driver/plan-default/cron"
	"errors"
	"time"
	"sync"
)


type (
	//驱动
	DefaultPlanDriver struct {
	}
	//连接
	DefaultPlanConnect struct {
		mutex   sync.Mutex

		config Map
		handler    driver.PlanHandler

		cron *cron.Cron

		//计划数据保存，三方驱动可以持久化
		plans   map[string]string   //map[name]time
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
func (drv *DefaultPlanDriver) Connect(config Map) (driver.PlanConnect,error) {
	return &DefaultPlanConnect{
		config: config, plans: map[string]string{}, datas: map[string]PlanData{},
	}, nil
}




//打开连接
func (drv *DefaultPlanConnect) Open() error {
	drv.cron = cron.New()
	return nil
}
//关闭连接
func (drv *DefaultPlanConnect) Close() error {
	drv.cron.Stop()
	return nil
}




//注册回调
func (con *DefaultPlanConnect) Accept(handler driver.PlanHandler) error {
	con.handler = handler
	return nil
}
//注册计划
func (con *DefaultPlanConnect) Register(name, time string) error {
	con.mutex.Lock()
	defer con.mutex.Unlock()

	//保存计划列表
	con.plans[name] = time

	return nil
}
//开始计划
func (con *DefaultPlanConnect) Start() error {

	for name,time := range con.plans {
		con.cron.AddFunc(time, func() {
			//新建计划
			id := NewMd5Id()
			plan := PlanData{
				Name: name, Time: time, Value: Map{},
			}
			//保存计划
			con.mutex.Lock()
			con.datas[id] = plan
			con.mutex.Unlock()
			//调用计划
			con.execute(id, plan.Name, plan.Time, plan.Value)

		}, name)

	}
	con.cron.Start()
	return nil
}



//执行统一到这里
func (connect *DefaultPlanConnect) execute(id string, name string, time string, value Map) {
	req := &driver.PlanRequest{ Id: id, Name: name, Time: time, Value: value }
	res := &DefaultPlanResponse{ connect }
	connect.handler(req, res)
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