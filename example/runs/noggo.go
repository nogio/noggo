package main

import (
	"github.com/nogio/noggo"
	_ "../modules/triggers"	//加载触发器定义
	"fmt"
)

func main() {

	fmt.Printf("config=%v\n", noggo.Config)

	noggo.Run()
}
