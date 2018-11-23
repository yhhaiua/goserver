package logic

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/yhhaiua/goserver/common/grouter"
	"github.com/yhhaiua/goserver/common/log4go"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"
)

type stRouterPost struct {
	wanbaDeal stWanba
	increaseId    int64
}

//机器人信息
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

//充值错误保存信息
type stRecharge struct {
	OrderId string
	Money string
	GoodsId string
	PlayerId string
	ErrorInfo string
	Operatorid string
}
func (myrouter *stRouterPost)getIncreaseID()int64  {
	return  atomic.AddInt64(&myrouter.increaseId,1)
}

func (myrouter *stRouterPost) rechargeDeal(w http.ResponseWriter, r *http.Request, _ grouter.Params) {

	operatorid := r.FormValue("operatorid")
	if operatorid != Instance().mstJSONConfig.operatorid{
		myrouter.send(w,-100,"operatorid error")
		return
	}
	if Instance().stopRecharge{
		myrouter.send(w,-100,"stop charge")
		return
	}
	switch operatorid {
	case "1":
		//玩吧渠道
		myrouter.wanbaDeal.rechargeDeal(w,r,myrouter)
	default:
		log4go.Error("错误渠道请求:%s",operatorid)
		myrouter.send(w,-100," operatorid no have error")
	}
}

func (router *stRouterPost)send(w http.ResponseWriter,ret int,msg string)  {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	var info RetInfo
	info.Ret = ret
	info.Msg = msg
	Message, err := json.Marshal(info)

	if err == nil {
		fmt.Fprintf(w, "%s", Message)
	}
}

//生成订单号
func (myrouter *stRouterPost)createorderId(playerId,str string) string {

	timestamp := time.Now().Unix()
	stime := strconv.FormatInt(timestamp, 10)
	var buffer bytes.Buffer
	buffer.WriteString(str)
	buffer.WriteString(playerId)
	buffer.WriteString(stime)
	buffer.WriteString( strconv.FormatInt(myrouter.getIncreaseID(),10))
	return buffer.String()
}
func (myrouter *stRouterPost)parseServer(playerId string) string  {
	pid,_:= strconv.ParseInt(playerId, 10, 64)
	serverid:= (pid % 1000000) / 100
	return  strconv.FormatInt(serverid, 10)
}
//像gm后台发送
func (myrouter *stRouterPost) sendgm(routers,orderId,money,goodsId,playerId string) (int,string) {
	success:= -1
	errorStr := orderId+":error"
	//playerId+商品id+时间+订单号+key
	timestamp := time.Now().Unix()
	severid := myrouter.parseServer(playerId)
	stime := strconv.FormatInt(timestamp, 10)
	md5str := playerId + goodsId + stime + orderId + Instance().mstJSONConfig.rechargekey
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(md5str))
	cipherStr := md5Ctx.Sum(nil)
	mysigon := hex.EncodeToString(cipherStr)

	var buffer bytes.Buffer
	buffer.WriteString("http://")
	buffer.WriteString(Instance().mstJSONConfig.gmhost)
	buffer.WriteString(routers)
	buffer.WriteString("?orderId=")
	buffer.WriteString(orderId)
	buffer.WriteString("&serverId=")
	buffer.WriteString(severid)
	buffer.WriteString("&money=")
	buffer.WriteString(money)
	buffer.WriteString("&goodsId=")
	buffer.WriteString(goodsId)
	buffer.WriteString("&playerId=")
	buffer.WriteString(playerId)
	buffer.WriteString("&time=")
	buffer.WriteString(stime)
	buffer.WriteString("&sign=")
	buffer.WriteString(mysigon)
	log4go.Info(buffer.String())

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := client.Get(buffer.String())
	if err == nil {
		defer req.Body.Close()
		body, err := ioutil.ReadAll(req.Body)
		if err == nil {
			valueRet := string(body)
			if valueRet == "ok"{
				success = 0
				errorStr = orderId+":ok"
				log4go.Info(errorStr)
			}else{
				log4go.Error("订单号:orderId:%s,充值错误返回:%s",orderId,valueRet)
				errorStr = orderId+":"+valueRet
				if valueRet == "-2"{
					success = -1
				}else{
					success = 1
				}
			}
		}else{
			log4go.Error("sendgm error2 订单号:orderId:%s,error:%s",orderId,err)
			success = 1
			errorStr = orderId+":-10002"
		}
	}else{
		log4go.Error("sendgm error1 订单号:orderId:%s,error:%s",orderId,err)
		success = -1
		errorStr = orderId+":-10001"
	}

	return success,errorStr
}
func (myrouter *stRouterPost) errorSave(orderId,money,goodsId,playerId,erorinfo string)  {
	var charge stRecharge
	charge.GoodsId = goodsId
	charge.Money = money
	charge.ErrorInfo = erorinfo
	charge.Operatorid = Instance().mstJSONConfig.operatorid
	charge.OrderId = orderId
	charge.PlayerId = playerId

	RedisMessage, err := json.Marshal(charge)
	if err == nil {
		log4go.Info("保存错误订单数据%s",RedisMessage)
		Key:= "ReCharge:" + orderId
		err = Instance().redisdb().Set(Key, RedisMessage)
		if err != nil {
			log4go.Error("errorSave error1:%s",err)
		}
	}
}

func (myrouter *stRouterPost) testsend(src string) {

	var data senddata
	data.Msgtype = "text"
	data.Text.Content = src
	data.At.IsAtAll = true
	b, err := json.Marshal(data)
	if err != nil {
		log4go.Error("json:%s", err)
		return
	}
	log4go.Info("send content :%s", string(b))

	body := bytes.NewBuffer(b)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	res, err := client.Post("https://oapi.dingtalk.com/robot/send?access_token=2eb8253aae5237588004af68512f5fa6205fe2f6b4f08fc15d603287e0376d40", "application/json;charset=utf-8", body)
	if err != nil {
		log4go.Error("testsend error1:%s", err)
		return
	}
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log4go.Error("testsend error2:%s", err)
		return
	}
	log4go.Info("testsend :%s", result)
}