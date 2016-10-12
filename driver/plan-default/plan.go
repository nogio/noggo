package plan_default


import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo"
	"github.com/nogio/noggo/driver/plan-default/cron"
	"errors"
)



//路由器定义
type (
	//驱动
	DefaultPlanDriver struct {
	}
	//连接
	DefaultPlanConnect struct {
		config Map
		cron *cron.Cron
		plans []string
	}
)


//返回驱动
func Driver() *DefaultPlanDriver {
	return &DefaultPlanDriver{}
}




//打开路由器
func (driver *DefaultPlanDriver) Connect(config Map) (noggo.PlanConnect) {
	//新建连接
	return &DefaultPlanConnect{
		config: config, plans: []string{},
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




//添加监听
func (driver *DefaultPlanConnect) Accept(name, time string, call func()) error {
	if driver.cron == nil {
		return errors.New("plan-default.accept: cron is nil")
	}

	//加入列表，记录所有name
	driver.plans = append(driver.plans, name)

	driver.cron.AddFunc(time, call, name)
	return nil
}





//移除监听
func (driver *DefaultPlanConnect) Remove(id string) error {
	if driver.cron == nil {
		return errors.New("plan-default.accept: cron is nil")
	}
	driver.cron.RemoveJob(id)
	return nil
}



//清理所有计划
func (driver *DefaultPlanConnect) Clear() error {
	if driver.cron == nil {
		return errors.New("plan-default.accept: cron is nil")
	}

	//移动所有
	for _,name := range driver.plans {
		driver.cron.RemoveJob(name)
	}

	return nil
}
