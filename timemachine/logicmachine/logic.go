package logicmachine

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/yhhaiua/goserver/common/glog"
)

//Logicsvr 服务器数据
type Logicsvr struct {
	morningtime time.Time
	nighttime   time.Time
	count       int
}

var (
	instance *Logicsvr
	mu       sync.Mutex
)

type mycontent struct {
	Content string `json:"content"`
}
type isAtdata struct {
	IsAtAll bool `json:"isAtAll"`
}
type senddata struct {
	Msgtype string    `json:"msgtype"`
	Text    mycontent `json:"text"`
	At      isAtdata  `json:"at"`
}

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

	//logic.testsend("测试，测试")
	logic.count = 1
	go logic.check()
	return true
}

func (logic *Logicsvr) check() {
	for {
		now := time.Now()
		// 计算下一个零点
		next := now.Add(time.Minute * 1)
		next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), next.Minute(), 0, 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		glog.Infof("next check time:%s", next.String())
		<-t.C
		if next.Weekday() == 1 {
			if next.Hour() == 9 && next.Minute() == 1 {
				if logic.count == 0 {
					logic.count = 1
				} else {
					logic.count = 0
				}

				if logic.count == 0 {
					logic.send("大家好，很荣幸为大家服务!本周服务时间周一到周五")
				} else {
					logic.send("大家好，很荣幸为大家服务!本周服务时间周一到周六")
				}
			}
		}
		if logic.count == 0 {
			if next.Weekday() == 0 || next.Weekday() == 6 {
				continue
			}
		} else {
			if next.Weekday() == 0 {
				continue
			}
		}
		if next.Weekday() != 0 {
			if next.Hour() == 9 && next.Minute() >= 20 && next.Minute() <= 30 {

				logic.send("注意注意，本神探提醒诸位早上时间可以打卡了!")
			}

			if next.Hour() == 21 && next.Minute() >= 30 && next.Minute() <= 50 {
				logic.send("注意注意，本神探提醒诸位晚上时间可以打卡了!")
			}
		}
	}
}

func (logic *Logicsvr) send(src string) {

	var data senddata
	data.Msgtype = "text"
	data.Text.Content = src
	data.At.IsAtAll = true
	b, err := json.Marshal(data)
	if err != nil {
		glog.Errorf("json:%s", err)
		return
	}
	glog.Infof("send content :%s", string(b))

	body := bytes.NewBuffer(b)
	res, err := http.Post("https://oapi.dingtalk.com/robot/send?access_token=2d3cabf7cf6dc79ada9f31658effc915bb976caaa788604d7aaebb0cea81ac6b", "application/json;charset=utf-8", body)
	if err != nil {
		glog.Errorf("error1:%s", err)
		return
	}
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		glog.Errorf("error2:%s", err)
		return
	}
	glog.Infof("%s", result)
}

func (logic *Logicsvr) testsend(src string) {

	var data senddata
	data.Msgtype = "text"
	data.Text.Content = src
	data.At.IsAtAll = true
	b, err := json.Marshal(data)
	if err != nil {
		glog.Errorf("json:%s", err)
		return
	}
	glog.Infof("send content :%s", string(b))

	body := bytes.NewBuffer(b)
	res, err := http.Post("https://oapi.dingtalk.com/robot/send?access_token=2eb8253aae5237588004af68512f5fa6205fe2f6b4f08fc15d603287e0376d40", "application/json;charset=utf-8", body)
	if err != nil {
		glog.Errorf("error1:%s", err)
		return
	}
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		glog.Errorf("error2:%s", err)
		return
	}
	glog.Infof("%s", result)
}
