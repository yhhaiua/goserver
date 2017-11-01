package protocol

import "github.com/yhhaiua/goserver/common"

//ServerCmdHeart 服务器间心跳包
type ServerCmdHeart struct {
	common.BaseCmd
	IsneedAck bool
	Checknum  int8
}

//Init ServerCmdHeart初始化
func (pcmd *ServerCmdHeart) Init() {
	pcmd.Cmd = 255
	pcmd.SupCmd = 254
}
