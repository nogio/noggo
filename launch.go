package noggo

import (
	"os"
	"os/signal"
	"syscall"
)

func Init() {
	Logger.init()
	Session.init()
	Trigger.init()
	Task.init()
	Logger.Info("noggo is running...")
}


func Exit() {
	wating()
	Logger.Info("noggo is exiting...")
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
