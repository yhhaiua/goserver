package logictime

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/yhhaiua/goserver/common/glog"
)

//Logicsvr 服务器数据
type Logicsvr struct {
	mstJSONConfig stJSONConfig
	recordtime    time.Time
	//twodaytime    time.Time
	//deltime       time.Time
}

var (
	instance *Logicsvr
	mu       sync.Mutex
)

//Instance 实例化logicsvr
func Instance() *Logicsvr {
	if instance == nil {
		mu.Lock()
		defer mu.Unlock()
		if instance == nil {
			instance = new(Logicsvr)
		}
	}
	return instance
}

//LogicInit 初始化
func (logic *Logicsvr) LogicInit() bool {

	//读取配置

	if logic.mstJSONConfig.configInit() {
		now := time.Now()
		logic.recordtime = now
		//logic.twodaytime = now
		//logic.deltime = now.Add(time.Hour * 1)
		logic.create()
		dir := "./save"
		os.MkdirAll(dir, os.ModePerm)
		go logic.timed()
		return true
	}
	return false
}

func (logic *Logicsvr) timeSub(t1, t2 time.Time) int {
	t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, time.Local)
	t2 = time.Date(t2.Year(), t2.Month(), t2.Day(), 0, 0, 0, 0, time.Local)

	return int(t1.Sub(t2).Hours())
}

func (logic *Logicsvr) timed() {
	for {
		now := time.Now()
		// 计算下一个零点
		next := now.Add(time.Hour * 1)
		next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), 0, 0, 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		glog.Infof("next save time:%s", next.String())
		<-t.C
		glog.Info("start save")
		logic.copy()
	}
}

func (logic *Logicsvr) create() {
	//now := time.Now()
	//	dir := fmt.Sprintf("./%04d%02d%02d%s",
	//		now.Year(),
	//		now.Month(),
	//		now.Day(),
	//		"saveone")
	os.MkdirAll("./saveone", os.ModePerm)
}
func (logic *Logicsvr) remove() {

	info := logic.dirents("./saveone/")
	if info != nil {
		if len(info) > 10 {
			dir := fmt.Sprintf("%s%s",
				"./saveone/",
				info[0].Name())
			err := os.Remove(dir)
			if err != nil {
				glog.Errorf("del %s error:%s", dir, err)
			} else {
				glog.Infof("del %s succ", dir)
			}

		}
	}

}
func (logic *Logicsvr) removeday() {

	info := logic.dirents("./save/")
	if info != nil {
		if len(info) > 10 {
			dir := fmt.Sprintf("%s%s",
				"./save/",
				info[0].Name())
			err := os.Remove(dir)
			if err != nil {
				glog.Errorf("del %s error:%s", dir, err)
			} else {
				glog.Infof("del %s succ", dir)
			}

		}
	}

}
func (logic *Logicsvr) dirents(dir string) []os.FileInfo {

	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		glog.Errorf("dirents error : %s\n", err)
		return nil
	}
	return entries
}
func (logic *Logicsvr) getfilecount(dir string) int {
	info := logic.dirents(dir)
	if info != nil {
		return len(info)
	}
	return 0
}
func (logic *Logicsvr) copy() {
	now := time.Now()
	daycount := logic.timeSub(now, logic.recordtime)
	if daycount >= 24 {
		logic.recordtime = now
		dir := fmt.Sprintf("%s/%04d%02d%02d%s",
			"./save",
			now.Year(),
			now.Month(),
			now.Day(),
			"dump.rdb")

		_, err := logic.copyFile(dir, "./dump.rdb")
		if err != nil {
			glog.Errorf("copy nextday error : %s\n", err)
		} else {
			glog.Infof("save next day succ")
			logic.removeday()
			//logic.create()
		}
	} else {
		//		wen := fmt.Sprintf("./%04d%02d%02d%s",
		//			now.Year(),
		//			now.Month(),
		//			now.Day(),
		//			"saveone")

		dir := fmt.Sprintf("%s/%04d%02d%02d%02d%s",
			"./saveone",
			now.Year(),
			now.Month(),
			now.Day(),
			now.Hour(),
			"dump.rdb")
		_, err := logic.copyFile(dir, "./dump.rdb")
		if err != nil {
			glog.Errorf("copy nexttime error : %s\n", err)
		} else {
			glog.Infof("save next time succ")
			logic.remove()
		}
	}

}

func (logic *Logicsvr) copyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return
	}
	defer dst.Close()
	return io.Copy(dst, src)
}
