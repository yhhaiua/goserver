package protocol

import "github.com/yhhaiua/goserver/common"

const (
	//CHECKDATACODE 服务器间数据检测
	CHECKDATACODE = 0x55884433
)

//ServerCmdLogin 服务器间登录包
type ServerCmdLogin struct {
	common.BaseCmd
	CheckData uint32
	Svrid     int32
	Svrtype   int32
	Now       int32
}

//Init ServerCmdLogin初始化
func (pcmd *ServerCmdLogin) Init() {
	pcmd.CheckData = CHECKDATACODE
	pcmd.Cmd = 254
	pcmd.SupCmd = 1
}

//ServerCmdLoginValue ServerCmdLogin的Value值
func ServerCmdLoginValue() uint16 {
	return common.GetValue(254, 1)
}
