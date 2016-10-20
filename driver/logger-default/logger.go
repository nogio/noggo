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
	DefaultLoggerDriver struct {
	}
	DefaultConnect struct {
		config Map
	}
)


//返回驱动
func Driver() (driver.LoggerDriver) {
	return &DefaultLoggerDriver{}
}




//连接驱动
func (logger *DefaultLoggerDriver) Connect(config Map) (driver.LoggerConnect,error) {
	return &DefaultConnect{
		config: config,
	},nil
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
