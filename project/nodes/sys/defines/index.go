package defines

import (
	"github.com/nogio/noggo"
	. "github.com/nogio/noggo/base"
)

func init() {

	//预定义后台管理员auth节点
	noggo.Define("sys.auth", Map{
		"admin": Map{},
	})
	noggo.Define("sys.auth.not", Map{
		"admin": Map{},
	})

}