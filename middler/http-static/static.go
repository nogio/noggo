package http_static

import (
	"github.com/nogio/noggo"
	"fmt"
	"os"
)



//返回中间件
func Middler() (noggo.HttpFunc) {
	return func(ctx *noggo.HttpContext) {
		var file string

		//先搜索节点所在目录
		// statics/node/xxx
		file = fmt.Sprintf("statics/%s%s", ctx.Node.Name, ctx.Path)
		fi, _ := os.Stat(file)
		if fi != nil && !fi.IsDir() {
			ctx.File(file)
		} else {

			//全局静态目录
			// statics/default
			file = fmt.Sprintf("statics/default%s", ctx.Path)
			fi, _ := os.Stat(file)
			if fi != nil && !fi.IsDir() {
				ctx.File(file)
			} else {
				//都不存在， 继续走
				ctx.Next()
			}

		}
	}
}
