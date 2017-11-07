package protocol

import (
	"github.com/yhhaiua/goserver/common/gpacket"
)

//ServerCmdCon redis间连接
type ServerCmdCon struct {
	gpacket.BaseCmd
	IsneedAck bool
	Szchannel string
}

//Init ServerCmdCon初始化
func (pcmd *ServerCmdCon) Init() {
	pcmd.Cmd = 254
	pcmd.SupCmd = 3
}
