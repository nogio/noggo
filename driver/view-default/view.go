package view_default



import (
	. "github.com/nogio/noggo/base"
	"github.com/nogio/noggo/driver"
	"errors"
	"fmt"
	"html/template"
	"path/filepath"
	"os"
	"bytes"
	"strings"
	"encoding/json"
	"time"
)


type (
	DefaultView struct {
		root    string  //根目录
		parse   *driver.ViewParse

		engine *template.Template
		helper template.FuncMap

		body string     //解析后的body暂存
		path string     //记录body当前的目录
		layout string   //设置的layout view
		model Map       //此为layout所用的model

		title,author,description,keywords string
		metas, styles, scripts []string
		heads, headers, footers []string
		blocks map[string][]string
	}
)

func newDefaultView(root string, parse *driver.ViewParse) (*DefaultView) {
	view := &DefaultView{ root: root, parse: parse }
	view.metas = []string{}
	view.styles = []string{}
	view.scripts = []string{}
	view.heads = []string{}
	view.headers = []string{}
	view.footers = []string{}
	view.blocks = map[string][]string{}

	//工具方法
	view.helper = template.FuncMap{

		//支持布局页
		"layout": func(name string, vals ...Any) string {

			args := []Map{}
			for _,v := range vals {
				switch t := v.(type) {
				case Map:
					args = append(args, t)
				case string:
					m := Map{}
					e := json.Unmarshal([]byte(t), &m)
					if e == nil {
						args = append(args, m)
					}
				}
			}

			view.layout = name
			if len(args) > 0 {
				view.model = args[0]
			}

			return ""
		},
		"title": func(args ...string) template.HTML {
			if len(args) > 0 {
				//设置TITLE
				view.title = args[0]
				return template.HTML("")
			} else {
				if view.title != "" {
					return template.HTML(view.title)
				} else {
					return template.HTML("")
				}
			}
		},
		"author": func(args ...string) template.HTML {
			if len(args) >0 {
				view.author = args[0]
				return template.HTML("")
			} else {
				if view.author != "" {
					return template.HTML(view.author)
				} else {
					return template.HTML("")
				}
			}
		},
		"description": func(args ...string) template.HTML {
			if len(args) >0 {
				view.description = args[0]
				return template.HTML("")
			} else {
				if view.description != "" {
					return template.HTML(view.description)
				} else {
					return template.HTML("")
				}
			}
		},
		"keywords": func(args ...string) template.HTML {
			if len(args) >0 {
				view.keywords = args[0]
				return template.HTML("")
			} else {
				if view.author != "" {
					return template.HTML(view.keywords)
				} else {
					return template.HTML("")
				}
			}
		},


		"body": func() template.HTML {
			return template.HTML(view.body)
		},
		"render": func(name string, vals ...Any) template.HTML {
			args := []Map{}
			for _,v := range vals {
				switch t := v.(type) {
				case Map:
					args = append(args, t)
				case string:
					m := Map{}
					e := json.Unmarshal([]byte(t), &m)
					if e == nil {
						args = append(args, m)
					}
				}
			}

			s,e := view.Render(name, args...)
			if e == nil {
				return template.HTML(s)
			} else {
				return template.HTML(fmt.Sprintf("render error: %v", e))
			}
		},
		"head": func(args ...string) template.HTML {

			if len(args) > 0 {
				//有参数,是写入数据
				for _,v := range args {
					view.heads = append(view.heads, v)
				}
				return template.HTML("")
			} else {
				html := strings.Join(view.heads, "\n")
				return template.HTML(html)
			}
		},

		"header": func(args ...string) template.HTML {

			if len(args) > 0 {
				//有参数,是写入数据
				for _,v := range args {
					view.headers = append(view.headers, v)
				}
				return template.HTML("")
			} else {
				html := strings.Join(view.headers, "\n")
				return template.HTML(html)
			}
		},
		"footer": func(args ...string) template.HTML {

			if len(args) > 0 {
				//有参数,是写入数据
				for _,v := range args {
					view.footers = append(view.footers, v)
				}
				return template.HTML("")
			} else {
				html := strings.Join(view.footers, "\n")
				return template.HTML(html)
			}
		},
		"block": func(name string,args ...string) template.HTML {

			if len(args) > 0 {
				if view.blocks[name] == nil {
					view.blocks[name] = []string{}
				}

				//有参数,是写入数据
				for _,v := range args {
					view.footers = append(view.blocks[name], v)
				}
				return template.HTML("")
			} else {
				html := ""
				if v,ok := view.blocks[name]; ok {
					html = strings.Join(v, "\n")
				}
				return template.HTML(html)
			}
		},
		"meta": func(name,content string, https ...bool) (template.HTML) {
			isHttp := false
			if len(https) > 0 {
				isHttp = https[0]
			}
			if isHttp {
				view.metas = append(view.metas, fmt.Sprintf(`<meta http-equiv="%v" content="%v" />`, name, content))
			} else {
				view.metas = append(view.metas, fmt.Sprintf(`<meta name="%v" content="%v" />`, name, content))
			}

			return template.HTML("")
		},
		"metas": func() template.HTML {
			html := ""
			if len(view.metas) > 0 {
				html = strings.Join(view.metas, "\n")
			}
			return template.HTML(html)
		},
		"style": func(path string, args ...string) (template.HTML) {
			media := ""
			if len(args) > 0 {
				media = args[0]
			}
			if media == "" {
				view.styles = append(view.styles, fmt.Sprintf(`<link type="text/css" rel="stylesheet" href="%v" />`, path))
			} else {
				view.styles = append(view.styles, fmt.Sprintf(`<link type="text/css" rel="stylesheet" href="%v" media="%v" />`, path, media))
			}
			return template.HTML("")
		},
		"styles": func() template.HTML {
			html := ""
			if len(view.styles) > 0 {
				html = strings.Join(view.styles, "\n")
			}
			return template.HTML(html)
		},
		"script": func(path string, args ...string) (template.HTML) {
			tttt := "text/javascript"
			if len(args) > 0 {
				tttt = args[0]
			}
			view.scripts = append(view.scripts, fmt.Sprintf(`<script type="%v" src="%v"></script>`, tttt, path))

			return template.HTML("")
		},
		"scripts": func() template.HTML {
			html := ""
			if len(view.scripts) > 0 {
				html = strings.Join(view.scripts, "\n")
			}
			return template.HTML(html)
		},







		"html": func(text Any) (template.HTML) {
			if text != nil {
				return template.HTML(fmt.Sprintf("%v", text))
			}
			return template.HTML("")
		},
		"join": func(a Any, s string) template.HTML {
			strs := []string{}

			if a != nil {


				switch vv := a.(type) {
				case []string:
					for _,v := range vv {
						strs = append(strs, v)
					}
				case []Any:
					for _,v := range vv {
						strs = append(strs, fmt.Sprintf("%v", v))
					}
				case []interface{}:
					for _,v := range vv {
						strs = append(strs, fmt.Sprintf("%v", v))
					}
				case []int:
					for _,v := range vv {
						strs = append(strs, fmt.Sprintf("%v", v))
					}
				case []int8:
					for _,v := range vv {
						strs = append(strs, fmt.Sprintf("%v", v))
					}
				case []int16:
					for _,v := range vv {
						strs = append(strs, fmt.Sprintf("%v", v))
					}
				case []int32:
					for _,v := range vv {
						strs = append(strs, fmt.Sprintf("%v", v))
					}
				case []int64:
					for _,v := range vv {
						strs = append(strs, fmt.Sprintf("%v", v))
					}
				case []float32:
					for _,v := range vv {
						strs = append(strs, fmt.Sprintf("%v", v))
					}
				case []float64:
					for _,v := range vv {
						strs = append(strs, fmt.Sprintf("%v", v))
					}
				}
			}

			html := strings.Join(strs, s)
			return template.HTML(html)
		},
		"json": func(data Any) (template.HTML) {
			if data != nil {
				bytes, err := json.Marshal(data)
				if err == nil {
					return template.HTML(string(bytes))
				}
			}
			return template.HTML("")
		},
		"format": func(format string, args ...interface{}) (string) {
			//支持一下显示时间
			if len(args) == 1{
				if ttt,ok := args[0].(time.Time); ok {
					return ttt.Format(format)
				} else if ttt,ok := args[0].(int64); ok {
					//时间戳是大于1971年是, 千万级, 2016年就是10亿级了

					if ttt >= int64(31507200) {
						sss := time.Unix(ttt, 0).Format(format)
						//if sss != "%25v" {
						if strings.HasPrefix(sss, "%")==false || format != sss {
							return sss
						}
					}
				}
			}
			return fmt.Sprintf(format, args...)
		},
		"now": func() time.Time {
			return time.Now()
		},
		"in": func(val Any, arr Any) (bool) {

			strVal := fmt.Sprintf("%v", val)
			strArr := []string{}

			switch vv := arr.(type) {
			case []string:
				for _,v := range vv {
					strArr = append(strArr, v)
				}
			case []int:
				for _,v := range vv {
					strArr = append(strArr, fmt.Sprintf("%v", v))
				}
			case []int8:
				for _,v := range vv {
					strArr = append(strArr, fmt.Sprintf("%v", v))
				}
			case []int16:
				for _,v := range vv {
					strArr = append(strArr, fmt.Sprintf("%v", v))
				}
			case []int32:
				for _,v := range vv {
					strArr = append(strArr, fmt.Sprintf("%v", v))
				}
			case []int64:
				for _,v := range vv {
					strArr = append(strArr, fmt.Sprintf("%v", v))
				}
			}

			for _,v := range strArr {
				if v == strVal {
					return true
				}
			}

			return false
		},

	}
	for k,v := range parse.Helpers {
		view.helper[k] = v
	}

	return view
}


