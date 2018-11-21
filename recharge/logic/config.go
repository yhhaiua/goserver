package logic

import (
	"encoding/xml"
	"github.com/yhhaiua/goserver/common/gjson"
	"github.com/yhhaiua/goserver/common/gredis"
	"github.com/yhhaiua/goserver/common/log4go"
	io "io/ioutil"
	"strconv"
)

type stJSONConfig struct {
	nloglvl      int                //日志等级
	sport        string             //http端口
	gmhost		 string				//发送的后台地址
	rechargekey  string
	clientkey 	 string
	mredisconfig gredis.RedisConfig //redis连接信息
	operatorid 	 string				//对应平台
	chargeconfig map[int]int
}

type UnitOne struct {
	Id   string `xml:"id,attr"`
	Money string `xml:"money,attr"`
}

type UnitConfig struct {
	Unit []UnitOne `xml:"unit"`
}

func (Config *stJSONConfig) configInit() bool {

	path := "./config/config.json"
	key := "recharge"
	data, err := io.ReadFile(path)
	if err != nil {
		log4go.Error("Failed to open config file '%s': %s\n", path, err)
		return false
	}

	jsondata, err := gjson.NewJSONByte(data)
	if err != nil {
		log4go.Error("Failed to NewJsonByte config file '%s': %s\n", path, err)
		return false
	}

	keydata := gjson.NewGet(jsondata, key)

	if keydata.IsValid() {

		i := 0

		logindata := gjson.NewGetindex(keydata, i)

		if logindata.IsValid() {

			Config.nloglvl = logindata.Getint("loglvl")
			Config.sport = logindata.Getstring("port")
			Config.gmhost = logindata.Getstring("gmhost")
			Config.rechargekey = logindata.Getstring("rechargekey")
			Config.clientkey = logindata.Getstring("clientkey")
			Config.operatorid = logindata.Getstring("operatorid")
			redata := gjson.NewGet(logindata, "redis")

			if redata.IsValid() {
				Config.mredisconfig.Shostport = redata.Getstring("host")
				Config.mredisconfig.Maxopen = redata.Getint("open")
				Config.mredisconfig.Maxidle = redata.Getint("idle")
				Config.mredisconfig.Password = redata.Getstring("password")
			} else {
				log4go.Error("Failed to redis config file '%s'", path)
				return false
			}

		} else {
			log4go.Error("Failed to config file '%s'", path)
			return false
		}
	}

	return Config.configXml()
}

func (Config *stJSONConfig) configXml() bool {
	path := "./config/charge.xml"
	content, err := io.ReadFile(path)
	if err != nil {
		log4go.Error("Failed to open config file '%s': %s\n", path, err)
		return false
	}
	var tempConfig	 UnitConfig
	err = xml.Unmarshal(content, &tempConfig)
	if err != nil {
		log4go.Error("Failed to Unmarshal config file '%s': %s\n", path, err)
		return false
	}
	Config.chargeconfig = make(map[int]int)
	for _,temp := range tempConfig.Unit {
		id,_:=strconv.Atoi(temp.Id)
		money,_:= strconv.Atoi(temp.Money)
		Config.chargeconfig[id] = money
	}
	return true
}