package gtcp

import (
	"net"

	"github.com/yhhaiua/goserver/common"

	"github.com/yhhaiua/goserver/common/glog"
)

//ClientConnecter 请求连接结构
type ClientConnecter struct {
	*baseSession
	myTCPAddr *net.TCPAddr
	nServerID int32
	sendOnce  func()
}

//AddConnect 添加请求信息
func AddConnect(serverip, port string, serverid int32, servername string) *ClientConnecter {
	Connecter := new(ClientConnecter)
	connectadd := serverip + ":" + port
	Connecter.nServerID = serverid
	var err error
	Connecter.myTCPAddr, err = net.ResolveTCPAddr("tcp", connectadd)
	if err != nil {
		glog.Errorf("AddConnect error:%s", err)
		Connecter = nil
		return nil
	}
	Connecter.baseSession = addbase(nil, int64(serverid), servername)
	Connecter.newcodec(newcodecBinary)
	glog.Infof("尝试连接ip:[%s],prot:[%s],serverid:[%d]", serverip, port, serverid)
	return Connecter
}

//Start 开始连接
func (connect *ClientConnecter) Start() {
	//尝试第一次连接
	connect.startconnect()
	//压人请求连接队列
	mTCPConnMap.Put(connect)
}
func (connect *ClientConnecter) startconnect() bool {
	if !connect.boConnected {
		var err error
		connect.conn, err = net.DialTCP("tcp", nil, connect.myTCPAddr)
		if err != nil {
			glog.Errorf("startconnect连接失败4秒后再次连接 error:%s", err)
		} else {
			connect.start()
			connect.sendOnce()
			return true
		}
	}
	return false
}

//Cmdcodec 包解析
func (connect *ClientConnecter) Cmdcodec() common.CmdCodec {
	return connect.cmdcodec
}

//SetFunc 发送验证包的函数、读取数据包的函数
func (connect *ClientConnecter) SetFunc(Queue func(pcmd *common.BaseCmd, data []byte) bool, Once func()) {
	connect.msgQueue = Queue
	connect.sendOnce = Once
}

//SetValid 设置是否为激活状态
func (connect *ClientConnecter) SetValid(bodata bool) {
	connect.bovalid = bodata
}

//SendCmd 发送数据包
func (connect *ClientConnecter) SendCmd(data interface{}) {
	connect.sendCmd(data)
}
