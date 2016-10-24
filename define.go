/*
	预定义个啥，然后在到处可以方便的引用
	免的同样的定义写N遍，纯为了方便代码
	如你所见，此框架为了简化开发真是够了
*/

package noggo

import (
	. "github.com/nogio/noggo/base"
	"sync"
)

var (
	defines map[string]Any
	defineMutex sync.Mutex
)



//定义对象
func Define(key string, args ...Any) (Any) {
	defineMutex.Lock()
	defer defineMutex.Unlock()

	if defines == nil {
		defines = make(map[string]Any)
	}

	if len(args) > 0 {
		//set
		defines[key] = args[0]
		return args[0]
	} else {
		//get
		return defines[key]
	}

	return nil
}


//获取已定义的对象
//特别处理map对象，可以扩展或重写
func Defined(key string, extends ...Any) (Any) {
	defineMutex.Lock()
	defer defineMutex.Unlock()
	if defines == nil {
		defines = make(map[string]Any)
	}


	if vs,ok := defines[key].(Map); ok {
		m := Map{}
		for k,v := range vs {
			m[k] = v
		}

		//扩展
		if len(extends) > 0 {
			if vs,ok := extends[0].(Map); ok {
				for k,v := range vs {
					m[k] = v
				}
			}
		}

		return m
	} else {
		//非map对象暂不处理
		return defines[key]
	}
}



