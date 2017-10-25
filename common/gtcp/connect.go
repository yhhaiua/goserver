package gtcp

import (
	"bytes"
	"encoding/binary"
	"net"

	"github.com/yhhaiua/goserver/common"

	"github.com/yhhaiua/goserver/common/glog"
)

//ClientConnecter 请求连接结构
type ClientConnecter struct {
	myTCPAddr   *net.TCPAddr
	nServerID   int32
	bovalid     bool
	boConnected bool
	conn        *net.TCPConn
	mrecvMybuf  loopBuf
	sendMybuf   loopBuf
	MsgQueue    func(pcmd *BaseCmd, data []byte)
	SendOnce    func()
}

const (
	maxSendbufLen   = 1024 * 4        //一次发送长度
	maxforcedbufLen = 1024 * 1024 * 5 //强制发送长度
)

//AddConnect 添加请求信息
func AddConnect(serverip, port string, serverid int32) *ClientConnecter {
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
	//压人请求连接队列
	mTCPConnMap.Put(Connecter)

	return Connecter
}

//SetFunc 发送验证包的函数、读取数据包的函数
func (connect *ClientConnecter) SetFunc(Queue func(pcmd *BaseCmd, data []byte), Once func()) {
	connect.MsgQueue = Queue
	connect.SendOnce = Once
}
func (connect *ClientConnecter) runRead() {
	tempbuf := make([]byte, 65536)

	for {
		if connect.boConnected && connect.conn != nil {
			len, err := connect.conn.Read(tempbuf)
			if err != nil {
				connect.doClose()
				glog.Errorf("socket连接断开 %d,%s", connect.nServerID, err)
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
func (connect *ClientConnecter) doClose() {
	connect.conn.Close()
	connect.bovalid = false
	connect.boConnected = false

}

func (connect *ClientConnecter) doInit() {
	connect.mrecvMybuf.newLoopBuf(2048)
	connect.sendMybuf.newLoopBuf(2048)
	connect.boConnected = true
}

func (connect *ClientConnecter) doRead() bool {

	for {
		if connect.boConnected && connect.mrecvMybuf.canreadlen >= 8 {

			tembuf := connect.mrecvMybuf.buf[connect.mrecvMybuf.getreadadd():connect.mrecvMybuf.getreadlenadd()]

			bytesBuffer := bytes.NewBuffer(tembuf[:8])
			var packet PacketBase
			binary.Read(bytesBuffer, binary.BigEndian, &packet)
			if packet.Size >= 1024*64 || packet.Size < 2 {
				glog.Errorf("收到恶意攻击包 %d,%d", connect.nServerID, packet.Size)
				return false
			}
			newlen := common.Alignment(int(packet.Size+6), 8)

			if connect.mrecvMybuf.canreadlen >= newlen {

				//包处理
				connect.MsgQueue(&packet.Pcmd, tembuf[6:packet.Size+6])

				connect.mrecvMybuf.setReadPtr(newlen)
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
func (connect *ClientConnecter) SendCmd(data interface{}) {

	if connect.boConnected && connect.bovalid {

		var packet Packet
		packet.Size = uint32(binary.Size(data))
		packet.data = data
		bytesBuffer := new(bytes.Buffer)
		binary.Write(bytesBuffer, binary.LittleEndian, &packet)

		connect.sendMybuf.addSendBuf(bytesBuffer.Bytes(), bytesBuffer.Len())

		if connect.sendMybuf.canreadlen >= maxforcedbufLen {
			//强制发送一次
			connect.startSend()
		}
	}

}

func (connect *ClientConnecter) startSend() {

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

func (connect *ClientConnecter) runWrite() {

	for {
		if connect.boConnected && connect.conn != nil {
			connect.startSend()
		}
	}

}
