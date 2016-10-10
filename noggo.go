package noggo


type (
	Noggo struct {
		//节点名称和唯一标识
		Id		string
		Name	string

		Plan	*planModule
	}
)



//创建新节点
func New(name string) (*Noggo) {


	return &Noggo{}
}




//启动节点
func (node *Noggo) Run() {

}