package main

import (
	"github.com/nogio/noggo"
	_ "../drivers"
	_ "../middlers"
	_ "../modules/triggers"
	_ "../modules/plans"
	_ "../modules/tasks"
	_ "../modules/https"
)

func main() {
	noggo.Launch()
}
