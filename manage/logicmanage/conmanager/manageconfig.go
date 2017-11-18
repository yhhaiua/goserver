package conmanager

import (
	io "io/ioutil"

	"github.com/yhhaiua/goserver/common/gjson"
	"github.com/yhhaiua/goserver/common/glog"
)

type stManageConfig struct {
	ContypeMap    map[int32]int32 //conmap
	ServertypeMap map[int32]int32 //servermap
}

func (Config *stManageConfig) configInit() bool {

	path := "./config/manageconfig.json"
	key := "list"
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

		num := keydata.Getnum()
		Config.ContypeMap = make(map[int32]int32)
		Config.ServertypeMap = make(map[int32]int32)
		for i := 0; i < num; i++ {
			data := gjson.NewGetindex(keydata, i)
			if data.IsValid() {
				contype := data.Getint32("connect")
				servertype := data.Getint32("server")
				Config.ContypeMap[contype] = servertype
				Config.ServertypeMap[servertype] = contype
				glog.Infof("内部连接配置 %d-%d", contype, servertype)
			} else {
				glog.Errorf("Failed to config file '%s'", path)
				return false
			}
		}
	}

	return true
}