func (view *DefaultView) Parse() (string,error) {
	view.title = ""

	view.engine = template.New("default").Delims("<%", "%>").Funcs(view.helper)
	return view.Layout(view.parse.View,view.parse.Model)
}



//这里实际是在解析layout
//注意，这里的name,model是body的
//layout的name,model要在 layout方法中调用， 记录到  view对象中的 layout, model
func (view *DefaultView) Layout(name string, model Map) (string,error) {

	bodyText,bodyError := view.Body(name, model)
	if bodyError != nil {
		return "",bodyError
	} else {


		if view.layout == "" {
			//没有使用布局，直接返回BODY
			return bodyText,nil
		} else {

			//body赋值
			view.body = bodyText

			//先搜索layout所在目录
			viewpaths := []string{};
			if view.path != "" {
				viewpaths = append(viewpaths, fmt.Sprintf("%s/%s.html", view.path, view.layout))
			}
			viewpaths = append(viewpaths, fmt.Sprintf("%s/%s/default/%s.html", view.root, view.parse.Node, view.layout))
			viewpaths = append(viewpaths, fmt.Sprintf("%s/%s/%s.html", view.root, view.parse.Node, view.layout))
			viewpaths = append(viewpaths, fmt.Sprintf("%s/default/%s.html", view.root, view.layout))
			viewpaths = append(viewpaths, fmt.Sprintf("%s/%s.html", view.root, view.layout))


			var viewname string

			for _,s := range viewpaths {
				if f, _ := os.Stat(s); f != nil && !f.IsDir() {
					viewname = s
					break
				}
			}

			//如果view不存在
			if viewname == "" {
				return "",errors.New(fmt.Sprintf("view %s not exist", name))
			} else {

				//不直接使用 view.engine 来new,而是克隆一份
				engine,_ := view.engine.Clone()
				t,e := engine.New(filepath.Base(viewname)).ParseFiles(viewname)
				if e != nil {
					return "",errors.New(fmt.Sprintf("view %s parse error: %v", viewname, e))
				} else {

					//缓冲
					buf := bytes.NewBuffer(make([]byte, 0))

					e := t.Execute(buf, Map{
						"data":	view.parse.Data,
						"model": view.model,
					})
					if e != nil {
						return "",errors.New(fmt.Sprintf("view %s parse error: %v", viewname, e))
					} else {
						return buf.String(),nil
					}
				}
			}
		}
	}
}



