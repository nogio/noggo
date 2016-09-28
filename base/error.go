package base




type (
	Error struct {
		Type	string		//错误类型
		Code	int			//错误编码
		Text	string		//错误信息文本
	}
)



/*
	创建错误对象
*/
func NewError(text string) (*Error) {
	return &Error{ "", 0, text }
}
func NewCodeError(code int, text string) (*Error) {
	return &Error{ "", code, text }
}
func NewTypeError(tttt string, text string) (*Error) {
	return &Error{ tttt, 0, text }
}
func NewTypeCodeError(tttt string, code int, text string) (*Error) {
	return &Error{ tttt, 0, text }
}
