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
	DELIMS	= `"`	//字段以及表名边界符，自己实现数据驱动才需要处理这个，必须能启标识作用

	IS		= "="	//等于
	NOT 	= "!="	//不等于
	EQ		= "="	//等于
	NE		= "!="	//不等于
	NEQ		= "!="	//不等于

	//约等于	正则等于
	AE		= "~*"		//正则等于，约等于
	AEC		= "~"		//正则等于，区分大小写，
	RE		= "~*"		//正则等于，约等于
	REC		= "~"		//正则等于，区分大小写，
	REQ		= "~*"		//正则等于，约等于
	REQC	= "~"		//正则等于，区分大小写，

	NAE		= "!~*"		//正则不等于，
	NAEC	= "!~"		//正则不等于，区分大小写，
	NRE		= "!~*"		//正则不等于，
	NREC	= "!~"		//正则不等于，区分大小写，
	NREQ	= "!~*"		//正则不等于，
	NREQC	= "!~"		//正则不等于，区分大小写，

	//换位约等于，值在前，字段在后，用于黑名单查询
	EA		= "$$~*$$"		//正则等于，约等于
	EAC		= "$$~$$"		//正则等于，区分大小写，
	ER		= "$$~*$$"		//正则等于，约等于
	ERC		= "$$~$$"		//正则等于，区分大小写，
	EQR		= "$$~*$$"		//正则等于，约等于
	EQRC	= "$$~$$"		//正则等于，区分大小写，

	NEA		= "$$!~*$$"		//正则不等于，
	NEAC	= "$$!~$$"		//正则不等于，区分大小写，
	NER		= "$$!~*$$"		//正则不等于，
	NERC	= "$$!~$$"		//正则不等于，区分大小写，
	NEQR	= "$$!~*$$"		//正则不等于，
	NEQRC	= "$$!~$$"		//正则不等于，区分大小写，


	GT	= ">"	//大于
	GE	= ">="	//大于等于
	GTE	= ">="	//大于等于
	LT	= "<"	//小于
	LE	= "<="	//小于等于
	LTE	= "<="	//小于等于

	IN	= "$$IN$$"	//支持  WHERE id IN (1,2,3)			//这条还没支持
	NIN = "$$NOTIN$$"	//支持	WHERE id NOT IN(1,2,3)
	ANY = "$$ANY"		//支持数组字段的   xxx=ANY("field")	//这条还没支持

	LIKE = "$$full$$"		//like搜索
	FULL = "$$full$$"		//like搜索
	LEFT = "$$left$$"		//like left搜索
	RIGHT = "$$right$$"	//like right搜索

	INC	= "$$inc$$"	//累加，    UPDATE时用，解析成：views=views+value

)

type (
	dataNil		struct {}
	dataNol		struct {}
	dataAsc		struct {}
	dataDesc	struct {}
)

var (
	NIL		dataNil		//为空	IS NULL
	NULL	dataNil		//为空	IS NULL
	NOL		dataNol		//不为空	IS NOT NULL
	NOTNULL	dataNol		//不为空	IS NOT NULL
	ASC		dataAsc		//正序	asc
	DESC	dataDesc	//倒序	desc
)




