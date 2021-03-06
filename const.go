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
	KeyMapNone 	    = "none"
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
func (global *constGlobal) Mime(name string, values ...string) (string) {
	if global.mimes == nil {
		global.mimes = map[string]string{}
	}

	if len(values) > 0 {
		global.mimes[name] = values[0]
		return values[0]
	} else {

		if v,ok := global.mimes[name]; ok {
			return v
		} else {
			return ""
		}
	}
}
func (global *constGlobal) Mimes(args ...Map) (map[string]string) {
	if global.mimes == nil {
		global.mimes = map[string]string{}
	}

	if len(args) > 0 {
		for k,v := range args[0] {
			switch vvvv := v.(type) {
			case string:
				global.mimes[k] = vvvv
			}
		}
	}

	return global.mimes
}
//获取mimetype
func (global *constGlobal) MimeType(name string, defs ...string) string {
	if v,ok := global.mimes[name]; ok {
		return v
	} else {
		if len(defs) > 0 {
			return defs[0]
		}
		return "application/octet-stream"
	}
}


//注册State
func (global *constGlobal) State(name string, codes ...int) (int) {
	if global.states == nil {
		global.states = map[string]int{}
	}


	if len(codes) > 0 {
		global.states[name] = codes[0]
		return codes[0]
	} else {

		if v,ok := global.states[name]; ok {
			return v
		} else {
			return 0
		}
	}
}

func (global *constGlobal) States(args ...Map) (map[string]int) {
	if global.states == nil {
		global.states = map[string]int{}
	}

	if len(args) > 0 {
		for k,v := range args[0] {
			switch vvvv := v.(type) {
			case int:
				global.states[k] = vvvv
			}
		}
	}

	return global.states
}

//获取获取码
func (global *constGlobal) StateCode(name string) int {
	if v,ok := global.states[name]; ok {
		return v
	} else {
		return -1
	}
}


//返回状态定字串列表
func (global *constGlobal) StateStrings(args ...string) (Map) {
	m := Map{}

	if len(args) > 0 {
		for _,n := range args {
			m[fmt.Sprintf("%d", global.StateCode(n))] = global.LangString(n)
		}
	} else {
		for k,v := range global.states {
			m[fmt.Sprintf("%d", v)] = global.LangString(k)
		}
	}


	return m
}









//注册Regular
func (global *constGlobal) Regular(name string, value Any) {
	if global.regulars == nil {
		global.regulars = map[string][]string{}
	}
	switch vvvv := value.(type) {
	case string:
		global.regulars[name] = []string{ vvvv }
	case []string:
		global.regulars[name] = vvvv
	}
}
//注册Regular
func (global *constGlobal) Regulars(regulars Map) {
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




//注册Lang
func (global *constGlobal) Lang(lang string, name string, values ...string) (string) {
	if global.langs == nil {
		global.langs = map[string]map[string]string{}
	}
	if global.langs[lang] == nil {
		global.langs[lang] = map[string]string{}
	}


	if len(values) > 0 {
		global.langs[lang][name] = values[0]
		return values[0]
	} else {
		if v,ok := global.langs[lang][name]; ok {
			return v
		} else {
			return ""
		}
	}
}
func (global *constGlobal) Langs(lang string, args ...Map) (map[string]string) {
	if global.langs == nil {
		global.langs = map[string]map[string]string{}
	}
	if global.langs[lang] == nil {
		global.langs[lang] = map[string]string{}
	}

	if len(args) > 0 {
		for k,v := range args[0] {
			switch vvvv := v.(type) {
			case string:
				global.langs[lang][k] = vvvv
			}
		}
	}

	return global.langs[lang]
}





//获取语言字串
//(key, lang)
func (global *constGlobal) LangString(name string, langs ...string) string {
	lang := ConstLangDefault
	if len(langs) > 0 {
		lang = langs[0]
	}

	if langs,ok := global.langs[lang]; ok {
		if v,ok := langs[name]; ok {
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
func (global *constGlobal) Check(value, regular string) bool {

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



//按状态生成错误
//args是用来Format错误信息的
//state 和langstring 联动
func (global *constGlobal) NewTypeLangStateError(tttt, lang, state string, args ...interface{}) *Error {

	stateCode := global.StateCode(state)
	stateText := global.LangString(state, lang)

	return NewTypeCodeError(tttt, stateCode, fmt.Sprintf(stateText, args...))
}



//按状态生成错误
//args是用来Format错误信息的
//state 和langstring 联动
func (global *constGlobal) NewTypeStateError(tttt,state string, args ...interface{}) *Error {

	stateCode := global.StateCode(state)
	stateText := global.LangString(state)

	return NewTypeCodeError(tttt,stateCode, fmt.Sprintf(stateText, args...))
}
