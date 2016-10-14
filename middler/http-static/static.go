package http_static

import (
	"github.com/nogio/noggo"
	"fmt"
	"os"
)



//返回中间件
func Middler(paths ...string) (noggo.HttpFunc) {

	path := "statics"
	if len(paths) > 0 {
		path = paths[0]
	}

	return func(ctx *noggo.HttpContext) {
		var file string

		//先搜索节点所在目录
		// statics/node/xxx
		file = fmt.Sprintf("%s/%s%s", path, ctx.Node.Name, ctx.Path)
		fi, _ := os.Stat(file)
		if fi != nil && !fi.IsDir() {
			ctx.File(file)
		} else {

			//全局静态目录
			// statics/default
			file = fmt.Sprintf("%s/default%s", path, ctx.Path)
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
