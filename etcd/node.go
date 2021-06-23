package etcd

type INodeInfo interface {
	GetRegisterInfo() string //获取节点的注册信息
	SetRegisterInfo(string)  //设置节点的注册信息
	GetUuid() string         //UUID
	GetStatus() string       //节点状态
	GetRegistTime() int64    //注册时间
	GetLastTimes() int64     //最后心跳时间
}

type NodeInfo struct {
	Uuid         string //UUID
	RegisterInfo string //UUID  watch的groupname
	Online       uint32
	RegistTime   int64 //注册时间
	LastTime     int64 //最后心跳时间
	status       func() string
}

func (this *NodeInfo) GetRegistTime() int64 {
	return this.RegistTime

}
func (this *NodeInfo) GetLastTimes() int64 {
	return this.LastTime
}
func (this *NodeInfo) GetUuid() string {
	return this.Uuid
}

func (this *NodeInfo) GetRegisterInfo() string {
	return this.RegisterInfo
}

func (this *NodeInfo) SetRegisterInfo(regInfo string) {
	this.RegisterInfo = regInfo
}

func (this *NodeInfo) GetStatus() string {
	return this.status()
}

func (this *NodeInfo) SetStatus(f func() string) {
	this.status = f
}
