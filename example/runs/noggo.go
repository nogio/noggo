package main

import (
	"github.com/nogio/noggo"
	_ "../contexts"
	_ "../drivers"
	_ "../modules/triggers"
	_ "../modules/plans"
)

func main() {
	noggo.Launch()
}
