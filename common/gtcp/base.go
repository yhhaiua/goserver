package gtcp

import (
	"net"
	"time"

	"github.com/yhhaiua/goserver/common"

	"github.com/yhhaiua/goserver/common/glog"
)

const (
	newcodecBinary = 1
)
const (
	maxSendbufLen   = 1024 * 4        //一次发送长度
	maxforcedbufLen = 1024 * 1024 * 5 //强制发送长度
)

//baseSession 连接结构
type baseSession struct {
	servertag   int32
	bovalid     bool
	boConnected bool
	conn        *net.TCPConn
	mrecvMybuf  loopBuf
	sendMybuf   loopBuf
	msgQueue    func(pcmd *common.BaseCmd, data []byte) bool
	cmdcodec    common.CmdCodec
}

func addbase(conn *net.TCPConn, servertag int32) *baseSession {
	Session := new(baseSession)
	Session.conn = conn
	Session.servertag = servertag
	return Session
}
func (connect *baseSession) newcodec(codectype int) {
	switch codectype {
	case newcodecBinary:
		connect.cmdcodec = new(common.BinaryCodec)
	default:
	}
}

//start 开始连接
func (connect *baseSession) start() {
	connect.doInit()
	go connect.runRead()
	go connect.runWrite()
}
func (connect *baseSession) runRead() {
	tempbuf := make([]byte, 65536)

	for {
		if connect.boConnected && connect.conn != nil {
			len, err := connect.conn.Read(tempbuf)
			if err != nil {
				connect.doClose()
				glog.Errorf("socket连接断开 %d,%s", connect.servertag, err)
				return
			}
			connect.mrecvMybuf.putData(tempbuf, len, len)
			//处理包
			if !connect.doRead() {
				connect.doClose()
				return
			}
		} else {
			break
		}
	}

}
func (connect *baseSession) doClose() {
	connect.conn.Close()
	connect.bovalid = false
	connect.boConnected = false

}

func (connect *baseSession) doInit() {
	glog.Infof("ServerSession 连接成功 %d", connect.servertag)
	connect.mrecvMybuf.newLoopBuf(2048)
	connect.sendMybuf.newLoopBuf(2048)
	connect.boConnected = true
}

func (connect *baseSession) doRead() bool {

	for {
		if connect.boConnected && connect.mrecvMybuf.canreadlen >= 8 {

			tembuf := connect.mrecvMybuf.buf[connect.mrecvMybuf.getreadadd():connect.mrecvMybuf.getreadlenadd()]
			var packet common.PacketBase

			err := connect.cmdcodec.Decode(tembuf[:8], &packet)

			if err != nil {
				glog.Errorf("收到恶意攻击包 %s", err)
				return false
			}
			if packet.Size >= 1024*64 || packet.Size < 2 {
				glog.Errorf("收到恶意攻击包 %d,%d", connect.servertag, packet.Size)
				return false
			}

			newlen := common.Alignment(int(packet.Size+6), 8)

			if connect.mrecvMybuf.canreadlen >= newlen {

				//包处理
				if connect.msgQueue(&packet.Pcmd, tembuf[6:packet.Size+6]) {
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

//sendCmd 发送数据包
func (connect *baseSession) sendCmd(data interface{}) {

	if connect.boConnected {

		var packet common.Packet
		packet.Size = uint32(connect.cmdcodec.Size(data))
		packet.Data = data
		bytedata, err := connect.cmdcodec.Encode(&packet)
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

func (connect *baseSession) startSend() {

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

func (connect *baseSession) runWrite() {

	for {
		if connect.boConnected && connect.conn != nil {
			connect.startSend()
			time.Sleep(time.Millisecond * 1)
		} else {
			break
		}
	}

}
