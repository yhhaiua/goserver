package conmanager

import "sync"
import "github.com/yhhaiua/goserver/protocol"
import "github.com/yhhaiua/goserver/manage/logicmanage"
import "github.com/yhhaiua/goserver/common/glog"

type stConnectInfo struct {
	nsvrtype int32
	sip      string
	sport    string
	keylink  int64
}

//ConManager 连接管理类
type ConManager struct {
	svrMap       map[int32]stConnectInfo
	manageConfig stManageConfig
}

var (
	instance *ConManager
	mu       sync.Mutex
)

//Instance 实例化ConManager
func Instance() *ConManager {
	if instance == nil {
		mu.Lock()
		defer mu.Unlock()
		if instance == nil {
			instance = new(ConManager)
		}
	}
	return instance
}

//Init 读取连接信息
func (manager *ConManager) Init() {
	manager.manageConfig.configInit()
	manager.svrMap = make(map[int32]stConnectInfo)
	go manager.thenpackage()
}

func (manager *ConManager) thenpackage() {

	for {
		tempbuf := <-logicmanage.Instance().Infolist
		manager.putMsgQueue(&tempbuf)
	}
}
func (manager *ConManager) putMsgQueue(pcmd *logicmanage.PackBaseInfo) {
	switch pcmd.Value() {
	case protocol.ServerCmdLoginCode:
		data := pcmd.Data.(*protocol.ServerCmdLogin)
		manager.connectadd(data.Svrid, data.Svrtype, data.Sip, data.Sport, pcmd.KeyLink)
	}
}
func (manager *ConManager) connectadd(nsvrid, nsvrtype int32, sip, sport string, keylink int64) {

	var info stConnectInfo
	info.nsvrtype = nsvrtype
	info.sip = sip
	info.sport = sport
	info.keylink = keylink
	manager.svrMap[nsvrid] = info
	glog.Infof("记录成功 id:[%d],type:[%d],ip:[%s],port:[%s],key:[%d]", nsvrid, nsvrtype, sip, sport, keylink)
	value, ok := manager.manageConfig.ContypeMap[nsvrtype]
	if ok {
		//是连接方，需要监听方的所有ip端口
		manager.sendAllCmd(value, keylink)
	} else {
		value, ok = manager.manageConfig.ServertypeMap[nsvrtype]
		if ok {
			//是监听方，把自己的ip端口发送给所有连接方
			manager.sendOneCmd(nsvrid, value)
		}
	}
}

func (manager *ConManager) sendOneCmd(svrid, servertype int32) {
	var retcmd protocol.ServerCmdConData
	retcmd.Init()
	ref := manager.getConData(svrid)
	retcmd.ConDataInfo = append(retcmd.ConDataInfo, ref)

	for _, value := range manager.svrMap {
		if value.nsvrtype == servertype {
			logicmanage.Instance().SendGateCmd(value.keylink, &retcmd)
		}
	}
}
func (manager *ConManager) getConData(svrid int32) (ref protocol.RefConDataInfo) {
	value, ok := manager.svrMap[svrid]
	if ok {
		ref.Sip = value.sip
		ref.Sport = value.sport
		ref.Svrid = svrid
		ref.Svrtype = value.nsvrtype
	}
	return
}

func (manager *ConManager) sendAllCmd(servertype int32, keylink int64) {

	var retcmd protocol.ServerCmdConData
	retcmd.Init()
	var ref protocol.RefConDataInfo
	for key, value := range manager.svrMap {
		if value.nsvrtype == servertype {
			ref.Sip = value.sip
			ref.Sport = value.sport
			ref.Svrid = key
			ref.Svrtype = value.nsvrtype
			retcmd.ConDataInfo = append(retcmd.ConDataInfo, ref)
		}
	}
	logicmanage.Instance().SendGateCmd(keylink, &retcmd)
}
