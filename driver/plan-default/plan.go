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
	DefaultPlan struct {
	}
	//连接
	DefaultConnect struct {
		config Map
		cron *cron.Cron
	}
)


//返回驱动
func Driver() *DefaultPlan {
	return &DefaultPlan{}
}




//打开路由器
func (plan *DefaultPlan) Connect(config Map) (noggo.PlanConnect) {
	//新建路由器连接
	return &DefaultConnect{
		config: config,
	}
}




//打开连接
func (plan *DefaultConnect) Open() error {
	plan.cron = cron.New()
	return nil
}


//关闭连接
func (plan *DefaultConnect) Close() error {
	return plan.Stop()
}










//开始
func (plan *DefaultConnect) Start() error {
	if plan.cron == nil {
		return errors.New("plan-default.accept: cron is nil")
	}
	plan.cron.Start()
	return nil
}

//停止
func (plan *DefaultConnect) Stop() error {
	if plan.cron != nil {
		return errors.New("plan-default.accept: cron is nil")
	}
	plan.cron.Stop()
	return nil
}




//添加监听
func (plan *DefaultConnect) Accept(job noggo.PlanJob) error {
	if plan.cron == nil {
		return errors.New("plan-default.accept: cron is nil")
	}
	plan.cron.AddFunc(job.Time, job.Call, job.Id)
	return nil
}





//移除监听
func (plan *DefaultConnect) Remove(id string) error {
	if plan.cron == nil {
		return errors.New("plan-default.accept: cron is nil")
	}
	plan.cron.RemoveJob(id)
	return nil
}
