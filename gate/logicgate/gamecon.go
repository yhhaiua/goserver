//读取配置文件

package logicgate

import (
	"github.com/yhhaiua/goserver/common"
	"github.com/yhhaiua/goserver/common/glog"
	"github.com/yhhaiua/goserver/common/gtcp"
	"github.com/yhhaiua/goserver/protocol"
)

type stGameCon struct {
	*gtcp.ClientConnecter
}

func (con *stGameCon) create() bool {
	con.ClientConnecter = gtcp.AddConnect("172.16.3.141", Instance().config().sport, Instance().serverid)

	if con.ClientConnecter != nil {
		con.SetFunc(con.putMsgQueue, con.sendOnceCmd)
		con.Start()
		return true
	}
	return false
}

func (con *stGameCon) putMsgQueue(pcmd *common.BaseCmd, data []byte) bool {
	switch pcmd.Value() {
	case protocol.ServerCmdLoginValue():
		return con.loginCmd(data)
	default:
	}
	return true
}

func (con *stGameCon) sendOnceCmd() {
	var retcmd protocol.ServerCmdLogin
	retcmd.Init()
	retcmd.Svrid = Instance().serverid
	retcmd.Svrtype = SERVERTYPE

	con.SendCmd(&retcmd)
}

func (con *stGameCon) loginCmd(data []byte) bool {
	var retcmd protocol.ServerCmdLogin
	err := con.Cmdcodec.Decode(data, &retcmd)

	if common.CheckError(err, "ServerCmdLogin") && retcmd.CheckData == protocol.CHECKDATACODE {
		con.SetValid(true)
		glog.Infof("game服务器 %d-%d 连接效验成功", retcmd.Svrid, retcmd.Svrtype)
		return true
	}
	return false
}
