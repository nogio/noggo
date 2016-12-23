package tempdb

import (
	"github.com/nogio/noggo"
)

func init() {
	//标记当前目录的模型注册到temp库下
	noggo.Current = "temp"
}
