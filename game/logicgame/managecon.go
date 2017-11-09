package logicgame

import (
	"github.com/yhhaiua/goserver/common"
	"github.com/yhhaiua/goserver/common/glog"
	"github.com/yhhaiua/goserver/common/gpacket"
	"github.com/yhhaiua/goserver/common/gtcp"
	"github.com/yhhaiua/goserver/comsvrsrc"
	"github.com/yhhaiua/goserver/protocol"
)

type stManageCon struct {
	*gtcp.ClientConnecter
}

//create 创建连接
func (con *stManageCon) create(game *stManageConfig) bool {
	con.ClientConnecter = gtcp.AddConnect(game.sip, game.sport, game.serverid, "manage务器")

	if con.ClientConnecter != nil {
		con.SetFunc(con.putMsgQueue, con.sendOnceCmd)
		con.Start()
		return true
	}
	return false
}

//putMsgQueue 消息队列
func (con *stManageCon) putMsgQueue(pcmd *gpacket.BaseCmd, data []byte) bool {
	switch pcmd.Value() {
	case protocol.ServerCmdLoginCode:
		return con.loginCmd(data)
	case protocol.ServerCmdHeartCode:
		con.heartCmd(data)
	default:
	}
	return true
}

//sendOnceCmd 连接发送验证包
func (con *stManageCon) sendOnceCmd() {
	var retcmd protocol.ServerCmdLogin
	retcmd.Init()

	retcmd.CheckData = comsvrsrc.CHECKDATACODE
	retcmd.Svrid = Instance().serverid
	retcmd.Svrtype = SERVERTYPE

	con.SendCmd(&retcmd)
}

//loginCmd 收到验证返回包
func (con *stManageCon) loginCmd(data []byte) bool {
	var retcmd protocol.ServerCmdLogin

	err := con.Cmdcodec().Decode(data, &retcmd)

	if common.CheckError(err, "ServerCmdLogin") && retcmd.CheckData == comsvrsrc.CHECKDATACODE {
		con.SetValid(true)
		glog.Infof("manage服务器 %d-%d 连接效验成功", retcmd.Svrid, retcmd.Svrtype)
		return true
	}
	return false
}

//heartCmd 心跳包
func (con *stManageCon) heartCmd(data []byte) bool {
	var retcmd protocol.ServerCmdHeart

	err := con.Cmdcodec().Decode(data, &retcmd)

	if common.CheckError(err, "ServerCmdHeart") {
		con.Setheartbeat(&retcmd)
		return true
	}
	return false
}
