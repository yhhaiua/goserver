package protocol

import "github.com/yhhaiua/goserver/common/gpacket"

//ServerCmdHeart 服务器间心跳包
type ServerCmdHeart struct {
	gpacket.BaseCmd
	IsneedAck bool
	Checknum  int8
}

//Init ServerCmdHeart初始化
func (pcmd *ServerCmdHeart) Init() {
	pcmd.Cmd = 254
	pcmd.SupCmd = 2
}
