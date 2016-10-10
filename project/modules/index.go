/*
	模块包
	此包其实没什么用，只是用来加载其它包
	当然，也可以直接在runs/noggo.go中直接加载以下包

	在此目录下所加载的所有模块相关的内容
	都相当于加载到所有节点

	比如，有以下节点：
	ing, admin, www
	那么在此目录下加载的所有计划，WEB，队列等等
	nog := noggo.New("ing") 的时候， 就自动加载到此节点实例
*/

package modules

import (
	//_ "./https"
	//_ "./plans"
)
