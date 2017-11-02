//读取配置文件

package logicgate

import (
	io "io/ioutil"

	"github.com/yhhaiua/goserver/common/gjson"
	"github.com/yhhaiua/goserver/common/glog"
)

type stGameConfig struct {
	sip      string
	sport    string
	serverid int32
}
type stJSONConfig struct {
	nloglvl     int    //日志等级
	sport       string //端口
	gameconfing []stGameConfig
}

func (Config *stJSONConfig) configInit(serverid int) bool {

	path := "./config/config.json"
	key := "gate"
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
		} else {
			glog.Errorf("Failed to config file '%s'", path)
			return false
		}
		if Config.nloglvl > 0 {
			glog.Setloglvl(Config.nloglvl)
		}
	}
	//读取game
	key = "game"
	keydata = gjson.NewGet(jsondata, key)
	if keydata.IsValid() {
		num := keydata.Getnum()
		var gameconfig stGameConfig
		for i := 0; i < num; i++ {
			data := gjson.NewGetindex(keydata, i)
			if data.IsValid() {
				gameconfig.sip = data.Getstring("ip")
				gameconfig.sport = data.Getstring("port")
				gameconfig.serverid = data.Getint32("id")
				Config.gameconfing = append(Config.gameconfing, gameconfig)
			} else {
				glog.Errorf("game Failed to config file '%s'", path)
				return false
			}
		}
	}
	glog.Infof("game连接数 %d", len(Config.gameconfing))
	return true
}
