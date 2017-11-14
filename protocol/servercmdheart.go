package protocol

import "github.com/yhhaiua/goserver/common/gpacket"

//ServerCmdHeart 服务器间心跳包
type ServerCmdHeart struct {
	gpacket.BaseCmd
	IsneedAck bool //是否需要回包
	Checknum  int8 //检测次数
}

//Init ServerCmdHeart初始化
func (pcmd *ServerCmdHeart) Init() {
	pcmd.Cmd = 254
	pcmd.SupCmd = 2
}
