package logicduty

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/yhhaiua/goserver/common/glog"
)

type stRedisVersion struct {
	Pid     int
	Version string
	URL     string
}

type stVersionData struct {
	pname     string
	redisdata stRedisVersion
}

type stZoneData struct {
	Pname     string
	Pid       int
	Zoneid    int
	Zonename  string
	Loginip   string
	Loginport string
}

type stMysqlRead struct {
	versionMap  map[string]stVersionData //平台配置
	zonedataMap map[string]stZoneData    //区服配置
}

const (
	datetickertime = 20 * time.Second
	versionData    = "VersionData_" //平台配置表
	gameData       = "game_"        //区服表

	cAccountData = "AccountData"  //帐号数据
	globalonly   = "GlobalOnlyid" //全局唯一id数据
)

func (mydata *stMysqlRead) Read() {

	//平台配置
	mydata.versionMap = make(map[string]stVersionData)
	//区服配置
	mydata.zonedataMap = make(map[string]stZoneData)

	if Instance().mysqldb().HaveConnect() {

		mydata.dataRead()

		if Instance().config().isRead() {
			mydata.onceRead()
		}

		//for {
		//	Instance().redisdb().Publish(comsvrsrc.SUBCHANNELlogin, "成功")
		//	time.Sleep(time.Second * 5)
		//}

		for _ = range time.NewTicker(datetickertime).C {

			mydata.dataRead()
		}
	}

}

func (mydata *stMysqlRead) dataRead() {

	if Instance().redisdb() != nil {
		//平台配置读取
		mydata.versionDataRead()
		//区服配置读取
		mydata.zoneDataRead()
	}

}

func (mydata *stMysqlRead) onceRead() {

	if Instance().redisdb() != nil {
		//玩家id
		mydata.onlyidDataRead()
		//玩家帐号数据
		mydata.accountDataRead()

		glog.Info("onceRead 读取数据完成")
	}

}
func (mydata *stMysqlRead) versionDataRead() {

	rows, err := Instance().mysqldb().Query("SELECT * FROM VersionData")

	if err == nil {
		defer rows.Close()
		var Data stVersionData

		tempversionMap := make(map[string]stVersionData)

		for rows.Next() {

			err = rows.Scan(&Data.redisdata.Pid, &Data.pname, &Data.redisdata.Version, &Data.redisdata.URL)
			if err != nil {
				glog.Errorf("mysql读取错误 VersionData 2 %s", err)
			} else {

				tempversionMap[Data.pname] = Data

				var boaddredis bool
				captial, ok := mydata.versionMap[Data.pname]
				if ok {
					if captial != Data {
						boaddredis = true
					}
				} else {
					boaddredis = true
				}

				if boaddredis {
					mydata.versionMap[Data.pname] = Data

					//保存到redis
					TempDataName := versionData + Data.pname

					RedisMessage, err := json.Marshal(Data.redisdata)
					if err == nil {
						Instance().redisdb().Set(TempDataName, RedisMessage)
						glog.Infof("mysql转存到redis成 %s, %s", TempDataName, RedisMessage)
					} else {
						glog.Errorf("mysql转存到redis失败 %s, %s", TempDataName, err)
					}
				}
			}
		}

		//没有的删除
		for mapkey := range mydata.versionMap {
			_, ok := tempversionMap[mapkey]
			if !ok {
				TempDataName := versionData + mapkey
				Instance().redisdb().Delete(TempDataName)
				glog.Infof("reidis成功删除 %s", TempDataName)
				delete(mydata.versionMap, mapkey)
			}
		}
	} else {

		glog.Errorf("mysql读取错误 VersionData 1 %s", err)
	}
}

func (mydata *stMysqlRead) zoneDataRead() {

	rows, err := Instance().mysqldb().Query("SELECT pname,pid,zoneid,zonename,loginip,loginport FROM game")
	if err == nil {
		defer rows.Close()
		var Data stZoneData

		var tempMap map[string]stZoneData

		tempMap = make(map[string]stZoneData)

		for rows.Next() {

			err = rows.Scan(&Data.Pname, &Data.Pid, &Data.Zoneid, &Data.Zonename, &Data.Loginip, &Data.Loginport)
			if err != nil {
				glog.Errorf("mysql读取错误 game 2 %s", err)
			} else {

				if !mydata.panmeHave(Data.Pname, Data.Pid) {
					glog.Errorf("表 game 中有错误平台 Zoneid:[%d],Panme:[%s],Pid:[%d]", Data.Zoneid, Data.Pname, Data.Pid)
					continue
				}
				keydata := Data.Pname + strconv.Itoa(Data.Zoneid)

				tempMap[keydata] = Data

				var boaddredis bool
				captial, ok := mydata.zonedataMap[keydata]
				if ok {
					if captial != Data {
						boaddredis = true
					}
				} else {
					boaddredis = true
				}

				if boaddredis {
					mydata.zonedataMap[keydata] = Data
					//创建表
					Instance().mysqldb().Create(cAccountData + keydata)
					//保存到redis
					TempDataName := gameData + keydata

					RedisMessage, err := json.Marshal(Data)
					if err == nil {
						Instance().redisdb().Set(TempDataName, RedisMessage)
						glog.Infof("mysql转存到redis成 %s, %s", TempDataName, RedisMessage)
					} else {
						glog.Errorf("mysql转存到redis失败 %s, %s", TempDataName, err)
					}
				}
			}
		}

		//没有的删除
		for mapkey := range mydata.zonedataMap {
			_, ok := tempMap[mapkey]
			if !ok {
				TempDataName := gameData + mapkey
				Instance().redisdb().Delete(TempDataName)
				glog.Infof("reidis成功删除 %s", TempDataName)
				delete(mydata.zonedataMap, mapkey)
			}
		}
	} else {

		glog.Errorf("mysql读取错误 game 1 %s", err)
	}
}

func (mydata *stMysqlRead) panmeHave(key string, pid int) bool {
	value, ok := mydata.versionMap[key]
	if ok {
		if value.redisdata.Pid != pid {
			return false
		}
	}
	return ok
}

func (mydata *stMysqlRead) onlyidDataRead() {

	rows, err := Instance().mysqldb().Query("SELECT * FROM GlobalOnlyid")
	if err == nil {
		defer rows.Close()
		var i64onlyid int64
		var Key, Value string
		for rows.Next() {

			err = rows.Scan(&Key, &Value)
			if err != nil {
				glog.Errorf("mysql读取错误 GlobalOnlyid 2 %s", err)
			} else {
				i64onlyid, _ = strconv.ParseInt(Value, 10, 64)
				if i64onlyid > 0 {
					//保存到redis
					skey := globalonly + "_" + Key
					Instance().redisdb().Set(skey, i64onlyid)
					glog.Infof("redis添加玩家id %s, %d", skey, i64onlyid)
				}
			}
		}

	} else {

		glog.Errorf("mysql读取错误 GlobalOnlyid 1 %s", err)
	}
}

func (mydata *stMysqlRead) accountDataRead() {

	for mapkey := range mydata.zonedataMap {

		tablename := cAccountData + mapkey

		Instance().mysqldb().SavetoRedis(Instance().redisdb(), tablename)
	}

}
