package gredis

import (
	"github.com/garyburd/redigo/redis"
	"github.com/yhhaiua/goserver/common"
	"github.com/yhhaiua/goserver/common/glog"
	"github.com/yhhaiua/goserver/common/gpacket"
)

//RedisPacket use
type RedisPacket struct {
	*RedisPool
	cmdcodec common.CmdCodec
	msgQueue func(pcmd *gpacket.BaseCmd, data []byte) bool
}

//SendCmd 向redis频道发送数据
func (rc *RedisPacket) SendCmd(sChannel string, data interface{}) {

	var packet gpacket.Packet
	packet.Size = uint32(rc.cmdcodec.Size(data))
	packet.Data = data
	bytedata, err := rc.cmdcodec.Encode(&packet)
	if err != nil {
		glog.Errorf("data err:%s", err)
		return
	}
	rc.Publish(sChannel, bytedata)
}

//Publish 发布
func (rc *RedisPacket) Publish(sChannel string, value interface{}) error {
	var err error
	if _, err = rc.do("PUBLISH", sChannel, value); err != nil {
		return err
	}
	return err
}

//Subscribe 订阅
func (rc *RedisPacket) Subscribe(sChannel string) {
	go rc.subscribe(sChannel)
}
func (rc *RedisPacket) subscribe(sChannel string) {
	c := rc.p.Get()
	defer c.Close()
	psc := new(redis.PubSubConn)
	psc.Conn = c
	psc.Subscribe(sChannel)
	for {
		switch v := psc.Receive().(type) {
		case redis.Subscription:
			glog.Infof("%s: %s %d\n", v.Channel, v.Kind, v.Count)
		case redis.Message: //单个订阅subscribe
			glog.Infof("%s: message: %s\n", v.Channel, v.Data)
			rc.doRead(v.Data)
		case redis.PMessage: //模式订阅psubscribe
			glog.Infof("PMessage: %s %s %s\n", v.Pattern, v.Channel, v.Data)
		case error:
			return

		}

	}
}

func (rc *RedisPacket) doRead(data []byte) {
	datalen := len(data)
	if datalen >= 8 {
		var packet gpacket.PacketBase

		err := rc.cmdcodec.Decode(data[:8], &packet)

		if err != nil {
			glog.Errorf("收到恶意攻击包%s", err)
			return
		}
		if packet.Size >= 1024*64 || packet.Size < 2 {
			glog.Errorf("收到恶意攻击包 %d", packet.Size)
			return
		}
		newlen := int(packet.Size + 6)
		if datalen == newlen {
			//包处理
			if rc.msgQueue(&packet.Pcmd, data[6:packet.Size+6]) {

			} else {
				return
			}
		} else {
			glog.Errorf("数据长度不对 %d-%d", datalen, newlen)
		}
	} else {
		glog.Errorf("数据长度不对 %d", datalen)
	}
}

//Cmdcodec 包解析
func (rc *RedisPacket) Cmdcodec() common.CmdCodec {
	return rc.cmdcodec
}

//SetFunc 收包函数
func (rc *RedisPacket) SetFunc(Queue func(pcmd *gpacket.BaseCmd, data []byte) bool) {
	rc.msgQueue = Queue
}

//NewRedis redis创建
func NewRedis(config *RedisConfig) (adapter *RedisPacket, err error) {
	adapter = new(RedisPacket)
	adapter.cmdcodec = new(common.BinaryCodec)
	adapter.RedisPool = newRedis()
	err = adapter.start(config)
	if err != nil {
		adapter = nil
	}
	return
}
