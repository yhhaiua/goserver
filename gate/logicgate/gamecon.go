package logicgate

import (
	"github.com/yhhaiua/goserver/common"
	"github.com/yhhaiua/goserver/common/ginter"
	"github.com/yhhaiua/goserver/common/glog"
	"github.com/yhhaiua/goserver/common/gpacket"
	"github.com/yhhaiua/goserver/common/gtcp"
	"github.com/yhhaiua/goserver/comsvrsrc"
	"github.com/yhhaiua/goserver/protocol"
)

type stGameCon struct {
	ginter.NetWorker
	codec common.BinaryCodec
}

//create 创建连接
func (con *stGameCon) create(sip, sport string, serverid int32) bool {
	con.NetWorker = gtcp.AddConnect(sip, sport, serverid, "game服务器", con)

	if con.NetWorker != nil {
		con.Start()
		return true
	}
	return false
}

//MsgQueue 消息队列
func (con *stGameCon) MsgQueue(pcmd *gpacket.BaseCmd, data []byte) bool {
	switch pcmd.Value() {
	case protocol.ServerCmdLoginCode:
		return con.loginCmd(data)
	case protocol.ServerCmdHeartCode:
		con.heartCmd(data)
	default:
	}
	return true
}

//CloseLink 关闭回调
func (con *stGameCon) CloseLink(servertag int64) {

}

//StartLink 启动回调
func (con *stGameCon) StartLink(servertag int64) {

}

//CmdCodec 解析函数
func (con *stGameCon) CmdCodec() common.CmdCodec {
	return &con.codec
}

//SendOnceCmd 连接发送验证包
func (con *stGameCon) SendOnceCmd() {
	var retcmd protocol.ServerCmdLogin
	retcmd.Init()

	retcmd.CheckData = comsvrsrc.CHECKDATACODE
	retcmd.Svrid = Instance().serverid
	retcmd.Svrtype = SERVERTYPE
	retcmd.Sip = Instance().config().sip
	retcmd.Sport = Instance().config().sport
	con.SendCmd(&retcmd)
}

//loginCmd 收到验证返回包
func (con *stGameCon) loginCmd(data []byte) bool {
	var retcmd protocol.ServerCmdLogin

	err := con.CmdCodec().Decode(data, &retcmd)

	if common.CheckError(err, "ServerCmdLogin") && retcmd.CheckData == comsvrsrc.CHECKDATACODE {
		con.SetValid(true)
		glog.Infof("game服务器 %d-%d 连接效验成功ip:[%s],port:[%s]", retcmd.Svrid, retcmd.Svrtype, retcmd.Sip, retcmd.Sport)
		return true
	}
	return false
}

//heartCmd 心跳包
func (con *stGameCon) heartCmd(data []byte) bool {
	var retcmd protocol.ServerCmdHeart

	err := con.CmdCodec().Decode(data, &retcmd)

	if common.CheckError(err, "ServerCmdHeart") {
		con.Setheartbeat(&retcmd)
		return true
	}
	return false
}
