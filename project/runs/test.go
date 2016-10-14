package main

import (
	"github.com/nogio/noggo"
)

func main() {
	nog := noggo.New()

	nog.Http.Route("/index.html", func(ctx *noggo.HttpContext) {
		ctx.Text("hahahahaha")
	})

	nog.Run()
}