//读取配置文件

package logicduty

import (
	io "io/ioutil"

	"github.com/yhhaiua/goserver/common"
	"github.com/yhhaiua/goserver/common/gjson"
	"github.com/yhhaiua/goserver/common/glog"
)

type stJSONConfig struct {
	nloglvl      int                //日志等级
	readdata     int                //是否读取数据(0读取1不读取)
	mredisconfig common.RedisConfig //redis连接信息
	mmysqlconfig common.MysqlConfig //mysql连接信息
}

func (Config *stJSONConfig) configInit(serverid int) bool {

	path := "./config/config.json"
	key := "duty"
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
			Config.readdata = logindata.Getint("readdata")
			redata := gjson.NewGet(logindata, "redis")
			if redata.IsValid() {
				Config.mredisconfig.Shostport = redata.Getstring("host")
				Config.mredisconfig.Maxopen = redata.Getint("open")
				Config.mredisconfig.Maxidle = redata.Getint("idle")
			} else {
				glog.Errorf("Failed to redis config file '%s'", path)
				return false
			}

			mysqldata := gjson.NewGet(logindata, "mysql")
			if mysqldata.IsValid() {

				Config.mmysqlconfig.Shost = mysqldata.Getstring("host")
				Config.mmysqlconfig.Sdbname = mysqldata.Getstring("dbname")
				Config.mmysqlconfig.Suser = mysqldata.Getstring("user")
				Config.mmysqlconfig.Spassword = mysqldata.Getstring("password")
				Config.mmysqlconfig.Maxopen = mysqldata.Getint("open")
				Config.mmysqlconfig.Maxidle = mysqldata.Getint("idle")
			} else {
				glog.Errorf("Failed to mysql config file '%s'", path)
				return false
			}
		} else {
			glog.Errorf("Failed to loginserver config file '%s'", path)
			return false
		}
		if Config.nloglvl > 0 {
			glog.Setloglvl(Config.nloglvl)
		}
	}

	return true
}
