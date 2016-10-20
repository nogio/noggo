package noggo

import (
	"os"
	"os/signal"
	"syscall"
)

var (
	initialized bool
)

func Init() {
	if initialized == false {
		Logger.init()
		Session.init()

		Trigger.init()
		Task.init()

		Plan.init()
		Event.init()
		Queue.init()
		Http.init()
		View.init()

		Data.init()

		Logger.Info("noggo is initiating...")

		initialized = true
	}
}


func Exit() {
	wating()
	Logger.Info("noggo is exiting...")

	//退出所有节点
	for _,node := range nodes {
		node.End()
	}

	Data.exit()

	View.exit()
	Http.exit()
	Queue.exit()
	Event.exit()
	Plan.exit()
	Task.exit()
	Trigger.exit()
	Session.exit()
	Logger.exit()
}


//使用管道监听退出信号
func wating() {
	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM)
	<-exitChan
}
