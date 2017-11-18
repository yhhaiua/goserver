package gtcp

import (
	"net"

	"github.com/yhhaiua/goserver/common/ginter"

	"github.com/yhhaiua/goserver/common/glog"
)

//ClientConnecter 请求连接结构
type ClientConnecter struct {
	*baseSession
	myTCPAddr  *net.TCPAddr
	nServerID  int32
	clientname string
	agent      ginter.ConnectAgenter
}

//AddConnect 添加请求信息
func AddConnect(serverip, port string, serverid int32, servername string, agent ginter.ConnectAgenter) *ClientConnecter {
	Connecter := new(ClientConnecter)
	connectadd := serverip + ":" + port
	Connecter.nServerID = serverid
	Connecter.clientname = servername
	Connecter.agent = agent
	var err error
	Connecter.myTCPAddr, err = net.ResolveTCPAddr("tcp", connectadd)
	if err != nil {
		glog.Errorf("AddConnect error:%s", err)
		Connecter = nil
		return nil
	}

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
	if connect.baseSession == nil || connect.boRecon {

		conn, err := net.DialTCP("tcp", nil, connect.myTCPAddr)
		if err != nil {
			glog.Errorf("startconnect连接失败4秒后再次连接 error:%s", err)
		} else {
			connect.baseInit(conn)
			connect.start()
			connect.agent.SendOnceCmd()
			return true
		}
	}
	return false
}

func (connect *ClientConnecter) baseInit(conn *net.TCPConn) {
	//baseSession结构中的参数初始化
	connect.baseSession = addbase(conn, int64(connect.nServerID), connect.clientname, connectbaseType, connect.agent)
}

//SetValid 设置是否为激活状态
func (connect *ClientConnecter) SetValid(bodata bool) {
	connect.bovalid = bodata
}

//SendCmd 发送数据包
func (connect *ClientConnecter) SendCmd(data interface{}) {
	connect.sendCmd(data)
}

//Close 关闭连接
func (connect *ClientConnecter) Close() {
	glog.Infof("主动调用关闭 %s,%d", connect.sname, connect.servertag)
	connect.close()
}
