package noggo

import (
	"os"
	"os/signal"
	"syscall"
	"fmt"
)



func Launch() {
	nogInit()
	nogExit()
}




//初始化
func nogInit() {
	fmt.Println("nogo init...")

	Logger.init()
	Router.init()
	Session.init()
	Trigger.init()
	Plan.init()
	Http.init()


}

func nogExit() {

	//使用管道监听退出信号
	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM)
	<-exitChan

	Logger.exit()
	Router.exit()
	Session.exit()
	Trigger.exit()
	Plan.exit()
	Http.exit()
	fmt.Println("nogo exit...")
}
