package logicduty

import "sync"

//SERVERTYPE 服务器类型
const SERVERTYPE = 1100

//Logicsvr 服务器数据
type Logicsvr struct {
	mstJSONConfig stJSONConfig
}

var instance *Logicsvr

var mu sync.Mutex

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
	if logic.mstJSONConfig.configInit(serverid) {

	}
	return false
}