/* 返回view */
func (view *DefaultView) Body(name string, args ...Map) (string,error) {
	bodyModel := Map{}
	if len(args) > 0 {
		bodyModel = args[0]
	}

	//定义View搜索的路径
	viewpaths := []string{
		fmt.Sprintf("%s/%s/%s.html", view.root, view.parse.Node, name),
		fmt.Sprintf("%s/%s/default/%s.html", view.root, view.parse.Node, name),
		fmt.Sprintf("%s/%s.html", view.root, name),
		fmt.Sprintf("%s/%s/index.html", view.root, name),
		fmt.Sprintf("%s/%s/%s/index.html", view.root, view.parse.Node, name),
		fmt.Sprintf("%s/default/%s.html", view.root, name),
		fmt.Sprintf("%s/default/%s/index.html", view.root, name),
	};




	var viewname string

	for _,s := range viewpaths {
		if f, _ := os.Stat(s); f != nil && !f.IsDir() {
			viewname = s
			//这里要保存body所在的目录，为当前目录
			view.path = filepath.Dir(s)
			break
		}
	}

	//如果view不存在
	if viewname == "" {
		return "",errors.New(fmt.Sprintf("view %s not exist", name))
	} else {

		//不直接使用 view.engine 来new,而是克隆一份
		engine,_ := view.engine.Clone()
		t,e := engine.New(filepath.Base(viewname)).ParseFiles(viewname)
		if e != nil {
			return "",errors.New(fmt.Sprintf("view %s parse error: %v", viewname, e))
		} else {

			//缓冲
			buf := bytes.NewBuffer(make([]byte, 0))

			e := t.Execute(buf, Map{
				"data":	view.parse.Data,
				"model": bodyModel,
			})
			if e != nil {
				return "",errors.New(fmt.Sprintf("view %s parse error: %v", viewname, e))
			} else {
				return buf.String(),nil
			}
		}
	}

}

