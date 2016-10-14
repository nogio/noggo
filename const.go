package noggo

import (
	. "github.com/nogio/noggo/base"
	"regexp"
	"fmt"
)



const (

	ConstDriverDefault  = "default"
	ConstLangDefault	= "default"
	ConstNodeGlobal		= "$global$"




	KeySign				= "noggo.sign."
	KeySignId			= "id"
	KeySignName			= "name"
	KeySignData			= "data"

	KeyMapUri		= "uri"
	KeyMapTime		= "time"
	KeyMapLine		= "line"
	KeyMapName		= "name"
	KeyMapText		= "text"
	KeyMapType		= "type"
	KeyMapMust		= "must"
	KeyMapAuto		= "auto"
	KeyMapJson		= "json"
	KeyMapEnum		= "enum"
	KeyMapArgs		= "args"
	KeyMapItem		= "item"
	KeyMapAuth		= "auth"

	KeyMapAction	= "action"
	KeyMapBranch	= "branch"
	KeyMapRoute		= "route"
	KeyMapMatch		= "match"

	KeyMapFound		= "found"
	KeyMapError		= "error"
	KeyMapFailed	= "failed"
	KeyMapDenied	= "denied"

	KeyMapEncode	= "encode"
	KeyMapDecode	= "decode"
	KeyMapValid		= "valid"
	KeyMapValue		= "value"
)





type (
	constGlobal struct {
		mimes		map[string]string
		states		map[string]int
		regulars	map[string][]string
		langs		map[string]map[string]string
	}
)



//注册mime
func (global *constGlobal) Mime(mimes Map) {
	if global.mimes == nil {
		global.mimes = map[string]string{}
	}
	for k,v := range mimes {
		switch vvvv := v.(type) {
		case string:
			global.mimes[k] = vvvv
		}
	}
}
//获取mimetype
func (global *constGlobal) MimeType(name string) string {
	if v,ok := global.mimes[name]; ok {
		return v
	} else {
		return name
	}
}


//注册State
func (global *constGlobal) State(states Map) {
	if global.states == nil {
		global.states = map[string]int{}
	}
	for k,v := range states {
		switch vvvv := v.(type) {
		case int:
			global.states[k] = vvvv
		}
	}
}
//获取获取码
func (global *constGlobal) StateCode(name string) int {
	if v,ok := global.states[name]; ok {
		return v
	} else {
		return -1
	}
}

//注册Regular
func (global *constGlobal) Regular(regulars Map) {
	if global.regulars == nil {
		global.regulars = map[string][]string{}
	}
	for k,v := range regulars {
		switch vvvv := v.(type) {
		case string:
			global.regulars[k] = []string{ vvvv }
		case []string:
			global.regulars[k] = vvvv
		}
	}
}
//获取正则的值
func (global *constGlobal) RegularExp(name string) []string {
	if v,ok := global.regulars[name]; ok {
		return v
	} else {
		return []string{ name }
	}
}




//注册String
func (global *constGlobal) Lang(lang string, strings Map) {
	if global.langs == nil {
		global.langs = map[string]map[string]string{}
	}
	if global.langs[lang] == nil {
		global.langs[lang] = map[string]string{}
	}


	for k,v := range strings {
		switch vvvv := v.(type) {
		case string:
			global.langs[lang][k] = vvvv
		}
	}
}
//获取语言字串
func (global *constGlobal) LangString(name string, langs ...string) string {
	lang := ConstLangDefault
	if len(langs) > 0 {
		lang = langs[0]
	}

	if langs,ok := global.langs[lang]; ok {
		if v,ok := langs[lang]; ok {
			return v
		}
	}

	return name
}










//使用正在验证
func (global *constGlobal) Valid(value, regular string) bool {

	exps := global.RegularExp(regular)
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
func (global *constGlobal) NewStateError(state string, args ...interface{}) *Error {

	stateCode := global.StateCode(state)
	stateText := global.LangString(state)

	return NewCodeError(stateCode, fmt.Sprintf(stateText, args...))
}


//按状态生成错误
//args是用来Format错误信息的
//state 和langstring 联动
func (global *constGlobal) NewLangStateError(lang, state string, args ...interface{}) *Error {

	stateCode := global.StateCode(state)
	stateText := global.LangString(state, lang)

	return NewCodeError(stateCode, fmt.Sprintf(stateText, args...))
}


