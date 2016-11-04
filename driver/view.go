/*
	视图接口
	2016-10-22  定稿
*/

package driver

import (
	"github.com/nogio/noggo"
	"github.com/nogio/noggo/driver/view-default"
)



//默认视图引擎
func ViewDefault(roots ...string) (noggo.ViewDriver) {
	return view_default.Driver(roots...)
}
