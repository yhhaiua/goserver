package logic

import (
	"github.com/yhhaiua/goserver/common/glog"
	"github.com/yhhaiua/goserver/common/gredis"
	"github.com/yhhaiua/goserver/common/grouter"
	"net/http"
	"sync"
	"time"
)

var (
	instance *LogicSvr
	mu       sync.Mutex
)
//LogicSvr 服务器数据
type LogicSvr struct {
	mstJSONConfig stJSONConfig
	routerConnect stRouterPost
	redisConnect  *gredis.RedisPacket
	playerBuyMap     *sync.Map
	backConnect   stBackstage
	stopRecharge	bool
}


//Instance 实例化LogicSvr
func Instance() *LogicSvr {
	if instance == nil {
		mu.Lock()
		defer mu.Unlock()
		if instance == nil {
			instance = new(LogicSvr)
		}
	}
	return instance
}

//LogicInit 初始化
func (logic *LogicSvr) LogicInit() bool {
	if logic.mstJSONConfig.configInit(){
		logic.redisCon()
		logic.playerBuyMap = new(sync.Map)
		logic.stopRecharge = false
		return  logic.routerInit()
	}
	return false
}

func (logic *LogicSvr) redisCon() {
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
}

func (logic *LogicSvr) routerInit() bool{

	router := grouter.New()

	//Get service information
	router.GET("/recharge", logic.routerConnect.rechargeDeal)
	router.GET("/stopcharge", logic.backConnect.stopCharge)
	router.GET("/makeuporder", logic.backConnect.makeUpOrder)

	glog.Infof("http监听开启%s", logic.mstJSONConfig.sport)
	glog.Infoln("当前版本:v1.0.1")

	srv := &http.Server{
		ReadTimeout: 30 * time.Second,
		WriteTimeout: 30 * time.Second,
		Addr:logic.mstJSONConfig.sport,
		Handler : router,
	}

	err := srv.ListenAndServe()
	if err != nil {
		glog.Errorf("http监听失败 %s", err)
		return false
	}
	return true
}

//获取物品价格
func (logic *LogicSvr) getMoney(itemid int) (int,bool)  {
	value,ok := logic.mstJSONConfig.chargeconfig[itemid]
	return value,ok
}

func (logic *LogicSvr) redisdb() *gredis.RedisPacket {
	return logic.redisConnect
}

func (logic *LogicSvr)addBuyMap(playerid string)  {
	logic.playerBuyMap.Store(playerid,playerid)
}
func (logic *LogicSvr) checkBuyMapKey(playerid string) bool  {
	_,ok:=logic.playerBuyMap.Load(playerid)
	return ok
}
func (logic *LogicSvr)delBuyMap(playerid string)  {
	logic.playerBuyMap.Delete(playerid)
}

//获取文件ip
func (logic *LogicSvr)getUserIp(r *http.Request) string  {
	userip := r.Header.Get("X-Real-IP")
	if userip == ""{
		userip = r.RemoteAddr
	}
	return userip
}