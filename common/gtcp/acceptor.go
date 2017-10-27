package gtcp

import (
	"net"

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
