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
	playerMap     *sync.Map
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
	num := len(logic.mstJSONConfig.gameconfing)
	for i := 0; i < num; i++ {
		con := new(stGameCon)
		game := &(logic.mstJSONConfig.gameconfing[i])
		if con.create(game) {
			logic.gameconmap[game.serverid] = con
		} else {
			return false
		}
	}
	return true
}

//allListen所有的监听
func (logic *Logicsvr) allListen() bool {

	logic.playerMap = new(sync.Map)

	success := gtcp.AddListen("0.0.0.0", logic.config().sport, callbackPLAYER, logic.ListenCallback)

	return success
}

func (logic *Logicsvr) getLinkKey() int64 {
	return atomic.AddInt64(&logic.linkKey, 1)
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
	key := logic.getLinkKey()
	if session.create(con, key) {
		logic.playerMap.Store(key, session)
	}
}
func (logic *Logicsvr) config() *stJSONConfig {
	return &logic.mstJSONConfig
}

func (logic *Logicsvr) syncMap() *sync.Map {
	return logic.playerMap
}

//SendPlayerCmd 发送给玩家信息
func (logic *Logicsvr) SendPlayerCmd(key int64, data interface{}) {

	value, ok := logic.playerMap.Load(key)
	if ok {

		session, zok := value.(*stPlayerSession)
		if zok {
			session.SendCmd(data)
		}
	}

}
