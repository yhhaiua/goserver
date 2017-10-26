package logicgate

import (
	"sync"

	"github.com/yhhaiua/goserver/common"
)

//SERVERTYPE 服务器类型
const SERVERTYPE = common.SERVERTYPEGATE

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

		if logic.gameConInit() {

			return true
		}
	}
	return false
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
func (logic *Logicsvr) config() *stJSONConfig {
	return &logic.mstJSONConfig
}
