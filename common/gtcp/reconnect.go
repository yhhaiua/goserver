package gtcp

import (
	"sync"
	"time"
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
		value.startconnect()
	}
}
func (m *TCPConnMap) newTCPConnMap() {
	m.mymap = make(map[int32]*ClientConnecter)
}

func init() {
	mTCPConnMap.newTCPConnMap()
	go mTCPConnMap.Run()
}
