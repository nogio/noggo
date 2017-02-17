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
	uuu := ""

	//逻辑有小调整，先从config.site中获取，如果没有，再从config.node获取，再没有，就显示127.0.0.1+port
	if vv,ok := Config.Site[name]; ok && vv!=nil && len(vv)>0 {
		//如果是多个，可以随机选一个
		uuu = vv[0]
	} else if vv,ok := Config.Node[name]; ok {
		if vv.Url != "" {
			uuu = vv.Url
		} else {
			uuu = "http://127.0.0.1" + vv.Port
		}
	} else {
		uuu = ""
	}

	if len(args) > 0 {
		return fmt.Sprintf("%s%s", uuu, args[0])
	} else {
		return uuu
	}
}

func (url *httpUrl) Root(args ...string) string {
	root := ""
	if url.ctx != nil && url.ctx.Node != nil && url.ctx.Node.Config != nil {
		root = url.ctx.Node.Config.Root
	}

	if len(args) > 0 {
		return fmt.Sprintf("%s%s", root, args[0])
	} else {
		return root
	}
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
		values,options := Map{}, Map{}
		if len(args) > 0 {
			if len(args) == 1 {
				values = args[0]
			} else {
				values = args[0]
				options = args[1]
			}
		}




		queryValues := Map{}

		//选项处理
		if (options["back"] != nil && url.ctx != nil) {
			var url = url.Back()
			queryValues[BACKURL] = encode64(url);
		}
		//选项处理
		if (options["last"] != nil && url.ctx != nil) {
			var url = url.Last()
			queryValues[BACKURL] = encode64(url);
		}
		//自动携带原有的query信息
		if (options["query"] != nil && url.ctx != nil) {
			for k,v := range url.ctx.Query {
				queryValues[k] = v
			}
		}





		//所以，解析uri中的参数，值得分几类：
		//1传的值，2param值, 3默认值
		//其中主要问题就是，传的值，需要到args解析，用于加密，这个值和auto值完全重叠了，除非分2次解析
		//为了框架好用，真是操碎了心
		dataValues, paramValues, autoValues := Map{},Map{},Map{}

		//1. 处理传过来的值
		//从value中获取
		//如果route不定义args，这里是拿不到值的
		dataArgsValues, dataParseValues := Map{},Map{}
		for k,v := range values {
			if k[0:1] == "{" {
				k = strings.Replace(k, "{","",-1)
				k = strings.Replace(k, "}","",-1)
				dataArgsValues[k] = v
			} else {
				//这个也要？要不然指定的一些page啥的不行？
				dataArgsValues[k] = v
				//另外的是query的值
				queryValues[k] = v
			}
		}
		dataErr := Mapping.Parse([]string{}, argsConfig, dataArgsValues, dataParseValues, false, true)
		if dataErr == nil {
			for k,v := range dataParseValues {

				//注意，这里能拿到的，还有非param，所以不能直接用加{}写入
				if _,ok := values[k]; ok {
					dataValues[k] = v
				} else if _,ok := values["{"+k+"}"]; ok {
					dataValues["{"+k+"}"] = v
				} else {
					//这里是默认值应该，就不需要了
				}
			}
		}
		//所以这里还得处理一次，如果route不定义args，parse就拿不到值，就直接用values中的值
		for k,v := range values {
			if k[0:1] == "{" && dataValues[k] == nil {
				dataValues[k] = v
			}
		}


		//2.params中的值
		//从params中来一下，直接参数解析
		if url.ctx != nil {
			for k,v := range url.ctx.Param {
				paramValues["{"+k+"}"] = v
			}
		}


		//3. 默认值
		//从value中获取
		autoArgsValues, autoParseValues := Map{},Map{}
		autoErr := Mapping.Parse([]string{}, argsConfig, autoArgsValues, autoParseValues, false, true)
		if autoErr == nil {
			for k,v := range autoParseValues {
				autoValues["{"+k+"}"] = v
			}
		}

		//Logger.Debug(name, "data", dataValues, "param", paramValues, "auto", autoValues)

		//开始替换值
		regx := regexp.MustCompile(`\{[_\*A-Za-z0-9]+\}`)
		uri = regx.ReplaceAllStringFunc(uri, func(p string) string {
			key := strings.Replace(p, "*", "", -1)

			if v,ok := dataValues[key]; ok {
				//先从传的值去取
				return fmt.Sprintf("%v", v)
			} else if v,ok := paramValues[key]; ok {
				//再从params中去取
				return fmt.Sprintf("%v", v)
			} else if v,ok := autoValues[key]; ok {
				//最后从默认值去取
				return fmt.Sprintf("%v", v)
			} else {
				//有参数没有值,
				return p
			}
		})



		/*


		//parse的时候， 直接传值进去，这样是不是科学一点？
		parseValues := Map{}

		for k,v := range values {
			//{开头的参数 才做转换处理，query中的参数 不处理
			//但是这里如果跳过的话route中可能有默认值，会变成默认值返回到argsValue
			//不行，因为page=%v这样的情况，而page有默认值，就不行了

			k = strings.Replace(k, "{","",-1)
			k = strings.Replace(k, "}","",-1)
			parseValues[k] = v
		}
		//从params中来一下，直接参数解析
		if url.ctx != nil {
			for k,v := range url.ctx.Param {
				paramValues["{"+k+"}"] = v
				if parseValues[k] == nil {
					//parse不能把param写进去，
					//因为如果param中的值加密过，在parse的时候， 就给解密了， 这不科学
					//所以还是要写进去，因为：
					//当route有默认值的时候， 而我调用route的时候， 只传了1个值，
					//其它的全都有默认值了，就不会跑param中去拿值了
					//最终还是不能复制param的过去

					//在下面args的时候判断，如果args中的k不在value中，就应该干掉，不写入datas

					//parseValues[k] = v
				}
			}
		}


		/*

		//再从args,是不是有默认值
		//小问题一个，如果目标路由，没有定义参数，比如/admin/remove/{id}
		//这里并不定义参数，所以在这里，上面的parseValue写入param其实就没有意义
		argsValue := Map{}
		e := Mapping.Parse([]string{}, argsConfig, parseValues, argsValue, false, true)


		//不直接写datas,而是在下面的params中,如果有,才写入
		if e == nil {
			for k,v := range argsValue {
				//注意，这里能拿到的，还有非param，所以不能直接用加{}写入
				if _,ok := values[k]; ok {
					datas[k] = v
				} else if _,ok := values["{"+k+"}"]; ok {
					datas["{"+k+"}"] = v
				} else {
					//到这里， 应该是默认值
					datas["{"+k+"}"] = v
				}
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

		/*


		//keys := []string{}

		//直接用正则替换
		regx := regexp.MustCompile(`\{[_\*A-Za-z0-9]+\}`)
		uri = regx.ReplaceAllStringFunc(uri, func(p string) string {
			key := strings.Replace(p, "*", "", -1)
			if v,ok := datas[key]; ok {
				delete(datas, key)
				return fmt.Sprintf("%v", v)
			} else if v,ok := paramValues[key]; ok {
				return fmt.Sprintf("%v", v)
			} else {
				//有参数没有值,
				return p
			}
		})

		*/





		//get参数
		querys := []string{}
		for k,v := range queryValues {
			//这里还要清理{}包裹的参数，因为上面parse的时候，如果路由有默认值的参数，而uri中没有，就会带进来
			if k[0:1] != "{" {
				querys = append(querys, fmt.Sprintf("%v=%v", k, v))
			}
		}
		if len(querys) > 0 {
			uri += "?" + strings.Join(querys, "&")
		}

		if url.ctx != nil && url.ctx.Node != nil && url.ctx.Node.Config != nil {
			return url.ctx.Node.Config.Root + uri
		} else {
			return uri
		}

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
		values,options := Map{}, Map{}
		if len(args) > 0 {
			if len(args) == 1 {
				values = args[0]
			} else {
				values = args[0]
				options = args[1]
			}
		}




		queryValues := Map{}

		//选项处理
		if (options["back"] != nil && url.ctx != nil) {
			var url = url.Back()
			queryValues[BACKURL] = encode64(url);
		}
		//选项处理
		if (options["last"] != nil && url.ctx != nil) {
			var url = url.Last()
			queryValues[BACKURL] = encode64(url);
		}
		//自动携带原有的query信息
		if (options["query"] != nil && url.ctx != nil) {
			for k,v := range url.ctx.Query {
				queryValues[k] = v
			}
		}





		//所以，解析uri中的参数，值得分几类：
		//1传的值，2param值, 3默认值
		//其中主要问题就是，传的值，需要到args解析，用于加密，这个值和auto值完全重叠了，除非分2次解析
		//为了框架好用，真是操碎了心
		dataValues, paramValues, autoValues := Map{},Map{},Map{}

		//1. 处理传过来的值
		//从value中获取
		//这里有个问题，如果route中没有定义args，那应该是拿不到值了
		dataArgsValues, dataParseValues := Map{},Map{}
		for k,v := range values {
			if k[0:1] == "{" {
				k = strings.Replace(k, "{","",-1)
				k = strings.Replace(k, "}","",-1)
				dataArgsValues[k] = v
			} else {
				//这个也要？要不然指定的一些page啥的不行？
				dataArgsValues[k] = v
				//另外的是query的值
				queryValues[k] = v
			}
		}
		dataErr := Mapping.Parse([]string{}, argsConfig, dataArgsValues, dataParseValues, false, true)
		if dataErr == nil {
			for k,v := range dataParseValues {
				//注意，这里能拿到的，还有非param，所以不能直接用加{}写入
				if _,ok := values[k]; ok {
					dataValues[k] = v
				} else if _,ok := values["{"+k+"}"]; ok {
					dataValues["{"+k+"}"] = v
				} else {
					//这里是默认值应该，就不需要了
				}
			}
		}
		//所以这里还得处理一次，如果route不定义args，parse就拿不到值，就直接用values中的值
		for k,v := range values {
			if k[0:1] == "{" && dataValues[k] == nil {
				dataValues[k] = v
			}
		}


		//2.params中的值
		//从params中来一下，直接参数解析
		if url.ctx != nil {
			for k,v := range url.ctx.Param {
				paramValues["{"+k+"}"] = v
			}
		}


		//3. 默认值
		//从value中获取
		autoArgsValues, autoParseValues := Map{},Map{}
		autoErr := Mapping.Parse([]string{}, argsConfig, autoArgsValues, autoParseValues, false, true)
		if autoErr == nil {
			for k,v := range autoParseValues {
				autoValues["{"+k+"}"] = v
			}
		}

		//开始替换值
		regx := regexp.MustCompile(`\{[_\*A-Za-z0-9]+\}`)
		uri = regx.ReplaceAllStringFunc(uri, func(p string) string {
			key := strings.Replace(p, "*", "", -1)

			if v,ok := dataValues[key]; ok {
				//先从传的值去取
				return fmt.Sprintf("%v", v)
			} else if v,ok := paramValues[key]; ok {
				//再从params中去取
				return fmt.Sprintf("%v", v)
			} else if v,ok := autoValues[key]; ok {
				//最后从默认值去取
				return fmt.Sprintf("%v", v)
			} else {
				//有参数没有值,
				return p
			}
		})



		/*


		//parse的时候， 直接传值进去，这样是不是科学一点？
		parseValues := Map{}

		for k,v := range values {
			//{开头的参数 才做转换处理，query中的参数 不处理
			//但是这里如果跳过的话route中可能有默认值，会变成默认值返回到argsValue
			//不行，因为page=%v这样的情况，而page有默认值，就不行了

			k = strings.Replace(k, "{","",-1)
			k = strings.Replace(k, "}","",-1)
			parseValues[k] = v
		}
		//从params中来一下，直接参数解析
		if url.ctx != nil {
			for k,v := range url.ctx.Param {
				paramValues["{"+k+"}"] = v
				if parseValues[k] == nil {
					//parse不能把param写进去，
					//因为如果param中的值加密过，在parse的时候， 就给解密了， 这不科学
					//所以还是要写进去，因为：
					//当route有默认值的时候， 而我调用route的时候， 只传了1个值，
					//其它的全都有默认值了，就不会跑param中去拿值了
					//最终还是不能复制param的过去

					//在下面args的时候判断，如果args中的k不在value中，就应该干掉，不写入datas

					//parseValues[k] = v
				}
			}
		}


		/*

		//再从args,是不是有默认值
		//小问题一个，如果目标路由，没有定义参数，比如/admin/remove/{id}
		//这里并不定义参数，所以在这里，上面的parseValue写入param其实就没有意义
		argsValue := Map{}
		e := Mapping.Parse([]string{}, argsConfig, parseValues, argsValue, false, true)


		//不直接写datas,而是在下面的params中,如果有,才写入
		if e == nil {
			for k,v := range argsValue {
				//注意，这里能拿到的，还有非param，所以不能直接用加{}写入
				if _,ok := values[k]; ok {
					datas[k] = v
				} else if _,ok := values["{"+k+"}"]; ok {
					datas["{"+k+"}"] = v
				} else {
					//到这里， 应该是默认值
					datas["{"+k+"}"] = v
				}
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

		/*


		//keys := []string{}

		//直接用正则替换
		regx := regexp.MustCompile(`\{[_\*A-Za-z0-9]+\}`)
		uri = regx.ReplaceAllStringFunc(uri, func(p string) string {
			key := strings.Replace(p, "*", "", -1)
			if v,ok := datas[key]; ok {
				delete(datas, key)
				return fmt.Sprintf("%v", v)
			} else if v,ok := paramValues[key]; ok {
				return fmt.Sprintf("%v", v)
			} else {
				//有参数没有值,
				return p
			}
		})

		*/





		//get参数
		querys := []string{}
		for k,v := range queryValues {
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