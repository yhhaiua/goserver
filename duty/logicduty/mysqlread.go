package logicduty

import (
	"time"
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

const datetickertime = 20 * time.Second

func (mydata *stMysqlRead) Read() {

	//平台配置
	mydata.versionMap = make(map[string]stVersionData)
	//区服配置
	mydata.zonedataMap = make(map[string]stZoneData)

	if Instance().mysqldb().HaveConnect() {

		mydata.dataRead()

		if Instance().config().isRead() {
		}

		for _ = range time.NewTicker(datetickertime).C {

			mydata.dataRead()
		}
	}

}

func (mydata *stMysqlRead) dataRead() {

}
