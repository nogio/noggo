/*
	队列接口
	2016-10-22  定稿
*/


package driver

import (
	"github.com/nogio/noggo"
	"github.com/nogio/noggo/driver/queue-default"
)

//默认队列驱动
func QueueDefault() (noggo.QueueDriver) {
	return queue_default.Driver()
}
