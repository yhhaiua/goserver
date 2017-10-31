package gtcp

import (
	"sync"
	"time"
)

var mTCPConnMap TCPConnMap

//TCPConnMap 所有请求连接队列
type TCPConnMap struct {
	mymap *sync.Map
}

//Put 向队列中压人新的请求
func (m *TCPConnMap) Put(clent *ClientConnecter) {
	m.mymap.Store(clent.nServerID, clent)
}

//Run 循环队列请求连接
func (m *TCPConnMap) Run() {

	for {
		m.TimeAction()
		time.Sleep(5 * time.Second)
	}

}

//TimeAction 激活连接
func (m *TCPConnMap) TimeAction() {
	m.mymap.Range(m.runCheck)
}

func (m *TCPConnMap) runCheck(key, value interface{}) bool {
	connect, zok := value.(*ClientConnecter)
	if zok {
		connect.startconnect()

		if connect.baseSession != nil {
			connect.runCheck()
		}
		return true
	}
	return false
}
func (m *TCPConnMap) newTCPConnMap() {
	m.mymap = new(sync.Map)
}

func init() {
	mTCPConnMap.newTCPConnMap()
	go mTCPConnMap.Run()
}
