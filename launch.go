package noggo

import (
	"os"
	"os/signal"
	"syscall"
	"fmt"
)



func Launch() {
	Init()
	Exit()
}




//初始化
func Init() {
	fmt.Println("nogo init...")

	Logger.init()
	Router.init()
	Session.init()
	Trigger.init()
	Plan.init()


}

func Exit() {

	//使用管道监听退出信号
	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM)
	<-exitChan

	Logger.exit()
	Router.exit()
	Session.exit()
	Trigger.exit()
	Plan.exit()
	fmt.Println("nogo exit...")
}
