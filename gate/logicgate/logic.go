package logicgate

import (
	"net"
	"sync"

	"github.com/yhhaiua/goserver/common"
	"github.com/yhhaiua/goserver/common/gtcp"
)

//SERVERTYPE 服务器类型
const SERVERTYPE = common.SERVERTYPEGATE
const (
	callbackPLAYER = 10000
)

//Logicsvr 服务器数据
type Logicsvr struct {
	mstJSONConfig stJSONConfig
	gameconmap    map[int32]*stGameCon
	serverid      int32
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

func (logic *Logicsvr) allconnect() bool {

	success := logic.gameConInit()

	return success
}
func (logic *Logicsvr) gameConInit() bool {
	logic.gameconmap = make(map[int32]*stGameCon)

	con := new(stGameCon)
	if con.create() {
		logic.gameconmap[logic.serverid] = con
		return true
	}
	return false
}
func (logic *Logicsvr) allListen() bool {

	success := gtcp.AddListen("0.0.0.0", logic.config().sport, callbackPLAYER, logic.ListenCallback)

	return success
}

//ListenCallback 监听回调
func (logic *Logicsvr) ListenCallback(con *net.TCPConn, backtype int32) {
	switch backtype {
	case callbackPLAYER:
	default:
	}
}
func (logic *Logicsvr) config() *stJSONConfig {
	return &logic.mstJSONConfig
}
