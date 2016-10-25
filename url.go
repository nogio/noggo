package noggo

import (
	. "github.com/nogio/noggo/base"
	"strings"
	"fmt"
	"regexp"
)

const (
	BACKURL = "_back_"
)

type (
	httpUrl struct {
		ctx *HttpContext
	}
)





func (url *httpUrl) Node(name string, args ...string) string {
	if c,ok := Config.Node[name]; ok {
		if len(args) > 0 {
			return fmt.Sprintf("%s%s", c.Url, args[0])
		} else {
			return c.Url
		}
	}
	return "[no node here]"
}



//当前站本身的路由
func (url *httpUrl) Route(name string, args ...Map) string {

	if url.ctx != nil {
		if name == "" {
			name = url.ctx.Name
		}
	}

	config := Map{}
	uri := name //支持外部URL加上参数
	if config,ok := url.ctx.Node.Http.routes[name]; ok {
		if c,ok := config["uri"]; ok {
			switch v := c.(type) {
			case string:
				uri = v
			case []string:
				uri = v[0]
			}
		}
	}



	if uri != "" {
		argsConfig := Map{}
		if c,ok := config["args"]; ok {
			argsConfig = c.(Map)
		} else if c,ok := config["route"]; ok {
			routeConfig := c.(Map)
			if d,ok := routeConfig["args"]; ok {
				argsConfig = d.(Map)
			} else if d,ok := routeConfig["get"]; ok {
				methodConfig := d.(Map)
				if e,ok := methodConfig["args"]; ok {
					argsConfig = e.(Map)
				}
			}
		} else {
			//其它无参数
		}




		//得到值和选项
		datas,values,options := Map{}, Map{}, Map{}
		if len(args) > 0 {
			if len(args) == 1 {
				values = args[0]
			} else {
				values = args[0]
				options = args[1]
			}
		}



		//选项处理
		if (options["back"] != nil && url.ctx != nil) {
			var url = url.Back()
			datas[BACKURL] = encode64(url);
		}
		//选项处理
		if (options["last"] != nil && url.ctx != nil) {
			var url = url.Last()
			datas[BACKURL] = encode64(url);
		}
		//自动携带原有的query信息
		if (options["query"] != nil && url.ctx != nil) {
			for k,v := range url.ctx.Query {
				datas[k] = v
			}
		}





		//再从args,是不是有默认值
		argsValue := Map{}
		e := Mapping.Parse([]string{}, argsConfig, Map{}, argsValue)
		//不直接写datas,而是在下面的params中,如果有,才写入
		if e == nil {
			/*
			for k,v := range argsValue {
				datas["{"+k+"}"] = v
			}
			*/
		}

		//这里后写, 才会盖掉args中的值
		//先从web.Params拿值,这样如果是同一页,就有的相同的参数了
		//当前页相同的路由,才加入Params参数
		//if url.ctx != nil && name == url.ctx.Name {
		//不能直接复制
		/*
		if url.ctx != nil {
			for k,v := range url.ctx.Params {
				datas["{"+k+"}"] = v
			}
		}
		*/

		webParams := Map{}
		if url.ctx != nil {
			webParams = url.ctx.Param
		}

		//然后是传过来的值
		for k,v := range values {
			datas[k] = v
		}




		//keys := []string{}

		//直接用正则替换
		regx := regexp.MustCompile(`\{[_\*A-Za-z0-9]+\}`)
		uri = regx.ReplaceAllStringFunc(uri, func(p string) string {
			key := strings.Replace(p, "*", "", -1)
			argsKey := key[1:len(key)-1]
			if v,ok := datas[key]; ok {
				delete(datas, key)
				return fmt.Sprintf("%v", v)
			} else if v,ok := webParams[argsKey]; ok {
				//params必须优先,因为args默认值,就变0了
				return fmt.Sprintf("%v", v)
			} else if v,ok := argsValue[argsKey]; ok {
				return fmt.Sprintf("%v", v)
			} else {
				//有参数没有值,
				return p
			}
		})


		//get参数
		querys := []string{}
		for k,v := range datas {
			querys = append(querys, fmt.Sprintf("%v=%v", k, v))
		}
		if len(querys) > 0 {
			uri += "?" + strings.Join(querys, "&")
		}

		return uri
	}

	return "[no route here]"
}



