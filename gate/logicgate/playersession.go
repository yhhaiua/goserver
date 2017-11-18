package logicgate

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

type stPlayerSession struct {
	ginter.NetWorker
	codec common.BinaryCodec
}

//create 创建连接
func (session *stPlayerSession) create(con *net.TCPConn, linkKey int64) bool {
	session.NetWorker = gtcp.AddSession(con, linkKey, "playersession", session)

	if session.NetWorker != nil {
		session.Start()
		return true
	}
	return false
}

//MsgQueue 消息队列
func (session *stPlayerSession) MsgQueue(pcmd *gpacket.BaseCmd, data []byte) bool {
	switch pcmd.Value() {
	case protocol.ServerCmdLoginCode:
		return session.loginCmd(data)
	default:
	}
	return true
}

//CloseLink 断开连接回调
func (session *stPlayerSession) CloseLink(servertag int64) {

	Instance().syncMap().Delete(servertag)

}

//StartLink 启动回调
func (session *stPlayerSession) StartLink(servertag int64) {

}

//CmdCodec 解析函数
func (session *stPlayerSession) CmdCodec() common.CmdCodec {
	return &session.codec
}

//loginCmd 收到验证包
func (session *stPlayerSession) loginCmd(data []byte) bool {
	var retcmd protocol.ServerCmdLogin

	err := session.CmdCodec().Decode(data, &retcmd)

	if common.CheckError(err, "ServerCmdLogin") && retcmd.CheckData == comsvrsrc.CHECKDATACODE {
		session.SetValid(true)
		glog.Infof("player %d-%d 连接效验成功", retcmd.Svrid, retcmd.Svrtype)
		return true
	}
	return false
}
