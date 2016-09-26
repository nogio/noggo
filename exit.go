package noggo

import (
	"os"
	"os/signal"
	"syscall"
	"fmt"
)

func exitWaiting() {
	//使用管道监听退出信号
	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM)
	<-exitChan
}

func Exit() {
	exitWaiting()

	//退出前处理
	fmt.Println("noggo exit...")
}
