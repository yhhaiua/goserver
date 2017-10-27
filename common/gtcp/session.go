package gtcp

import (
	"net"

	"github.com/yhhaiua/goserver/common"

	"github.com/yhhaiua/goserver/common/glog"
)

//ServerSession 请求连接结构
type ServerSession struct {
	backtype    int32
	bovalid     bool
	boConnected bool
	conn        *net.TCPConn
	mrecvMybuf  loopBuf
	sendMybuf   loopBuf
	MsgQueue    func(pcmd *common.BaseCmd, data []byte) bool

	Cmdcodec common.CmdCodec
}

//AddSession 添加请求信息
func AddSession(conn *net.TCPConn, backtype int32) *ServerSession {
	Session := new(ServerSession)
	Session.conn = conn
	Session.backtype = backtype
	//包解析
	Session.Cmdcodec = new(common.BinaryCodec)

	return Session
}

//SetFunc 发送验证包的函数、读取数据包的函数
func (connect *ServerSession) SetFunc(Queue func(pcmd *common.BaseCmd, data []byte) bool) {
	connect.MsgQueue = Queue
}

//Start 开始连接
func (connect *ServerSession) Start() {
	connect.doInit()
	go connect.runRead()
	go connect.runWrite()
}
func (connect *ServerSession) runRead() {
	tempbuf := make([]byte, 65536)

	for {
		if connect.boConnected && connect.conn != nil {
			len, err := connect.conn.Read(tempbuf)
			if err != nil {
				connect.doClose()
				glog.Errorf("socket连接断开 %d,%s", connect.backtype, err)
				return
			}
			connect.mrecvMybuf.putData(tempbuf, len, len)
			//处理包
			if !connect.doRead() {
				connect.doClose()
				return
			}
		}
	}

}
func (connect *ServerSession) doClose() {
	connect.conn.Close()
	connect.bovalid = false
	connect.boConnected = false

}

//SetValid 设置是否为激活状态
func (connect *ServerSession) SetValid(bodata bool) {
	connect.bovalid = bodata
}
func (connect *ServerSession) doInit() {
	glog.Infof("ServerSession 连接成功 %d", connect.backtype)
	connect.mrecvMybuf.newLoopBuf(2048)
	connect.sendMybuf.newLoopBuf(2048)
	connect.boConnected = true
}

func (connect *ServerSession) doRead() bool {

	for {
		if connect.boConnected && connect.mrecvMybuf.canreadlen >= 8 {

			tembuf := connect.mrecvMybuf.buf[connect.mrecvMybuf.getreadadd():connect.mrecvMybuf.getreadlenadd()]
			var packet common.PacketBase

			err := connect.Cmdcodec.Decode(tembuf[:8], &packet)

			if err != nil {
				glog.Errorf("收到恶意攻击包 %s", err)
				return false
			}
			if packet.Size >= 1024*64 || packet.Size < 2 {
				glog.Errorf("收到恶意攻击包 %d,%d", connect.backtype, packet.Size)
				return false
			}

			newlen := common.Alignment(int(packet.Size+6), 8)

			if connect.mrecvMybuf.canreadlen >= newlen {

				//包处理
				if connect.MsgQueue(&packet.Pcmd, tembuf[6:packet.Size+6]) {
					connect.mrecvMybuf.setReadPtr(newlen)
				} else {
					break
				}

			} else {
				break
			}
		} else {
			break
		}
	}

	return true
}

//SendCmd 发送数据包
func (connect *ServerSession) SendCmd(data interface{}) {

	if connect.boConnected {

		var packet common.Packet
		packet.Size = uint32(connect.Cmdcodec.Size(data))
		packet.Data = data
		bytedata, err := connect.Cmdcodec.Encode(&packet)
		if err != nil {
			glog.Errorf("data err:%s", err)
			return
		}
		connect.sendMybuf.addSendBuf(bytedata, len(bytedata))

		if connect.sendMybuf.canreadlen >= maxforcedbufLen {
			//强制发送一次
			connect.startSend()
		}
	}

}

func (connect *ServerSession) startSend() {

	connect.sendMybuf.Sendlock.Lock()
	defer connect.sendMybuf.Sendlock.Unlock()

	if connect.sendMybuf.canreadlen > 0 {
		for {
			sendlen := common.Min(connect.sendMybuf.canreadlen, maxSendbufLen)

			if sendlen > 0 {
				tembuf := connect.sendMybuf.buf[connect.sendMybuf.getreadadd() : connect.sendMybuf.getreadadd()+sendlen]
				writelen, err := connect.conn.Write(tembuf)
				if err == nil {
					connect.sendMybuf.setReadPtr(writelen)
				} else {
					glog.Errorf("写入错误%s", err)
					break
				}
			} else {
				break
			}
		}

	}
}

func (connect *ServerSession) runWrite() {

	for {
		if connect.boConnected && connect.conn != nil {
			connect.startSend()
		}
	}

}
