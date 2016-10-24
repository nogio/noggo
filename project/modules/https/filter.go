package https

import (
	"github.com/nogio/noggo"
	"github.com/nogio/noggo/middler"
)

func init() {
	//三基础中间件
	noggo.Use(middler.HttpLogger())
	noggo.Use(middler.HttpStatic("statics"))
	noggo.Use(middler.HttpForm("uploads"))
}