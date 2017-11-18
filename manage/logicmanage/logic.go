package logicmanage

import (
	"net"
	"sync"
	"sync/atomic"

	"github.com/yhhaiua/goserver/common/gpacket"

	"github.com/yhhaiua/goserver/common/gtcp"
	"github.com/yhhaiua/goserver/comsvrsrc"
)

//SERVERTYPE 服务器类型
const SERVERTYPE = comsvrsrc.SERVERTYPEMANAGE
const (
	callbackGate = 10000
	callbackGame = 10001
)

//PackBaseInfo 内部通讯包
type PackBaseInfo struct {
	KeyLink int64
	*gpacket.BaseCmd
	Data interface{}
}

//Logicsvr 服务器数据
type Logicsvr struct {
	mstJSONConfig stJSONConfig
	gateMap       *sync.Map
	serverid      int32
	linkKey       int64
	Infolist      chan PackBaseInfo
}

var (
	instance *Logicsvr
	mu       sync.Mutex
)

//Instance 实例化logicsvr
func Instance() *Logicsvr {
	if instance == nil {
		mu.Lock()
		defer mu.Unlock()
		if instance == nil {
			instance = new(Logicsvr)
		}
	}
	return instance
}

//LogicInit 初始化
func (logic *Logicsvr) LogicInit(serverid int) bool {

	//读取配置
	logic.serverid = int32(serverid)
	if logic.mstJSONConfig.configInit(serverid) {
		logic.Infolist = make(chan PackBaseInfo, 10)
		//连接与监听
		if logic.allconnect() && logic.allListen() {

			return true
		}
	}
	return false
}

//allconnect所有的连接
func (logic *Logicsvr) allconnect() bool {

	success := true
	return success
}

//allListen所有的监听
func (logic *Logicsvr) allListen() bool {

	logic.gateMap = new(sync.Map)

	successgate := gtcp.AddListen("0.0.0.0", logic.config().sport, callbackGate, logic.ListenCallback)

	return successgate
}

func (logic *Logicsvr) getLinkKey() int64 {
	return atomic.AddInt64(&logic.linkKey, 1)
}

//ListenCallback 监听回调
func (logic *Logicsvr) ListenCallback(con *net.TCPConn, backtype int32) {
	switch backtype {
	case callbackGate:
		logic.gateSessionInit(con)
	default:
	}
}

//玩家连接调用
func (logic *Logicsvr) gateSessionInit(con *net.TCPConn) {
	session := new(stGateSession)

	key := logic.getLinkKey()

	if session.create(con, key) {
		logic.gateMap.Store(key, session)
	}
}

func (logic *Logicsvr) config() *stJSONConfig {
	return &logic.mstJSONConfig
}

func (logic *Logicsvr) syncgateMap() *sync.Map {
	return logic.gateMap
}

//SendGateCmd 发送给网关信息
func (logic *Logicsvr) SendGateCmd(key int64, data interface{}) {

	value, ok := logic.gateMap.Load(key)
	if ok {

		session, zok := value.(*stGateSession)
		if zok {
			session.SendCmd(data)
		}
	}

}
func (logic *Logicsvr) putMsgList(keyLink int64, pcmd *gpacket.BaseCmd, data interface{}) {
	var info PackBaseInfo
	info.KeyLink = keyLink
	info.BaseCmd = pcmd
	info.Data = data
	logic.Infolist <- info
}
