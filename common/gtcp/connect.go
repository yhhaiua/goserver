package gtcp

import (
	"bytes"
	"encoding/binary"
	"net"

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
}

//AddConnect 添加请求信息
func AddConnect(serverip, port string, serverid int32) {
	Connecter := new(ClientConnecter)
	connectadd := serverip + ":" + port
	Connecter.nServerID = serverid
	var err error
	Connecter.myTCPAddr, err = net.ResolveTCPAddr("tcp", connectadd)
	if err != nil {
		glog.Errorf("AddConnect error:%s", err)
		Connecter = nil
		return
	}
	//压人请求连接队列
	mTCPConnMap.Put(Connecter)
}

func (connect *ClientConnecter) sendVerificationCmd() {

}
func (connect *ClientConnecter) runRead() {
	tempbuf := make([]byte, 65536)

	connect.mrecvMybuf.newLoopBuf(2048)

	for {
		if connect.boConnected && connect.conn != nil {
			len, err := connect.conn.Read(tempbuf)
			if err != nil || len == 0 {
				connect.boConnected = false
				connect.bovalid = false
				glog.Errorf("socket连接断开 %d,%s", connect.nServerID, err)
				return
			}
			connect.mrecvMybuf.putData(tempbuf, len)
			//处理包
			if !connect.doRead() {
				connect.boConnected = false
				connect.bovalid = false
				connect.conn.Close()
				return
			}
		}
	}

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
			newlen := alignment(int(packet.Size+6), 8)

			if connect.mrecvMybuf.canreadlen >= newlen {

				//包处理
				pushMsgQueue(&packet.Pcmd, tembuf[6:packet.Size+6])

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
