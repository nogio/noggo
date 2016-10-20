package native

import (
	"github.com/nogio/noggo"
	"github.com/nogio/noggo/driver/logger-default"
	"github.com/nogio/noggo/driver/session-default"
	"github.com/nogio/noggo/driver/task-default"
	"github.com/nogio/noggo/driver/plan-default"
	"github.com/nogio/noggo/driver/event-default"
	"github.com/nogio/noggo/driver/http-default"
	"github.com/nogio/noggo/driver/view-default"
	"github.com/nogio/noggo/driver/queue-default"
)

func init() {
	noggo.Logger.Driver(noggo.ConstDriverDefault, logger_default.Driver())
	noggo.Session.Driver(noggo.ConstDriverDefault, session_default.Driver())
	noggo.Task.Driver(noggo.ConstDriverDefault, task_default.Driver())
	noggo.Plan.Driver(noggo.ConstDriverDefault, plan_default.Driver())
	noggo.Event.Driver(noggo.ConstDriverDefault, event_default.Driver())
	noggo.Queue.Driver(noggo.ConstDriverDefault, queue_default.Driver())
	noggo.Http.Driver(noggo.ConstDriverDefault, http_default.Driver())
	noggo.View.Driver(noggo.ConstDriverDefault, view_default.Driver())

}
