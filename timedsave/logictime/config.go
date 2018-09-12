package logictime

import ()

type stJSONConfig struct {
	source string //redis 目标地址
	target string //保存目录
}

func (Config *stJSONConfig) configInit() bool {
	// path := "./config/config.json"
	// data, err := ioutil.ReadFile(path)
	// if err != nil {
	// 	glog.Errorf("Failed to open config file '%s': %s\n", path, err)
	// 	return false
	// }
	// err = json.Unmarshal(data, Config)
	// if err != nil {
	// 	glog.Errorf("json error '%s': %s\n", path, err)
	//     return false
	// }
	return true
}
