package logicduty

import (
	"time"

	"github.com/yhhaiua/goserver/common/glog"
	"github.com/yhhaiua/goserver/common/gpacket"
	"github.com/yhhaiua/goserver/comsvrsrc"
	"github.com/yhhaiua/goserver/protocol"
)

type stRedisMsg struct {
	boConnection bool
}

func (msg *stRedisMsg) putMsgQueue(pcmd *gpacket.BaseCmd, data []byte) bool {

	switch pcmd.Value() {
	case protocol.ServerCmdConCode:
		msg.boConnection = true
		glog.Info("login服务器连接成功")
	default:
	}
	return true
}

func (msg *stRedisMsg) runSendLogin() {
	for {
		if msg.boConnection {
			break
		}
		msg.sendlogincmd()
		time.Sleep(5 * time.Second)
	}
}

func (msg *stRedisMsg) sendlogincmd() {
	var retcmd protocol.ServerCmdCon
	retcmd.Init()
	retcmd.IsneedAck = true
	retcmd.Szchannel = comsvrsrc.SUBCHANNELduty
	Instance().redisdb().SendCmd(comsvrsrc.SUBCHANNELlogin, &retcmd)
	glog.Infof("发送连接请求 %s", comsvrsrc.SUBCHANNELlogin)
}
