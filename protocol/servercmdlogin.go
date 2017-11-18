package protocol

import "github.com/yhhaiua/goserver/common/gpacket"

//ServerCmdLogin 服务器间登录包
type ServerCmdLogin struct {
	gpacket.BaseCmd
	CheckData uint32 //效验码
	Svrid     int32  //服务器id
	Svrtype   int32  //服务器类型
	Sip       string //服务器ip
	Sport     string //服务器端口
}

//Init ServerCmdLogin初始化
func (pcmd *ServerCmdLogin) Init() {
	pcmd.Cmd = 254
	pcmd.SupCmd = 1
}
