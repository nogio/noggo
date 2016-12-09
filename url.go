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
	if cfg,ok := url.ctx.Node.Http.routes[name]; ok {
		if c,ok := cfg["uri"]; ok {
			switch v := c.(type) {
			case string:
				uri = v
			case []string:
				uri = v[0]
			}
		}
		config = cfg
	}



	if uri != "" {
		argsConfig := Map{}
		if c,ok := config["args"].(Map); ok {
			argsConfig = c
		} else if routeConfig,ok := config["route"].(Map); ok {
			if d,ok := routeConfig["args"].(Map); ok {
				argsConfig = d
			} else if methodConfig,ok := routeConfig["get"].(Map); ok {
				if e,ok := methodConfig["args"].(Map); ok {
					argsConfig = e
				}
			}
		} else {
			//其它无参数
		}


		//得到值和选项
		//值有几个来源
		//1 datas中传过来的，最优先
		//2 args中的默认值，
		//3 params中现有的值
		//而1中的值，也要args解析，用来保证加解密的话，1和2其实就混在一起了
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




		//parse的时候， 直接传值进去，这样是不是科学一点？
		parseValues := Map{}

		for k,v := range values {
			k = strings.Replace(k, "{","",-1)
			k = strings.Replace(k, "}","",-1)
			parseValues[k] = v
		}
		//从params中来一下，直接参数解析
		if url.ctx != nil {
			for k,v := range url.ctx.Param {
				if parseValues[k] == nil {
					parseValues[k] = v
				}
			}
		}

		//再从args,是不是有默认值
		argsValue := Map{}
		e := Mapping.Parse([]string{}, argsConfig, parseValues, argsValue, false, true)
		//不直接写datas,而是在下面的params中,如果有,才写入
		if e == nil {

			//这里解析完， 默认值都过来了，不能直接替换写入datas，要不然传过来的值就无效了
			for k,v := range argsValue {
				//注意
				datas["{"+k+"}"] = v
			}
		}

		//然后是传过来的值，先把传的值写入？
		for k,v := range values {
			if datas[k] == nil {
				datas[k] = v
			}
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


		//这个单独处理
		/*
		webParams := Map{}
		if url.ctx != nil {
			webParams = url.ctx.Param
		}
		*/



		//keys := []string{}

		//直接用正则替换
		regx := regexp.MustCompile(`\{[_\*A-Za-z0-9]+\}`)
		uri = regx.ReplaceAllStringFunc(uri, func(p string) string {
			key := strings.Replace(p, "*", "", -1)
			if v,ok := datas[key]; ok {
				delete(datas, key)
				return fmt.Sprintf("%v", v)
			} else {
				//有参数没有值,
				return p
			}
		})






		//get参数
		querys := []string{}
		for k,v := range datas {
			//这里还要清理{}包裹的参数，因为上面parse的时候，如果路由有默认值的参数，而uri中没有，就会带进来
			if k[0:1] != "{" {
				querys = append(querys, fmt.Sprintf("%v=%v", k, v))
			}
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
		if cfg,ok := siteConfig[name]; ok {
			if c,ok := cfg["uri"]; ok {
				switch v := c.(type) {
				case string:
					uri = v
				case []string:
					uri = v[0]
				}
			}
			config = cfg
		}
	}



	if uri != "" {
		argsConfig := Map{}
		if c,ok := config["args"].(Map); ok {
			argsConfig = c
		} else if routeConfig,ok := config["route"].(Map); ok {
			if d,ok := routeConfig["args"].(Map); ok {
				argsConfig = d
			} else if methodConfig,ok := routeConfig["get"].(Map); ok {
				if e,ok := methodConfig["args"].(Map); ok {
					argsConfig = e
				}
			}
		} else {
			//其它无参数
		}



		//得到值和选项
		//值有几个来源
		//1 datas中传过来的，最优先
		//2 args中的默认值，
		//3 params中现有的值
		//而1中的值，也要args解析，用来保证加解密的话，1和2其实就混在一起了
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




		//parse的时候， 直接传值进去，这样是不是科学一点？
		parseValues := Map{}

		for k,v := range values {
			k = strings.Replace(k, "{","",-1)
			k = strings.Replace(k, "}","",-1)
			parseValues[k] = v
		}
		//从params中来一下，直接参数解析
		if url.ctx != nil {
			for k,v := range url.ctx.Param {
				if parseValues[k] == nil {
					parseValues[k] = v
				}
			}
		}

		//再从args,是不是有默认值
		argsValue := Map{}
		e := Mapping.Parse([]string{}, argsConfig, parseValues, argsValue, false, true)
		//不直接写datas,而是在下面的params中,如果有,才写入
		if e == nil {

			//这里解析完， 默认值都过来了，不能直接替换写入datas，要不然传过来的值就无效了
			for k,v := range argsValue {
				//注意
				datas["{"+k+"}"] = v
			}
		}

		//然后是传过来的值，先把传的值写入？
		for k,v := range values {
			if datas[k] == nil {
				datas[k] = v
			}
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


		//这个单独处理
		/*
		webParams := Map{}
		if url.ctx != nil {
			webParams = url.ctx.Param
		}
		*/



		//keys := []string{}

		//直接用正则替换
		regx := regexp.MustCompile(`\{[_\*A-Za-z0-9]+\}`)
		uri = regx.ReplaceAllStringFunc(uri, func(p string) string {
			key := strings.Replace(p, "*", "", -1)
			if v,ok := datas[key]; ok {
				delete(datas, key)
				return fmt.Sprintf("%v", v)
			} else {
				//有参数没有值,
				return p
			}
		})






		//get参数
		querys := []string{}
		for k,v := range datas {
			//这里还要清理{}包裹的参数，因为上面parse的时候，如果路由有默认值的参数，而uri中没有，就会带进来
			if k[0:1] != "{" {
				querys = append(querys, fmt.Sprintf("%v=%v", k, v))
			}
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