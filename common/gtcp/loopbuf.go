package gtcp

import (
	"sync"

	"github.com/yhhaiua/goserver/common"
)

const (
	initMybufLen = 2048      //初始化buf长度
	maxMybufLen  = 64 * 1024 //最大buf长度
)

type loopBuf struct {
	buf         []byte     //包内容
	bufsize     int        //buff的最大长度
	canwritelen int        //能够写入的长度
	writeadd    int        //写地址
	canreadlen  int        //能够读取的长度
	readadd     int        //读地址
	freedatalen int        //空闲数据长度
	Sendlock    sync.Mutex //发送锁
	SendCond    *sync.Cond //发送信号
}

//新建一个buff缓存
func (loop *loopBuf) newLoopBuf(nmaxlen int) {

	loop.buf = make([]byte, nmaxlen)
	loop.bufsize = nmaxlen
	loop.readadd = 0
	loop.writeadd = 0
	loop.canreadlen = 0
	loop.canwritelen = loop.bufsize
	loop.freedatalen = 0
}

//向buff中压人数据
func (loop *loopBuf) putData(data []byte, len int, datelen int) {

	if len <= 0 {
		return
	}
	if loop.canwritelen < len {

		if loop.canwritelen+loop.freedatalen > 2*len {
			//自己内存挪动
			loop.moveData()
		} else {
			//开辟更大的内存
			loop.distributionData(len)
		}
	}
	loop.putRightData(data, len, datelen)
}

//压人数据
func (loop *loopBuf) putRightData(data []byte, len int, datelen int) {

	copy(loop.buf[loop.writeadd:], data[:datelen])
	loop.writeadd += len
	loop.canwritelen -= len
	loop.canreadlen += len
}

//自身数据腾挪
func (loop *loopBuf) moveData() {

	copy(loop.buf[:], loop.buf[loop.readadd:loop.readadd+loop.canreadlen])
	loop.readadd = 0
	loop.writeadd = loop.canreadlen
	loop.freedatalen = 0
	loop.canwritelen = loop.bufsize - loop.canreadlen
}

//重新开辟更大空间
func (loop *loopBuf) distributionData(len int) {
	newlen := common.Alignment(loop.bufsize+len+loop.bufsize, 1024)

	temp := loop.buf[loop.readadd : loop.readadd+loop.canreadlen]

	loop.buf = make([]byte, newlen)
	copy(loop.buf, temp)

	loop.bufsize = newlen
	loop.canwritelen = loop.bufsize - loop.canreadlen
	loop.freedatalen = 0
	loop.readadd = 0
	loop.writeadd = loop.canreadlen
}

//向buff中释放数据
func (loop *loopBuf) setReadPtr(nlen int) {
	if nlen <= 0 || loop.canreadlen <= 0 {
		return
	}
	if loop.canreadlen <= nlen {
		//正好读完数据
		loop.readadd = 0
		loop.writeadd = 0
		loop.canreadlen = 0
		loop.canwritelen = loop.bufsize
		loop.freedatalen = 0
		if loop.bufsize >= maxMybufLen {
			//缓存过大直接重置
			loop.newLoopBuf(initMybufLen)
		}
		return
	}
	loop.readadd += nlen
	loop.freedatalen += nlen
	loop.canreadlen -= nlen
}

//读取数据地址
func (loop *loopBuf) getreadadd() int {
	return loop.readadd
}

//读取数据末尾地址
func (loop *loopBuf) getreadlenadd() int {
	return loop.readadd + loop.canreadlen
}

//添加发送数据
func (loop *loopBuf) addSendBuf(data []byte, len int) {
	loop.Sendlock.Lock()
	loop.putData(data, common.Alignment(len, 8), len)
	loop.Sendlock.Unlock()

	loop.SendCond.Signal()
}

//初始化发送
func (loop *loopBuf) initSendBuf() {
	loop.SendCond = sync.NewCond(&loop.Sendlock)
}
