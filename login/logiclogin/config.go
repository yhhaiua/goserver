//读取配置文件

package logiclogin

import (
	io "io/ioutil"

	"github.com/yhhaiua/goserver/common/gjson"
	"github.com/yhhaiua/goserver/common/glog"
	"github.com/yhhaiua/goserver/common/gredis"
)

type stJSONConfig struct {
	nloglvl      int                //日志等级
	sport        string             //端口
	mredisconfig gredis.RedisConfig //redis连接信息
}

func (Config *stJSONConfig) configInit(serverid int) bool {

	path := "./config/config.json"
	key := "login"
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

		logindata := gjson.NewGetindex(keydata, i)

		if logindata.IsValid() {

			Config.nloglvl = logindata.Getint("loglvl")
			Config.sport = logindata.Getstring("port")
			redata := gjson.NewGet(logindata, "redis")
			if redata.IsValid() {
				Config.mredisconfig.Shostport = redata.Getstring("host")
				Config.mredisconfig.Maxopen = redata.Getint("open")
				Config.mredisconfig.Maxidle = redata.Getint("idle")
			} else {
				glog.Errorf("Failed to redis config file '%s'", path)
				return false
			}

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
