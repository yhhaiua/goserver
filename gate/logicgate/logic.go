package logicgate

import (
	"net"
	"sync"
	"sync/atomic"

	"github.com/yhhaiua/goserver/common/gtcp"
	"github.com/yhhaiua/goserver/comsvrsrc"
)

//SERVERTYPE 服务器类型
const SERVERTYPE = comsvrsrc.SERVERTYPEGATE
const (
	callbackPLAYER = 10000
)

//Logicsvr 服务器数据
type Logicsvr struct {
	mstJSONConfig stJSONConfig
	gameconmap    map[int32]*stGameCon
	playersessmap stPlayerSessMap
	serverid      int32
	linkKey       int64
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

		//连接与监听
		if logic.allconnect() && logic.allListen() {

			return true
		}
	}
	return false
}

//allconnect所有的连接
func (logic *Logicsvr) allconnect() bool {
	//连接gs队列
	logic.gameconmap = make(map[int32]*stGameCon)
	//连接gs
	success := logic.gameConInit()

	return success
}

//game连接
func (logic *Logicsvr) gameConInit() bool {
	con := new(stGameCon)
	if con.create() {
		logic.gameconmap[logic.serverid] = con
		return true
	}
	return false
}

//allListen所有的监听
func (logic *Logicsvr) allListen() bool {

	logic.playersessmap.sessMap = make(map[int64]*stPlayerSession)

	success := gtcp.AddListen("0.0.0.0", logic.config().sport, callbackPLAYER, logic.ListenCallback)

	return success
}

//ListenCallback 监听回调
func (logic *Logicsvr) ListenCallback(con *net.TCPConn, backtype int32) {
	switch backtype {
	case callbackPLAYER:
		logic.playerSessionInit(con)
	default:
	}
}

//玩家连接调用
func (logic *Logicsvr) playerSessionInit(con *net.TCPConn) {
	session := new(stPlayerSession)
	key := atomic.AddInt64(&logic.linkKey, 1)
	if session.create(con, key) {
		logic.playersessmap.add(session, key)
	}
}
func (logic *Logicsvr) config() *stJSONConfig {
	return &logic.mstJSONConfig
}

func (logic *Logicsvr) playerMap() *stPlayerSessMap {
	return &logic.playersessmap
}
