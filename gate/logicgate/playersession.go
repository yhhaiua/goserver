package logicgate

import (
	"net"

	"github.com/yhhaiua/goserver/common"
	"github.com/yhhaiua/goserver/common/glog"
	"github.com/yhhaiua/goserver/common/gpacket"
	"github.com/yhhaiua/goserver/common/gtcp"
	"github.com/yhhaiua/goserver/comsvrsrc"
	"github.com/yhhaiua/goserver/protocol"
)

type stPlayerSession struct {
	*gtcp.ServerSession
}

//create 创建连接
func (session *stPlayerSession) create(con *net.TCPConn, linkKey int64) bool {
	session.ServerSession = gtcp.AddSession(con, linkKey, "playersession")

	if session.ServerSession != nil {
		session.SetFunc(session.putMsgQueue, session.delCloseLink)
		session.Start()
		return true
	}
	return false
}

//putMsgQueue 消息队列
func (session *stPlayerSession) putMsgQueue(pcmd *gpacket.BaseCmd, data []byte) bool {
	switch pcmd.Value() {
	case protocol.ServerCmdLoginCode:
		return session.loginCmd(data)
	default:
	}
	return true
}

//delCloseLink 断开连接回调
func (session *stPlayerSession) delCloseLink(servertag int64) {

	Instance().syncMap().Delete(servertag)

}

//loginCmd 收到验证包
func (session *stPlayerSession) loginCmd(data []byte) bool {
	var retcmd protocol.ServerCmdLogin

	err := session.Cmdcodec().Decode(data, &retcmd)

	if common.CheckError(err, "ServerCmdLogin") && retcmd.CheckData == comsvrsrc.CHECKDATACODE {
		session.SetValid(true)
		glog.Infof("player %d-%d 连接效验成功", retcmd.Svrid, retcmd.Svrtype)
		return true
	}
	return false
}
