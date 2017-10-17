package logicduty

import (
	"sync"
	"time"

	"github.com/yhhaiua/goserver/common"
	"github.com/yhhaiua/goserver/common/glog"
)

//SERVERTYPE 服务器类型
const SERVERTYPE = 1100

//Logicsvr 服务器数据
type Logicsvr struct {
	mstJSONConfig stJSONConfig
	redisConnect  *common.RedisPool
	mysqlConnect  *common.MysqlDB
	mysqlread     *stMysqlRead
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
		//reidis连接
		logic.redisCon()
		//mysql连接
		logic.mysqlCon()
		//定时读取mysql数据到redis中
		logic.myslqRead()
		return true
	}
	return false
}

func (logic *Logicsvr) redisCon() {
	for {
		var err error

		logic.redisConnect, err = common.NewRedis(&logic.mstJSONConfig.mredisconfig)

		if err == nil {
			glog.Infof("redis连接成功%s", logic.mstJSONConfig.mredisconfig.Shostport)
			break
		} else {
			glog.Errorf("redis连接错误%s，等待5秒后再次连接 %s", logic.mstJSONConfig.mredisconfig.Shostport, err)
			time.Sleep(5 * time.Second)
		}
	}
}

func (logic *Logicsvr) mysqlCon() {

	var err error

	logic.mysqlConnect, err = common.NewMysql(&logic.mstJSONConfig.mmysqlconfig)

	if err == nil {
		for {
			err = logic.mysqlConnect.CheckPing()

			if err == nil {
				glog.Infof("mysql连接成功%s", logic.mstJSONConfig.mmysqlconfig.Shost)
				break
			} else {
				glog.Errorf("mysql连接错误%s，等待5秒后再次连接 %s", logic.mstJSONConfig.mmysqlconfig.Shost, err)
				time.Sleep(5 * time.Second)
			}
		}
	} else {
		glog.Errorf("mysql配置错误")
	}
}

func (logic *Logicsvr) myslqRead() {
	logic.mysqlread = new(stMysqlRead)
	logic.mysqlread.Read()
}
func (logic *Logicsvr) config() *stJSONConfig {
	return &logic.mstJSONConfig
}

func (logic *Logicsvr) mysqldb() *common.MysqlDB {
	return logic.mysqlConnect
}
