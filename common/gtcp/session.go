package gtcp

import (
	"net"

	"github.com/yhhaiua/goserver/common"
	"github.com/yhhaiua/goserver/common/glog"
	"github.com/yhhaiua/goserver/common/gpacket"
)

//ServerSession 请求连接结构
type ServerSession struct {
	*baseSession
	myDel func(servertag int64)
}

//AddSession 添加请求信息
func AddSession(conn *net.TCPConn, backtype int64, sname string) *ServerSession {
	Session := new(ServerSession)
	Session.baseSession = addbase(conn, backtype, sname)
	//包解析
	Session.newcodec(newcodecBinary)

	glog.Infof("有新的连接进入 %s,%d", sname, backtype)
	return Session
}

//Cmdcodec 包解析
func (connect *ServerSession) Cmdcodec() common.CmdCodec {
	return connect.cmdcodec
}

//SetFunc 发送验证包的函数、断开回调包
func (connect *ServerSession) SetFunc(Queue func(pcmd *gpacket.BaseCmd, data []byte) bool, Del func(servertag int64)) {
	connect.msgQueue = Queue
	connect.myDel = Del
	connect.delLink = connect.myDellink
}

//myDellink 删除回调
func (connect *ServerSession) myDellink(servertag int64) {

	mCheckSessionMap.Del(servertag)
	if connect.myDel != nil {
		connect.myDel(connect.servertag)
	}
}

//Start 开始连接
func (connect *ServerSession) Start() {
	connect.start()
	mCheckSessionMap.Put(connect)
}

//SetValid 设置是否为激活状态
func (connect *ServerSession) SetValid(bodata bool) {
	connect.bovalid = bodata
}

//SendCmd 发送数据包
func (connect *ServerSession) SendCmd(data interface{}) {
	connect.sendCmd(data)

}

//Close 关闭连接
func (connect *ServerSession) Close() {
	glog.Infof("主动调用关闭 %s,%d", connect.sname, connect.servertag)
	connect.close()
}
