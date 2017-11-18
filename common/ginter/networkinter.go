package ginter

import (
	"github.com/yhhaiua/goserver/common"
	"github.com/yhhaiua/goserver/common/gpacket"
	"github.com/yhhaiua/goserver/protocol"
)

//////////////////////////////外包类需要实现的接口///////////////////////////////////////////

//NetWorkAgenter 需要实现的接口
type NetWorkAgenter interface {
	MsgQueue(*gpacket.BaseCmd, []byte) bool
	CloseLink(int64)
	StartLink(int64)
	CmdCodec() common.CmdCodec
}

//SessionAgenter 接听回调需要实现的接口
type SessionAgenter interface {
	NetWorkAgenter
}

//ConnectAgenter 连接回调需要实现的接口
type ConnectAgenter interface {
	NetWorkAgenter
	SendOnceCmd()
}

//////////////////////////////网络接口/////////////////////////////////////////////////////

//NetWorker 需要的接口
type NetWorker interface {
	Start()
	Close()
	SetValid(bodata bool)
	SendCmd(data interface{})
	Setheartbeat(pcmd *protocol.ServerCmdHeart)
}