/* 返回view */
func (view *DefaultView) Render(name string, args ...Map) (string,error) {

	renderModel := Map{}
	if len(args) > 0 {
		renderModel = args[0]
	}

	//先搜索body所在目录
	viewpaths := []string{};
	if view.path != "" {
		viewpaths = append(viewpaths, fmt.Sprintf("%s/%s.html", view.path, name))
	}
	viewpaths = append(viewpaths, fmt.Sprintf("%s/%s/default/%s.html", view.root, view.parse.Node, name))
	viewpaths = append(viewpaths, fmt.Sprintf("%s/%s/%s.html", view.root, view.parse.Node, name))
	viewpaths = append(viewpaths, fmt.Sprintf("%s/default/%s.html", view.root, name))
	viewpaths = append(viewpaths, fmt.Sprintf("%s/%s.html", view.root, name))



	var viewname string
	for _,s := range viewpaths {
		if f, _ := os.Stat(s); f != nil && !f.IsDir() {
			viewname = s
			break
		}
	}


	//如果view不存在
	if viewname == "" {
		return "",errors.New(fmt.Sprintf("render %s not exist", name))
	} else {

		//不直接使用 view.engine 来new,而是克隆一份
		//因为1.6以后，不知道为什么，直接用，就会有问题
		//会报重复render某页面的问题
		engine,_ := view.engine.Clone()

		//如果一个模板被引用过了
		//不再重新加载文件
		//要不然, render某个页面,只能render一次
		n := filepath.Base(viewname)
		t := engine.Lookup(n)


		if t == nil {

			newT,e := engine.New(n).ParseFiles(viewname)
			if e != nil {
				return "",errors.New(fmt.Sprintf("render %s parse error: %v", name, e.Error()))
			} else {
				t = newT
			}
		}

		//缓冲
		buf := bytes.NewBuffer(make([]byte, 0))


		e := t.Execute(buf, Map{
			"data":	view.parse.Data,
			"model": renderModel,
		})
		if e != nil {
			return "",errors.New(fmt.Sprintf("view %s parse error: %v", viewname, e))
		} else {
			return buf.String(),nil
		}

	}


}

