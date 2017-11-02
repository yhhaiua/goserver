package logicgame

import (
	io "io/ioutil"

	"github.com/yhhaiua/goserver/common/gjson"
	"github.com/yhhaiua/goserver/common/glog"
)

type stJSONConfig struct {
	nloglvl int    //日志等级
	sip     string //ip
	sport   string //端口
}

func (Config *stJSONConfig) configInit(serverid int) bool {

	path := "./config/config.json"
	key := "game"
	data, err := io.ReadFile(path)
	if err != nil {
		glog.Errorf("Failed to open config file '%s': %s\n", path, err)
		return false
	}

	jsondata, err := gjson.NewJSONByte(data)
	if err != nil {
		glog.Errorf("Failed to NewJsonByte config file '%s': %s\n", path, err)
		return false
	}

	keydata := gjson.NewGet(jsondata, key)

	if keydata.IsValid() {

		i := serverid - SERVERTYPE

		data := gjson.NewGetindex(keydata, i)

		if data.IsValid() {

			Config.nloglvl = data.Getint("loglvl")
			Config.sport = data.Getstring("port")
			Config.sip = data.Getstring("ip")
		} else {
			glog.Errorf("Failed to config file '%s'", path)
			return false
		}
		if Config.nloglvl > 0 {
			glog.Setloglvl(Config.nloglvl)
		}
	}

	return true
}
