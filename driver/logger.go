/*
	日志接口
	2016-10-22  定稿
*/

package driver

import (
	"github.com/nogio/noggo"
	"github.com/nogio/noggo/driver/logger-default"
)


//默认任务驱动
func LoggerDefault() (noggo.LoggerDriver) {
	return logger_default.Driver()
}
