package logiclogin

import (
	"net/http"
	"sync"
	"time"

	"github.com/yhhaiua/goserver/common/glog"
	"github.com/yhhaiua/goserver/common/gredis"
	"github.com/yhhaiua/goserver/common/grouter"
	"github.com/yhhaiua/goserver/comsvrsrc"
)

//SERVERTYPE 服务器类型
const SERVERTYPE = comsvrsrc.SERVERTYPELOGIN

//Logicsvr 服务器数据
type Logicsvr struct {
	mstJSONConfig stJSONConfig
	redisConnect  *gredis.RedisPacket
	routerConnect stRouterPost
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
	if logic.mstJSONConfig.configInit(serverid) {
		//reidis连接
		logic.redisCon()
		//http监听
		logic.routerInit()
		return true
	}
	return false
}

func (logic *Logicsvr) redisCon() {
	for {
		var err error

		logic.redisConnect, err = gredis.NewRedis(&logic.mstJSONConfig.mredisconfig)

		if err == nil {
			glog.Infof("redis连接成功%s", logic.mstJSONConfig.mredisconfig.Shostport)
			break
		} else {
			glog.Errorf("redis连接错误%s，等待5秒后再次连接 %s", logic.mstJSONConfig.mredisconfig.Shostport, err)
			time.Sleep(5 * time.Second)
		}
	}
	//监听订阅频道
	logic.redisConnect.Subscribe(comsvrsrc.SUBCHANNELlogin)
}

func (logic *Logicsvr) routerInit() {

	router := grouter.New()

	//Get service information
	router.POST("/public/getPlatformZone", logic.routerConnect.getPlatformZone)
	//Get account information
	router.POST("/public/regAccount", logic.routerConnect.regAccount)

	glog.Infof("http监听开启%s", logic.mstJSONConfig.sport)
	err := http.ListenAndServe(logic.mstJSONConfig.sport, router)
	if err != nil {
		glog.Errorf("http监听s失败 %s", err)
	}

}

func (logic *Logicsvr) config() *stJSONConfig {
	return &logic.mstJSONConfig
}

func (logic *Logicsvr) redisdb() *gredis.RedisPacket {
	return logic.redisConnect
}
