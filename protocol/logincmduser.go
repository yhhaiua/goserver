package protocol

import "github.com/yhhaiua/goserver/common/gpacket"

//LoginCmdUser 玩家登录包（C->S）
type LoginCmdUser struct {
	gpacket.BaseCmd
	Account string //帐号
	Name    string //名字
	OnlyID  int64  //玩家id
	Paramp  string //md5加密字符串
	Paramt  string //时间戳
}

//Init LoginCmdUser初始化
func (pcmd *LoginCmdUser) Init() {
	pcmd.Cmd = 1
	pcmd.SupCmd = 1
}
