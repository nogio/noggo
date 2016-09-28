package main

import (
	"github.com/nogio/noggo"
	_ "../drivers"
	_ "../middlers"
	_ "../modules/triggers"
	_ "../modules/plans"
)

func main() {
	noggo.Launch()
}
