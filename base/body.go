package base

import "time"


type (
	//触发器响应完成
	BodyTriggerFinish struct {}
	//触发器响应重来
	BodyTriggerRetrigger struct {
		Delay	time.Duration
	}
)



type (
	//计划响应完成
	BodyPlanFinish struct {}
	//计划响应重来
	BodyPlanReplan struct {
		Delay	time.Duration
	}
)




type (
	//任务响应完成
	BodyTaskFinish struct {}
	//任务响应重来
	BodyTaskRetask struct {
		Delay	time.Duration
	}
)








