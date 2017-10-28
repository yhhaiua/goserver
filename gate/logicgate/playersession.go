package logicgate

import (
	"net"
	"sync"

	"github.com/yhhaiua/goserver/common"
	"github.com/yhhaiua/goserver/common/glog"
	"github.com/yhhaiua/goserver/common/gtcp"
	"github.com/yhhaiua/goserver/protocol"
)

type stPlayerSessMap struct {
	sync.Mutex
	sessMap map[int64]*stPlayerSession
}

func (sess *stPlayerSessMap) add(newPlayer *stPlayerSession, key int64) {
	sess.Lock()
	defer sess.Unlock()
	sess.sessMap[key] = newPlayer
}
func (sess *stPlayerSessMap) delete(key int64) {
	sess.Lock()
	defer sess.Unlock()
	delete(sess.sessMap, key)
}
func (sess *stPlayerSessMap) SendCmd(key int64, data interface{}) {
	sess.Lock()
	defer sess.Unlock()
	value, ok := sess.sessMap[key]
	if ok {
		value.SendCmd(data)
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////////////
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
func (session *stPlayerSession) putMsgQueue(pcmd *common.BaseCmd, data []byte) bool {
	switch pcmd.Value() {
	case protocol.ServerCmdLoginValue():
		return session.loginCmd(data)
	default:
	}
	return true
}

//delCloseLink 断开连接回调
func (session *stPlayerSession) delCloseLink(servertag int64) {

	Instance().playerMap().delete(servertag)

}

//loginCmd 收到验证包
func (session *stPlayerSession) loginCmd(data []byte) bool {
	var retcmd protocol.ServerCmdLogin

	err := session.Cmdcodec().Decode(data, &retcmd)

	if common.CheckError(err, "ServerCmdLogin") && retcmd.CheckData == protocol.CHECKDATACODE {
		session.SetValid(true)
		glog.Infof("player %d-%d 连接效验成功", retcmd.Svrid, retcmd.Svrtype)
		return true
	}
	return false
}
