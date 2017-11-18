package protocol

import "github.com/yhhaiua/goserver/common/gpacket"

//ServerCmdConData 服务器间连接数据包
type ServerCmdConData struct {
	gpacket.BaseCmd
	ConDataInfo []RefConDataInfo //连接数据
}

//Init ServerCmdConData初始化
func (pcmd *ServerCmdConData) Init() {
	pcmd.Cmd = 254
	pcmd.SupCmd = 3
}
