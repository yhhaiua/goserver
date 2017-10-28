package gtcp

import (
	"net"

	"github.com/yhhaiua/goserver/common"
)

//ServerSession 请求连接结构
type ServerSession struct {
	*baseSession
}

//AddSession 添加请求信息
func AddSession(conn *net.TCPConn, backtype int32) *ServerSession {
	Session := new(ServerSession)
	Session.baseSession = addbase(conn, backtype)
	//包解析
	Session.newcodec(newcodecBinary)

	return Session
}

//Cmdcodec 包解析
func (connect *ServerSession) Cmdcodec() common.CmdCodec {
	return connect.cmdcodec
}

//SetFunc 发送验证包的函数、读取数据包的函数
func (connect *ServerSession) SetFunc(Queue func(pcmd *common.BaseCmd, data []byte) bool) {
	connect.msgQueue = Queue
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
