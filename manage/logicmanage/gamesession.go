package logicmanage

import (
	"net"

	"github.com/yhhaiua/goserver/common"
	"github.com/yhhaiua/goserver/common/glog"
	"github.com/yhhaiua/goserver/common/gpacket"
	"github.com/yhhaiua/goserver/common/gtcp"
	"github.com/yhhaiua/goserver/comsvrsrc"
	"github.com/yhhaiua/goserver/protocol"
)

type stGameSession struct {
	*gtcp.ServerSession
}

//create 创建连接
func (session *stGameSession) create(con *net.TCPConn, linkKey int64) bool {
	session.ServerSession = gtcp.AddSession(con, linkKey, "game服务器")

	if session.ServerSession != nil {
		session.SetFunc(session.putMsgQueue, session.delCloseLink)
		session.Start()
		return true
	}
	return false
}

//putMsgQueue 消息队列
func (session *stGameSession) putMsgQueue(pcmd *gpacket.BaseCmd, data []byte) bool {
	switch pcmd.Value() {
	case protocol.ServerCmdLoginCode:
		return session.loginCmd(data)
	case protocol.ServerCmdHeartCode:
		session.heartCmd(data)
	default:
	}
	return true
}

//delCloseLink 断开连接回调
func (session *stGameSession) delCloseLink(servertag int64) {

	Instance().syncgameMap().Delete(servertag)

}

//sendOnceCmd 连接发送验证包
func (session *stGameSession) sendOnceCmd() {
	var retcmd protocol.ServerCmdLogin
	retcmd.Init()

	retcmd.CheckData = comsvrsrc.CHECKDATACODE
	retcmd.Svrid = Instance().serverid
	retcmd.Svrtype = SERVERTYPE

	session.SendCmd(&retcmd)
}

//loginCmd 收到验证包
func (session *stGameSession) loginCmd(data []byte) bool {
	var retcmd protocol.ServerCmdLogin

	err := session.Cmdcodec().Decode(data, &retcmd)

	if common.CheckError(err, "ServerCmdLogin") && retcmd.CheckData == comsvrsrc.CHECKDATACODE {
		session.SetValid(true)
		session.sendOnceCmd()
		glog.Infof("game %d-%d 连接效验成功", retcmd.Svrid, retcmd.Svrtype)
		return true
	}
	return false
}

//heartCmd 心跳包
func (session *stGameSession) heartCmd(data []byte) bool {
	var retcmd protocol.ServerCmdHeart

	err := session.Cmdcodec().Decode(data, &retcmd)

	if common.CheckError(err, "ServerCmdHeart") {
		session.Setheartbeat(&retcmd)
		return true
	}
	return false
}