package logicgame

import (
	"net"

	"github.com/yhhaiua/goserver/common"
	"github.com/yhhaiua/goserver/common/ginter"
	"github.com/yhhaiua/goserver/common/glog"
	"github.com/yhhaiua/goserver/common/gpacket"
	"github.com/yhhaiua/goserver/common/gtcp"
	"github.com/yhhaiua/goserver/comsvrsrc"
	"github.com/yhhaiua/goserver/protocol"
)

type stGateSession struct {
	ginter.NetWorker
	codec common.BinaryCodec
}

//create 创建连接
func (session *stGateSession) create(con *net.TCPConn, linkKey int64) bool {
	session.NetWorker = gtcp.AddSession(con, linkKey, "gate服务器", session)

	if session.NetWorker != nil {
		session.Start()
		return true
	}
	return false
}

//MsgQueue 消息队列
func (session *stGateSession) MsgQueue(pcmd *gpacket.BaseCmd, data []byte) bool {
	switch pcmd.Value() {
	case protocol.ServerCmdLoginCode:
		return session.loginCmd(data)
	case protocol.ServerCmdHeartCode:
		session.heartCmd(data)
	default:
	}
	return true
}

//CloseLink 断开连接回调
func (session *stGateSession) CloseLink(servertag int64) {

	Instance().syncMap().Delete(servertag)

}

//StartLink 启动回调
func (session *stGateSession) StartLink(servertag int64) {

}

//CmdCodec 解析函数
func (session *stGateSession) CmdCodec() common.CmdCodec {
	return &session.codec
}

//sendOnceCmd 连接发送验证包
func (session *stGateSession) sendOnceCmd() {
	var retcmd protocol.ServerCmdLogin
	retcmd.Init()

	retcmd.CheckData = comsvrsrc.CHECKDATACODE
	retcmd.Svrid = Instance().serverid
	retcmd.Svrtype = SERVERTYPE

	session.SendCmd(&retcmd)
}

//loginCmd 收到验证包
func (session *stGateSession) loginCmd(data []byte) bool {
	var retcmd protocol.ServerCmdLogin

	err := session.CmdCodec().Decode(data, &retcmd)

	if common.CheckError(err, "ServerCmdLogin") && retcmd.CheckData == comsvrsrc.CHECKDATACODE {
		session.SetValid(true)
		session.sendOnceCmd()
		glog.Infof("gate %d-%d 连接效验成功", retcmd.Svrid, retcmd.Svrtype)
		return true
	}
	return false
}

//heartCmd 心跳包
func (session *stGateSession) heartCmd(data []byte) bool {
	var retcmd protocol.ServerCmdHeart

	err := session.CmdCodec().Decode(data, &retcmd)

	if common.CheckError(err, "ServerCmdHeart") {
		session.Setheartbeat(&retcmd)
		return true
	}
	return false
}
