package base

/*
    用于生成查询条件，示例如下：

    Map{
        "views":	Map{ GT: 100, LT: 500 },			解析成：	views>100 AND views<500
        "hits":		Map{ GTE: 100, LTE: 500 },			解析成	hits>=100 AND hits<=500
        "name": 	Map{ EQ: "noggo", NEQ: "nogio" },	解析成：	name='noggo' AND name!='nogio'
        "tags": 	Map{ ANY: "nog" },					解析成：	ANY(tags)='nog'
        "id":		Map{ IN: []int{ 1,2,3 } },			解析成：	id IN (1,2,3)
        "email": 	nil,								解析成：	email IS NULL
        "email": 	NIL,								解析成：	email IS NULL
        "email": 	NOL,								解析成：	email IS NOT NULL

        "id":		ASC									解析成：id ASC
        "id":		DESC								解析成：id DESC
    }

*/

const (
	DELIMS	= "\t"	//字段以及表名边界符，自己实现数据驱动才需要处理这个，必须能启标识作用

	IS	= "="	//等于
	EQ	= "="	//等于
	NE	= "!="	//不等于
	NEQ	= "!="	//不等于
	NOT = "!="	//不等于

	GT	= ">"	//大于
	GE	= ">="	//大于等于
	GTE	= ">="	//大于等于
	LT	= "<"	//小于
	LE	= "<="	//小于等于
	LTE	= "<="	//小于等于

	IN	= "$$$IN$$$"	//支持  WHERE id IN (1,2,3)			//这条还没支持
	ANY = "$$$ANY"		//支持数组字段的   xxx=ANY("field")	//这条还没支持

	FULL = "$$$full$$$"		//like搜索
	LEFT = "$$$left$$$"		//like left搜索
	RIGHT = "$$$right$$$"	//like right搜索

	INC	= "$$$inc$$$"	//累加，    UPDATE时用，解析成：views=views+1

)

type (
	dataNil		struct {}
	dataNol		struct {}
	dataAsc		struct {}
	dataDesc	struct {}
	LIKE		string
)

var (
	NIL		dataNil		//为空	IS NULL
	NULL	dataNil		//为空	IS NULL
	NOL		dataNol		//不为空	IS NOT NULL
	NOTNULL	dataNol		//不为空	IS NOT NULL
	ASC		dataAsc		//正序	asc
	DESC	dataDesc	//倒序	desc
)




