package base




type (
	Error struct {
		Type	Type		//错误类型
		Code	int64		//错误编码
		Text	string		//错误信息文本
	}
)



/*
	创建错误对象
*/
func NewError(code int64, text string) (*Error) {
	return &Error{ TypeErrorNone, code, text }
}
func NewTypeError(tttt Type, code int64, text string) (*Error) {
	return &Error{ tttt, code, text }
}




/*
	创建类型错误
func NewTypeError(tttt ErrorType, code int64, text string) (*Error) {
	return &Error{ tttt, code, text }
}
*/