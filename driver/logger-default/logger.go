package logger_default


import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo/driver"
	"log"
	"time"
	"fmt"
)



//日志器定义
type (
	//驱动
	DefaultLogger struct {
	}
	DefaultConnect struct {
		config Map
	}
)


//返回驱动
func Driver() *DefaultLogger {
	return &DefaultLogger{}
}




//打开日志器
func (logger *DefaultLogger) Connect(config Map) (driver.LoggerConnect) {
	//新建日志器连接
	return &DefaultConnect{
		config: config,
	}
}














//打开连接
func (logger *DefaultConnect) Open() error {
	return nil
}



//关闭连接
func (logger *DefaultConnect) Close() error {
	return nil
}






//输出调试
func (logger *DefaultConnect) Debug(args ...interface{}) {
	log.Println(args...)
}



//输出信息
func (logger *DefaultConnect) Info(args ...interface{}) {
	newArgs := []interface{}{
		time.Now().Format("2006/01/02 15:04:05"),
	}
	newArgs = append(newArgs, args[:]...)

	fmt.Println(newArgs...)
}


//输出错误
func (logger *DefaultConnect) Error(args ...interface{}) {
	log.Println(args...)
}
