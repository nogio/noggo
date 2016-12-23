package defines


/*
    预定义一些对象，可以在项目里通用
*/


import (
    "github.com/nogio/noggo"
    . "github.com/nogio/noggo/base"
)

func init() {

    noggo.Define("auth.user", Map{
        "aaa": "aaa",
    })

}