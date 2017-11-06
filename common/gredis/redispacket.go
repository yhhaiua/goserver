package gredis

import (
	"github.com/yhhaiua/goserver/common/glog"
	"github.com/yhhaiua/goserver/common/goobjfmt"
	"github.com/yhhaiua/goserver/common/gpacket"
)

//RedisPacket use
type RedisPacket struct {
	*RedisPool
}

//SendCmd 向redis频道发送数据
func (rc *RedisPacket) SendCmd(sChannel string, data interface{}) {

	var packet gpacket.Packet
	packet.Size = uint32(goobjfmt.BinarySize(data))
	packet.Data = data
	bytedata, err := goobjfmt.BinaryWrite(&packet)
	if err != nil {
		glog.Errorf("data err:%s", err)
		return
	}
	rc.Publish(sChannel, bytedata)
}

//NewRedis redis创建
func NewRedis(config *RedisConfig) (adapter *RedisPacket, err error) {

	adapter.RedisPool = newRedis()
	err = adapter.start(config)
	if err != nil {
		adapter = nil
	}
	return
}
