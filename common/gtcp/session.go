package gtcp

import (
	"net"

	"github.com/yhhaiua/goserver/common/ginter"
	"github.com/yhhaiua/goserver/common/glog"
)

//ServerSession 请求连接结构
type ServerSession struct {
	*baseSession
}

//AddSession 添加请求信息
func AddSession(conn *net.TCPConn, backtype int64, sname string, agent ginter.SessionAgenter) *ServerSession {
	Session := new(ServerSession)
	Session.baseSession = addbase(conn, backtype, sname, sessionbaseType, agent)

	glog.Infof("有新的连接进入 %s,%d", sname, backtype)
	return Session
}

//Start 开始连接
func (connect *ServerSession) Start() {
	connect.start()

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
