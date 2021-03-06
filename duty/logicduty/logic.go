package logicduty

import (
	"sync"
	"time"

	"github.com/yhhaiua/goserver/common/glog"
	"github.com/yhhaiua/goserver/common/gmysql"
	"github.com/yhhaiua/goserver/common/gredis"
	"github.com/yhhaiua/goserver/comsvrsrc"
)

//SERVERTYPE 服务器类型
const SERVERTYPE = comsvrsrc.SERVERTYPEDUTY

//Logicsvr 服务器数据
type Logicsvr struct {
	mstJSONConfig stJSONConfig
	redisConnect  *gredis.RedisPacket
	mysqlConnect  *gmysql.MysqlDB
	mysqlread     *stMysqlRead
	mysqlwrite    *stMysqlWrite
	myredismsg    *stRedisMsg
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
		//mysql连接
		logic.mysqlCon()
		//定时读取redis改变数据到mysql中
		logic.myslqWrite()
		//定时读取mysql数据到redis中
		logic.myslqRead()
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
	//设置redis接包
	logic.myredismsg = new(stRedisMsg)
	logic.redisConnect.SetFunc(logic.myredismsg.putMsgQueue)
	//监听订阅频道
	logic.redisConnect.Subscribe(comsvrsrc.SUBCHANNELduty)
}

func (logic *Logicsvr) mysqlCon() {

	var err error

	logic.mysqlConnect, err = gmysql.NewMysql(&logic.mstJSONConfig.mmysqlconfig)

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
func (logic *Logicsvr) myslqWrite() {
	logic.mysqlwrite = new(stMysqlWrite)
	logic.mysqlwrite.Write()
}
func (logic *Logicsvr) config() *stJSONConfig {
	return &logic.mstJSONConfig
}

func (logic *Logicsvr) mysqldb() *gmysql.MysqlDB {
	return logic.mysqlConnect
}

func (logic *Logicsvr) redisdb() *gredis.RedisPacket {
	return logic.redisConnect
}
func (logic *Logicsvr) redismsg() *stRedisMsg {
	return logic.myredismsg
}
