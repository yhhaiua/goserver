package protocol

import "github.com/yhhaiua/goserver/common/gpacket"

//LoginCmdUserRet 玩家登录包返回（S->C）
type LoginCmdUserRet struct {
	gpacket.BaseCmd
	ErrorCode int32 //错误码
}

//Init LoginCmdUserRet初始化
func (pcmd *LoginCmdUserRet) Init() {
	pcmd.Cmd = 1
	pcmd.SupCmd = 2
}
