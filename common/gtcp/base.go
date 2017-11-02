package gtcp

import (
	"net"
	"sync"
	"time"

	"github.com/yhhaiua/goserver/protocol"

	"github.com/yhhaiua/goserver/common"

	"github.com/yhhaiua/goserver/common/glog"
)

const (
	newcodecBinary = 1
)
const (
	maxSendbufLen    = 1024 * 4 //一次发送长度
	initMybufLen     = 2048     //初始化buf长度
	verificationtime = 30       //连接验证时间
)

const (
	checkstatusNo   = 0
	checkstatusWait = 1
)

type heartbeat struct {
	checkstatus    int32
	checktime      int64
	checklock      sync.Mutex
	checkTimeStart int32
	checkTimeWait  int32
}

func (myheart *heartbeat) myInit() {
	myheart.checkTimeStart = 120
	myheart.checkTimeWait = 20
	myheart.checktime = time.Now().Unix() + int64(myheart.checkTimeStart)
}

//BaseSession 连接结构
type baseSession struct {
	servertag   int64
	bovalid     bool
	boConnected bool
	boRecon     bool
	conn        *net.TCPConn
	mrecvMybuf  loopBuf
	sendMybuf   loopBuf
	msgQueue    func(pcmd *common.BaseCmd, data []byte) bool
	cmdcodec    common.CmdCodec
	sname       string
	starttime   int64
	delLink     func(servertag int64)
	endSync     sync.WaitGroup
	myHeart     heartbeat
	myHeartFunc func()
}

func addbase(conn *net.TCPConn, servertag int64, sname string) *baseSession {
	Session := new(baseSession)
	Session.conn = conn
	Session.servertag = servertag
	Session.sname = sname
	Session.sendMybuf.initSendBuf()
	Session.myHeart.myInit()
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
	connect.endSync.Add(2)

	go connect.waitClose()

	go connect.runRead()
	go connect.runWrite()
}
func (connect *baseSession) waitClose() {
	connect.endSync.Wait()
	connect.doClose()
}
func (connect *baseSession) runRead() {
	tempbuf := make([]byte, 65536)

	for {
		if connect.boConnected && connect.conn != nil {
			len, err := connect.conn.Read(tempbuf)
			if err != nil {
				connect.close()
				glog.Errorf("socket连接断开 %s,%d,%s", connect.sname, connect.servertag, err)
				break
			}
			connect.mrecvMybuf.putData(tempbuf, len, len)
			//处理包
			if !connect.doRead() {
				connect.close()
				break
			}
		} else {
			break
		}
	}
	connect.endSync.Done()

}
func (connect *baseSession) doClose() {
	connect.bovalid = false
	connect.boRecon = true
	glog.Warningf("连接关闭%s,%d", connect.sname, connect.servertag)
	if connect.delLink != nil {
		connect.delLink(connect.servertag)
	}
}

func (connect *baseSession) close() {
	connect.boConnected = false
	connect.sendMybuf.SendCond.Signal()
}
func (connect *baseSession) isValid() bool {
	return connect.bovalid
}
func (connect *baseSession) doInit() {
	glog.Warningf("连接成功 %d,%s", connect.servertag, connect.sname)
	connect.mrecvMybuf.newLoopBuf(initMybufLen)
	connect.sendMybuf.newLoopBuf(initMybufLen)
	connect.boConnected = true
	connect.boRecon = false
	connect.starttime = time.Now().Unix()
}

