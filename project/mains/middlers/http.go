package middlers

import (
    "github.com/nogio/noggo"
    "github.com/nogio/noggo/middler"
)

func init() {

    //注册中间件
    noggo.Middler(middler.HttpLogger())
    noggo.Middler(middler.HttpStatic(noggo.Config.Custom.String("static")))
    noggo.Middler(middler.HttpForm(noggo.Config.Custom.String("upload")))

}
