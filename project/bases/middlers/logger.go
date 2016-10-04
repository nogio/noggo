package middlers

import (
	"github.com/nogio/noggo"
	//"github.com/nogio/noggo/middler/plan-logger"
	"github.com/nogio/noggo/middler/http-logger"
)

func init() {

	noggo.Http.RequestFilter("logger", http_logger.Logger())

}