func (connect *baseSession) doRead() bool {

	for {
		if connect.boConnected && connect.mrecvMybuf.canreadlen >= 8 {

			tembuf := connect.mrecvMybuf.buf[connect.mrecvMybuf.getreadadd():connect.mrecvMybuf.getreadlenadd()]
			var packet common.PacketBase

			err := connect.cmdcodec.Decode(tembuf[:8], &packet)

			if err != nil {
				glog.Errorf("收到恶意攻击包%s,%d,%s", connect.sname, connect.servertag, err)
				return false
			}
			if packet.Size >= 1024*64 || packet.Size < 2 {
				glog.Errorf("收到恶意攻击包 %s,%d,%d", connect.sname, connect.servertag, packet.Size)
				return false
			}

			newlen := common.Alignment(int(packet.Size+6), 8)

			if connect.mrecvMybuf.canreadlen >= newlen {

				//包处理
				if connect.msgQueue(&packet.Pcmd, tembuf[6:packet.Size+6]) {
					connect.mrecvMybuf.setReadPtr(newlen)
				} else {
					return false
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
			glog.Errorf("data err:%s,%d,%s", connect.sname, connect.servertag, err)
			return
		}
		connect.sendMybuf.addSendBuf(bytedata, len(bytedata))
	}

}

func (connect *baseSession) startSend() {

	connect.sendMybuf.Sendlock.Lock()
	if connect.sendMybuf.canreadlen == 0 {
		connect.sendMybuf.SendCond.Wait()
	}
	connect.sendMybuf.Sendlock.Unlock()

	connect.sendMybuf.Sendlock.Lock()
	for {
		sendlen := common.Min(connect.sendMybuf.canreadlen, maxSendbufLen)

		if sendlen > 0 {
			tembuf := connect.sendMybuf.buf[connect.sendMybuf.getreadadd() : connect.sendMybuf.getreadadd()+sendlen]
			writelen, err := connect.conn.Write(tembuf)
			if err == nil {
				connect.sendMybuf.setReadPtr(writelen)
			} else {
				glog.Errorf("写入错误%s,%d,%s", connect.sname, connect.servertag, err)
				break
			}
		} else {
			break
		}
	}
	connect.sendMybuf.Sendlock.Unlock()
}

func (connect *baseSession) runWrite() {

	for {
		if connect.boConnected && connect.conn != nil {
			connect.startSend()

			//time.Sleep(time.Millisecond * 1)
		} else {
			break
		}
	}
	connect.conn.Close()
	connect.endSync.Done()

}

func (connect *baseSession) runCheck() {
	if connect.boConnected && connect.conn != nil {

		if !connect.isValid() {
			nowtime := time.Now().Unix()
			if nowtime-connect.starttime > verificationtime {
				glog.Errorf("%d秒连接验证超时%s,%d", verificationtime, connect.sname, connect.servertag)
				connect.close()
			}
		} else {
			if connect.myHeartFunc != nil {
				connect.myHeartFunc()
			} else {
				//统一的心跳函数
				connect.heartbeatCheck()
			}

		}
	}
}

func (connect *baseSession) heartbeatCheck() {
	connect.myHeart.checklock.Lock()
	nowtime := time.Now().Unix()
	if nowtime > connect.myHeart.checktime {
		switch connect.myHeart.checkstatus {
		case checkstatusNo:
			var retcmd protocol.ServerCmdHeart
			retcmd.Init()
			retcmd.IsneedAck = true
			connect.sendCmd(&retcmd)
			connect.myHeart.checkstatus = checkstatusWait
			connect.myHeart.checktime = nowtime + int64(connect.myHeart.checkTimeWait)
			glog.Infof("发送心跳包 %s,%d", connect.sname, connect.servertag)
		case checkstatusWait:
			glog.Infof("没有收到心跳包关闭 %s,%d", connect.sname, connect.servertag)
			connect.close()
		}
	}
	connect.myHeart.checklock.Unlock()

}
func (connect *baseSession) Setheartbeat(pcmd *protocol.ServerCmdHeart) {
	glog.Infof("收到心跳包 %s,%d", connect.sname, connect.servertag)
	if pcmd.IsneedAck {
		pcmd.IsneedAck = false
		connect.sendCmd(pcmd)
	}
	if connect.myHeartFunc == nil {

		nowtime := time.Now().Unix()
		connect.myHeart.checklock.Lock()
		connect.myHeart.checkstatus = checkstatusNo
		connect.myHeart.checktime = nowtime + int64(connect.myHeart.checkTimeStart)
		connect.myHeart.checklock.Unlock()
	}
}