//站点级路由URL
func (url *httpUrl) Routo(site,name string, args ...Map) string {

	if url.ctx != nil {
		if site == "" { site = url.ctx.Node.Name }
		if name == "" { name = url.ctx.Name }
	}

	config := Map{}
	uri := ""
	if siteConfig,ok := Http.routes[site]; ok {
		if config,ok := siteConfig[name]; ok {
			if c,ok := config["uri"]; ok {
				switch v := c.(type) {
				case string:
					uri = v
				case []string:
					uri = v[0]
				}
			}

		}
	}



	if uri != "" {
		argsConfig := Map{}
		if c,ok := config["args"]; ok {
			argsConfig = c.(Map)
		} else if c,ok := config["route"]; ok {
			routeConfig := c.(Map)
			if d,ok := routeConfig["args"]; ok {
				argsConfig = d.(Map)
			} else if d,ok := routeConfig["get"]; ok {
				methodConfig := d.(Map)
				if e,ok := methodConfig["args"]; ok {
					argsConfig = e.(Map)
				}
			}
		} else {
			//其它无参数
		}


		//得到值和选项
		datas,values,options := Map{}, Map{}, Map{}
		if len(args) > 0 {
			if len(args) == 1 {
				values = args[0]
			} else {
				values = args[0]
				options = args[1]
			}
		}



		//选项处理
		if (options["back"] != nil && url.ctx != nil) {
			var url = url.Back()
			datas[BACKURL] = encode64(url);
		}
		//选项处理
		if (options["last"] != nil && url.ctx != nil) {
			var url = url.Last()
			datas[BACKURL] = encode64(url);
		}
		//自动携带原有的query信息
		if (options["query"] != nil && url.ctx != nil) {
			for k,v := range url.ctx.Query {
				datas[k] = v
			}
		}





		//再从args,是不是有默认值
		argsValue := Map{}
		e := Mapping.Parse([]string{}, argsConfig, Map{}, argsValue)
		//不直接写datas,而是在下面的params中,如果有,才写入
		if e == nil {
			/*
			for k,v := range argsValue {
				datas["{"+k+"}"] = v
			}
			*/
		}

		//这里后写, 才会盖掉args中的值
		//先从web.Params拿值,这样如果是同一页,就有的相同的参数了
		//当前页相同的路由,才加入Params参数
		//if url.ctx != nil && name == url.ctx.Name {
		//不能直接复制
		/*
		if url.ctx != nil {
			for k,v := range url.ctx.Params {
				datas["{"+k+"}"] = v
			}
		}
		*/
		webParams := Map{}
		if url.ctx != nil {
			webParams = url.ctx.Param
		}

		//然后是传过来的值
		for k,v := range values {
			datas[k] = v
		}




		//keys := []string{}

		//直接用正则替换

		regx := regexp.MustCompile(`\{[_\*A-Za-z0-9]+\}`)
		uri = regx.ReplaceAllStringFunc(uri, func(p string) string {
			key := strings.Replace(p, "*", "", -1)
			argsKey := key[1:len(key)-1]
			if v,ok := datas[key]; ok {
				delete(datas, key)
				return fmt.Sprintf("%v", v)
			} else if v,ok := webParams[argsKey]; ok {
				//params必须优先,因为args默认值,就变0了
				return fmt.Sprintf("%v", v)
			} else if v,ok := argsValue[argsKey]; ok {
				return fmt.Sprintf("%v", v)
			} else {
				//有参数没有值,
				return p
			}
		})


		//get参数
		querys := []string{}
		for k,v := range datas {
			querys = append(querys, fmt.Sprintf("%v=%v", k, v))
		}
		if len(querys) > 0 {
			uri += "?" + strings.Join(querys, "&")
		}

		return url.Node(site, uri)
	}

	return "[no route here]"
}





func (url *httpUrl) Back() string {
	if url.ctx == nil {
		return "/"
	}

	if s,ok := url.ctx.Query[BACKURL]; ok {
		return decode64(s.(string))
	} else if url.ctx.Req.Header.Get("referer") != "" {
		return url.ctx.Req.Header.Get("referer")
	} else {
		//都没有，就是当前URL
		return url.Current()
	}
}



func (url *httpUrl) Last() string {
	if url.ctx == nil {
		return "/"
	}

	if url.ctx.Req.Header.Get("referer") != "" {
		return url.ctx.Req.Header.Get("referer")
	} else {
		return "/"
	}
}

func (url *httpUrl) Current() string {
	if url.ctx == nil {
		return "/"
	}

	//获取点节URL
	return fmt.Sprintf("%s%s", url.Node(url.ctx.Node.Name), url.ctx.Req.URL.RequestURI())
}