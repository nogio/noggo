/*
	全局包
	此包其实没什么用，只是用来加载其它包
	当然，也可以直接在runs/noggo.go中直接加载以下包
*/


package globals


import (
	//_ "./https"
	//_ "./plans"
	_ "./tasks"
	_ "./triggers"
)
