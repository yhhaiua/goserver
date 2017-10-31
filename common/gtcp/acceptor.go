package gtcp

import (
	"net"
	"sync"
	"time"

	"github.com/yhhaiua/goserver/common/glog"
)

//CallbackCon 回调连接函数
type CallbackCon func(con *net.TCPConn, backtype int32)

//AddListen 监听信息
func AddListen(serverip, port string, backtype int32, callbackcon CallbackCon) bool {

	connectadd := serverip + ":" + port

	myTCPAddr, err := net.ResolveTCPAddr("tcp", connectadd)
	if err != nil {
		glog.Errorf("AddListen error:%s", err)
		return false
	}
	lister, err := net.ListenTCP("tcp", myTCPAddr)

	if err != nil {
		glog.Errorf("AddListen ListenTCP error:%s", err)
		return false
	}
	glog.Infof("服务器监听开始ip:[%s],prot:[%s],backtype:[%d]", serverip, port, backtype)
	go runAccept(lister, backtype, callbackcon)

	return true
}

func runAccept(lister *net.TCPListener, backtype int32, callbackcon CallbackCon) {

	for {
		con, err := lister.AcceptTCP()
		if err != nil {
			glog.Errorf("runAccept error:%s", err)
			return
		}
		callbackcon(con, backtype)
	}

}

///////////////////////////////////////////////////////////////////////////////////////////////
var mCheckSessionMap stCheckSession

type stCheckSession struct {
	mymap *sync.Map
}

//Put 向队列中压人新的请求
func (m *stCheckSession) Put(clent *ServerSession) {
	m.mymap.Store(clent.servertag, clent)
}
func (m *stCheckSession) Del(servertag int64) {
	m.mymap.Delete(servertag)
}

//Run 循环队列请求连接
func (m *stCheckSession) Run() {

	for {
		m.mymap.Range(m.runCheck)
		time.Sleep(5 * time.Second)
	}

}

func (m *stCheckSession) runCheck(key, value interface{}) bool {
	connect, zok := value.(*ServerSession)
	if zok {
		if connect.baseSession != nil {
			connect.runCheck()
		}
		return true
	}
	return false
}
func (m *stCheckSession) newTCPConnMap() {
	m.mymap = new(sync.Map)
}

func init() {
	mCheckSessionMap.newTCPConnMap()
	go mCheckSessionMap.Run()
}
