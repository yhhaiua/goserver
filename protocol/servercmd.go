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

//ServerCmdHeart 服务器间心跳包
type ServerCmdHeart struct {
	common.BaseCmd
	IsneedAck bool
	Checknum  int8
}

//Init ServerCmdLogin初始化
func (pcmd *ServerCmdLogin) Init() {
	pcmd.Cmd = 254
	pcmd.SupCmd = 1
}

//ServerCmdLoginValue ServerCmdLogin的Value值
func ServerCmdLoginValue() uint16 {
	return common.GetValue(254, 1)
}

//Init ServerCmdHeart初始化
func (pcmd *ServerCmdHeart) Init() {
	pcmd.Cmd = 254
	pcmd.SupCmd = 254
}

//ServerCmdHeartValue ServerCmdHeart的Value值
func ServerCmdHeartValue() uint16 {
	return common.GetValue(254, 254)
}
