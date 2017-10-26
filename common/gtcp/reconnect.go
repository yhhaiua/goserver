package gtcp

import (
	"net"
	"sync"
	"time"

	"github.com/yhhaiua/goserver/common/glog"
)

var mTCPConnMap TCPConnMap

//TCPConnMap 所有请求连接队列
type TCPConnMap struct {
	sync.Mutex
	mymap map[int32]*ClientConnecter
}

//Put 向队列中压人新的请求
func (m *TCPConnMap) Put(clent *ClientConnecter) {

	m.Lock()
	defer m.Unlock()
	m.mymap[clent.nServerID] = clent
}

//Run 循环队列请求连接
func (m *TCPConnMap) Run() {

	for {
		m.TimeAction()
		time.Sleep(4 * time.Second)
	}

}

//TimeAction 激活连接
func (m *TCPConnMap) TimeAction() {
	m.Lock()
	defer m.Unlock()
	for _, value := range m.mymap {
		if !value.boConnected {
			var err error
			value.conn, err = net.DialTCP("tcp", nil, value.myTCPAddr)
			if err != nil {
				glog.Errorf("timeAction重新连接失败4秒后再次连接 error:%s", err)
				continue
			} else {
				value.doInit()
				value.SendOnce()
				go value.runRead()
				go value.runWrite()
			}
		}
	}
}
func (m *TCPConnMap) newTCPConnMap() {
	m.mymap = make(map[int32]*ClientConnecter)
}

func init() {
	mTCPConnMap.newTCPConnMap()
	go mTCPConnMap.Run()
}
