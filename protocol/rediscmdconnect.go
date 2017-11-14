package protocol

import "github.com/yhhaiua/goserver/common/gpacket"

//RedisCmdConnect redis间连接
type RedisCmdConnect struct {
	gpacket.BaseCmd
	IsneedAck bool   //是否需要回包
	Szchannel string //通道名
}

//Init RedisCmdConnect初始化
func (pcmd *RedisCmdConnect) Init() {
	pcmd.Cmd = 253
	pcmd.SupCmd = 1
}
