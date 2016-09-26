package noggo

import (
	"os"
	"os/signal"
	"syscall"
	"fmt"
)



func Run() {
	Init()
	Exit()
}

func End() {
	fmt.Println("nogo exit...")
}




//初始化
func Init() {
	fmt.Println("nogo init...")
}









func exitWaiting() {
	//使用管道监听退出信号
	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM)
	<-exitChan
}

func Exit() {
	exitWaiting()
	End()
}
