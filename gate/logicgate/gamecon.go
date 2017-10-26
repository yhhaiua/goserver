//读取配置文件

package logicgate

import (
	"github.com/yhhaiua/goserver/common"
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
		return true
	}
	return false
}

func (con *stGameCon) putMsgQueue(pcmd *common.BaseCmd, data []byte) {

}

func (con *stGameCon) sendOnceCmd() {
	var retcmd protocol.ServerCmdLogin
	retcmd.Init()
	retcmd.Svrid = Instance().serverid
	retcmd.Svrtype = SERVERTYPE

	con.SendCmd(retcmd)
}
