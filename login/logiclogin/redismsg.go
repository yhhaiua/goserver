package logiclogin

import (
	"github.com/yhhaiua/goserver/common"
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
	case protocol.RedisCmdConnectCode:
		msg.logincmd(data)
	default:
	}
	return true
}
func (msg *stRedisMsg) logincmd(data []byte) {
	var retcmd protocol.RedisCmdConnect

	err := Instance().redisdb().Cmdcodec().Decode(data, &retcmd)

	if common.CheckError(err, "RedisCmdConnect") {
		getchannel := retcmd.Szchannel
		glog.Infof("连接成功 %s", retcmd.Szchannel)
		msg.boConnection = true
		retcmd.Szchannel = comsvrsrc.SUBCHANNELlogin
		Instance().redisdb().SendCmd(getchannel, &retcmd)
	}
}

func (msg *stRedisMsg) boCon() bool {
	return msg.boConnection
}
