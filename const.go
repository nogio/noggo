package noggo

import (
	. "github.com/nogio/noggo/base"
)


type (
	constModule struct {
		mimes		map[string]string
		states		map[string]int
		regulars	map[string][]string
		strings		map[string]map[string]string
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

//注册String
func (module *constModule) String(lang string, strings Map) {
	if module.strings == nil {
		module.strings = map[string]map[string]string{}
	}
	if module.strings[lang] == nil {
		module.strings[lang] = map[string]string{}
	}


	for k,v := range strings {
		switch vvvv := v.(type) {
		case string:
			module.strings[lang][k] = vvvv
		}
	}
}