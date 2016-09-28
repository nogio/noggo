package noggo

import (
	. "github.com/nogio/noggo/base"
	"regexp"
	"fmt"
)


type (
	constModule struct {
		mimes		map[string]string
		states		map[string]int
		regulars	map[string][]string
		langs		map[string]map[string]string
	}
)



//注册mime
func (module *constModule) Mime(mimes Map) {
	if module.mimes == nil {
		module.mimes = map[string]string{}
	}
	for k,v := range mimes {
		switch vvvv := v.(type) {
		case string:
			module.mimes[k] = vvvv
		}
	}
}
//获取mimetype
func (module *constModule) MimeType(name string) string {
	if v,ok := module.mimes[name]; ok {
		return v
	} else {
		return name
	}
}


//注册State
func (module *constModule) State(states Map) {
	if module.states == nil {
		module.states = map[string]int{}
	}
	for k,v := range states {
		switch vvvv := v.(type) {
		case int:
			module.states[k] = vvvv
		}
	}
}
//获取获取码
func (module *constModule) StateCode(name string) int {
	if v,ok := module.states[name]; ok {
		return v
	} else {
		return -1
	}
}

//注册Regular
func (module *constModule) Regular(regulars Map) {
	if module.regulars == nil {
		module.regulars = map[string][]string{}
	}
	for k,v := range regulars {
		switch vvvv := v.(type) {
		case string:
			module.regulars[k] = []string{ vvvv }
		case []string:
			module.regulars[k] = vvvv
		}
	}
}
//获取正则的值
func (module *constModule) RegularExp(name string) []string {
	if v,ok := module.regulars[name]; ok {
		return v
	} else {
		return []string{ name }
	}
}




//注册String
func (module *constModule) Lang(strings Map, langs ...string) {
	lang := ConstLangDefault
	if len(langs) > 0 {
		lang = langs[0]
	}

	if module.langs == nil {
		module.langs = map[string]map[string]string{}
	}
	if module.langs[lang] == nil {
		module.langs[lang] = map[string]string{}
	}


	for k,v := range strings {
		switch vvvv := v.(type) {
		case string:
			module.langs[lang][k] = vvvv
		}
	}
}
//获取语言字串
func (module *constModule) LangString(name string, langs ...string) string {
	lang := ConstLangDefault
	if len(langs) > 0 {
		lang = langs[0]
	}

	if langs,ok := module.langs[lang]; ok {
		if v,ok := langs[lang]; ok {
			return v
		}
	}

	return name
}










//使用正在验证
func (module *constModule) Valid(value, regular string) bool {

	exps := module.RegularExp(regular)
	for _,v := range exps {
		regx := regexp.MustCompile(v)
		if regx.MatchString(value) {
			return true
		}
	}

	return false
}





//按状态生成错误
//args是用来Format错误信息的
//state 和langstring 联动
func (module *constModule) NewStateError(state string, args ...interface{}) *Error {

	stateCode := module.StateCode(state)
	stateText := module.LangString(state)

	return NewCodeError(stateCode, fmt.Sprintf(stateText, args...))
}


//按状态生成错误
//args是用来Format错误信息的
//state 和langstring 联动
func (module *constModule) NewLangStateError(lang, state string, args ...interface{}) *Error {

	stateCode := module.StateCode(state)
	stateText := module.LangString(state, lang)

	return NewCodeError(stateCode, fmt.Sprintf(stateText, args...))
}


