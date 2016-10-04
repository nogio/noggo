package main

import (
	"github.com/nogio/noggo"

	_ "../bases/consts"
	_ "../bases/types"
	_ "../bases/cryptos"
	_ "../bases/drivers"
	_ "../bases/middlers"

	_ "../modules/triggers"
	_ "../modules/plans"
	_ "../modules/tasks"
	_ "../modules/https"
)

func main() {
	noggo.Launch()
}
