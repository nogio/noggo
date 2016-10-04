package router_default


import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo"
	"sync"
	"regexp"
	"strings"
)



//路由器定义
type (
	//驱动
	DefaultRouter struct {
	}
	DefaultRoute struct {
		name	string
		uri		string
	}
	DefaultConnect struct {
		config Map
		routes []DefaultRoute	//保证顺序
		routesMutex sync.Mutex
	}
)


//返回驱动
func Driver() *DefaultRouter {
	return &DefaultRouter{}
}




//打开路由器
func (router *DefaultRouter) Connect(config Map) (noggo.RouterConnect) {
	//新建路由器连接
	return &DefaultConnect{
		config: config, routes: []DefaultRoute{},
	}
}














//打开连接
func (router *DefaultConnect) Open() error {
	return nil
}



//关闭连接
func (router *DefaultConnect) Close() {

}




//注册路由
func (router *DefaultConnect) Accept(name, uri string) {
	router.routes = append(router.routes, DefaultRoute{name,uri})
}



//解析路由
func (router *DefaultConnect) Parse(host,path string) (*noggo.RouterResult){

	//路由解析
	for _,obj := range router.routes {
		uri := obj.uri

		keys, vals := []string{}, []string{}

		//用来匹配的正则表达式
		regs := strings.Replace(uri, ".", "\\.", -1)
		regs = strings.Replace(regs, "/", "\\/", -1)


		//正则对象
		regx := regexp.MustCompile(`\{[\*A-Za-z0-9_]+\}`)
		//拿到URI中的参数列表
		//keys := regx.FindAllString(uri, -1)
		//替换参数为正则
		regs = regx.ReplaceAllStringFunc(regs, func(p string) string {
			if (p[1:1] == "*") {
				keys = append(keys, p[2:len(p)-1])
				return `(.*)`
			} else {
				keys = append(keys, p[1:len(p)-1])
				return `([-_A-Za-z0-9\.]+)`
			}
		})



		//正则匹配当前URL，是否从域名开始匹配
		//忽略query部分
		url := path
		if (uri[0:1] != "/") {
			url = host + path
		}


		regx = regexp.MustCompile("^"+regs+"$")
		if regx.MatchString(url) {

			matchs := regx.FindAllStringSubmatch(url, -1)
			if matchs != nil && len(matchs) > 0 && len(matchs[0]) > 0 {
				vals = matchs[0][1:]
			}

			param := Map{}

			//处理params
			if len(keys) == len(vals) {
				for i,k := range keys {
					param[k] = vals[i]
				}
			}

			return &noggo.RouterResult{
				Name: obj.name, Uri: obj.uri, Param: param,
			}
		}
	}



	return nil
}
