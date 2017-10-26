package protocol

import "github.com/yhhaiua/goserver/common"

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
	pcmd.CheckData = 0x55884433
	pcmd.Cmd = 254
	pcmd.SupCmd = 1
}
