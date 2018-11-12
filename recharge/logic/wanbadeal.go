package logic

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/yhhaiua/goserver/common/gjson"
	"github.com/yhhaiua/goserver/common/glog"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// RetInfo 错误返回
type RetInfo struct {
	Ret int `json:"ret"`
	Msg string `json:"msg"`
}
type stWanba struct {

}

//urlencode
func (router *stWanba)encodeURIComponent(str string) string {
	r := url.QueryEscape(str)
	r = strings.Replace(r, "+", "%20", -1)
	return r
}

//生成sig
func (router *stWanba)sigCreate(urlinit,urlstr string) string {
	str0 :="POST"
	str1:= router.encodeURIComponent(urlinit)
	str2 := router.encodeURIComponent(urlstr)
	str3 := str0+"&"+str1+"&"+str2
	key :=  []byte("DWB13t84CoEL8eax&")
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(str3))
	ucnc := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	fmt.Println(ucnc)
	return router.encodeURIComponent(ucnc)
}

//func (router *stWanba) urlencode(urlStr string)  {
//	l3, err3 := url.Parse(urlStr)
//}
//玩吧处理
func (router *stWanba) rechargeDeal(w http.ResponseWriter, r *http.Request,myrouter *stRouterPost) ()  {

	ok,money := router.inquiryRet(w,r,myrouter)
	if ok{
		playerid := r.FormValue("playerid")

		if(Instance().checkBuyMapKey(playerid)){
			myrouter.send(w,-100,"有订单正在处理")
			return
		}

		Instance().addBuyMap(playerid)

		billno := myrouter.createorderId(playerid,"wb")

		ok = router.deductionRet(w,r,myrouter,money,billno)
		if ok{
			router.sendGm(w,r,myrouter,money,billno)
			myrouter.send(w,0,"sucess")
		}else{
			Instance().delBuyMap(playerid)
		}
	}

}

//判断玩家金钱
func (router *stWanba) inquiryRet(w http.ResponseWriter,r *http.Request,myrouter *stRouterPost) (bool,string) {
	openid := r.FormValue("openid")
	//openkey := r.FormValue("openkey")
	//appid := r.FormValue("appid")
	//userip := r.FormValue("userip")
	//count := r.FormValue("count")
	//zoneid := r.FormValue("zoneid")
	sign := r.FormValue("sign")
	stime := r.FormValue("time")
	playerid := r.FormValue("playerid")
	itemid := r.FormValue("itemid")
	//pf := r.FormValue("pf")
	//openid+itemid+time+playerid+key
	md5str := openid + itemid + stime + playerid + Instance().mstJSONConfig.clientkey
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(md5str))
	cipherStr := md5Ctx.Sum(nil)
	mysigon := hex.EncodeToString(cipherStr)
	if(mysigon != sign){
		glog.Errorf("md5 error : me:%s,client:%s",mysigon,sign)
		myrouter.send(w,-100,"md5 error")
		return false,""
	}
	id,_:= strconv.Atoi(itemid)
	value,ok := Instance().getMoney(id)
	if !ok{
		glog.Errorf("没有对应的商品:%s",itemid)
		myrouter.send(w,-100,"no item")
		return false,""
	}
	return true,strconv.Itoa(value)
	//info := router.inquiry(openid,openkey,appid,userip,count,zoneid,pf)
	//if(info != nil){
	//	jsondata, err := gjson.NewJSONByte(info)
	//	if err != nil {
	//		glog.Errorf("inquiryRet NewJsonByte: %s",err)
	//		myrouter.send(w,-100,"inquiryRet json error")
	//		return false,""
	//	}
	//	ret := jsondata.Getint("ret")
	//	if ret == 0{
	//		var score int
	//		keydata := gjson.NewGet(jsondata, "data")
	//		if keydata.IsValid() {
	//			logindata := gjson.NewGetindex(keydata, 0)
	//			if logindata.IsValid() {
	//				score= logindata.Getint("score")
	//			}
	//		}
	//		if(score < value){
	//			myrouter.send(w,-101,"score no have")
	//			return false,""
	//		}
	//		return true,strconv.Itoa(value)
	//
	//	}else{
	//		msg:= jsondata.Getstring("msg")
	//		myrouter.send(w,-100,"inquiry wanba ret:"+ msg)
	//	}
	//}else{
	//	myrouter.send(w,-100,"inquiry error")
	//}
	//return false,""
}
//玩吧查询玩家星币
func (router *stWanba) inquiry(openid,openkey,appid,userip,count,zoneid,pf string) []byte{

	//https://api.urlshare.cn/v3/user/get_playzone_userinfo?
	//openid=B624064BA065E01CB73F835017FE96FA&
	//	zoneid=1&
	//	openkey=5F154D7D2751AEDC8527269006F290F70297B7E54667536C&
	//	appid=2&
	//	sig=VrN%2BTn5J%2Fg4IIo0egUdxq6%2B0otk%3D&
	//	pf=wanba_ts&
	//	format=json&
	//	userip=112.90.139.30
	var buffer bytes.Buffer
	buffer.WriteString("appid=")
	buffer.WriteString(appid)
	buffer.WriteString("&format=json")
	buffer.WriteString("&openid=")
	buffer.WriteString(openid)
	buffer.WriteString("&openkey=")
	buffer.WriteString(openkey)
	buffer.WriteString("&pf=")
	buffer.WriteString(pf)
	buffer.WriteString("&userip=")
	buffer.WriteString(userip)
	buffer.WriteString("&zoneid=")
	buffer.WriteString(zoneid)

	sig:=router.sigCreate("/v3/user/get_playzone_userinfo",buffer.String())
	buffer.WriteString("&sig=")
	buffer.WriteString(sig)

	sendStr := "https://api.urlshare.cn/v3/user/get_playzone_userinfo?"+buffer.String()

	glog.Infoln(sendStr)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := client.Post("https://api.urlshare.cn/v3/user/get_playzone_userinfo","application/x-www-form-urlencoded",&buffer)
	if err == nil {
		defer req.Body.Close()
		body, err := ioutil.ReadAll(req.Body)
		if err == nil {
			glog.Infoln(string(body))
			return body
		}else{
			glog.Errorf("error2:%s",err)
		}
	}else{
		glog.Errorf("error1:%s",err)
	}
	return nil
}

func (router *stWanba)runSendGm(myrouter *stRouterPost,billno,money,itemid,playerid string)  {
	for  {
		time.Sleep(10*time.Second)
		ret,errStr := myrouter.sendgm("/tm_charge/pay/by/wanba",billno,money,itemid,playerid)
		if ret == 0 || ret == 1{
			myrouter.errorSave(billno,money,itemid,playerid,errStr)
			Instance().delBuyMap(playerid)
			go myrouter.testsend("订单补发:"+errStr)
			break
		}else{
			go myrouter.testsend("订单补发失败:"+errStr)
		}
	}
}
//发送gm
func (router *stWanba) sendGm(w http.ResponseWriter,r *http.Request,myrouter *stRouterPost,money,billno string) (bool) {
	playerid := r.FormValue("playerid")
	itemid := r.FormValue("itemid")
	ret,errStr := myrouter.sendgm("/tm_charge/pay/by/wanba",billno,money,itemid,playerid)
	if ret != 0{
		myrouter.errorSave(billno,money,itemid,playerid,errStr)
		if ret == -1{
			go router.runSendGm(myrouter,billno,money,itemid,playerid)
		}else{
			Instance().delBuyMap(playerid)
		}
		go myrouter.testsend("充值错误:"+errStr)
	}else{
		Instance().delBuyMap(playerid)
	}
	return ret == 0
}


//玩吧扣款
func (router *stWanba)deductionRet(w http.ResponseWriter, r *http.Request,myrouter *stRouterPost,money,billno string) bool {

	openid := r.FormValue("openid")
	openkey := r.FormValue("openkey")
	appid := r.FormValue("appid")
	userip := r.FormValue("userip")
	//count := r.FormValue("count")
	zoneid := r.FormValue("zoneid")
	pf := r.FormValue("pf")
	info := router.deduction(openid,openkey,appid,userip,zoneid,money,billno,pf)
	if(info != nil){
		jsondata, err := gjson.NewJSONByte(info)
		if err != nil {
			glog.Errorf("deductionRet NewJsonByte: %s",err)
			myrouter.send(w,-100,"deductionRet json error")
			return false
		}
		code := jsondata.Getint("code")
		if code == 0{
			return true
		}else {
			message:= jsondata.Getstring("message")
			if code == 1004{
				myrouter.send(w,-101,strconv.Itoa(code)+","+message)
			}else if code == 1002{
				if message == "白名单用户额度不够"{
					myrouter.send(w,-101,strconv.Itoa(code)+","+message)
				}else{
					myrouter.send(w,-100,strconv.Itoa(code)+","+message)
				}
			}else{
				myrouter.send(w,-100,strconv.Itoa(code)+","+message)
			}

		}
	}else{
		myrouter.send(w,-100,"deduction error")
	}
	return false
}
//玩吧扣除玩家星币
func (router *stWanba)deduction(openid,openkey,appid,userip,zoneid,money,billno,pf string) []byte {

	//https://api.urlshare.cn/v3/user/buy_playzone_item?
	//	billno=xxxxx&
	//		openid=B624064BA065E01CB73F835017FE96FA&
	//		zoneid=1&
	//		openkey=5F154D7D2751AEDC8527269006F290F70297B7E54667536C&
	//		appid=2&
	//		itemid=10&
	//		count=1&
	//		sig=VrN%2BTn5J%2Fg4IIo0egUdxq6%2B0otk%3D&
	//		pf=wanba_ts&
	//		format=json&
	//		userip=112.90.139.30

	var buffer bytes.Buffer
	buffer.WriteString("appid=")
	buffer.WriteString(appid)
	buffer.WriteString("&billno=")
	buffer.WriteString(billno)
	buffer.WriteString("&count=")
	buffer.WriteString(money)
	buffer.WriteString("&format=json")
	buffer.WriteString("&itemid=")
	if(zoneid == "1"){
		buffer.WriteString("38008")
	}else if(zoneid == "2"){
		buffer.WriteString("38011")
	}
	buffer.WriteString("&openid=")
	buffer.WriteString(openid)
	buffer.WriteString("&openkey=")
	buffer.WriteString(openkey)
	buffer.WriteString("&pf=")
	buffer.WriteString(pf)
	buffer.WriteString("&userip=")
	buffer.WriteString(userip)
	buffer.WriteString("&zoneid=")
	buffer.WriteString(zoneid)

	sig:=router.sigCreate("/v3/user/buy_playzone_item",buffer.String())
	buffer.WriteString("&sig=")
	buffer.WriteString(sig)

	sendStr := "https://api.urlshare.cn/v3/user/buy_playzone_item?" + buffer.String()
	glog.Infoln(sendStr)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := client.Post("https://api.urlshare.cn/v3/user/buy_playzone_item","application/x-www-form-urlencoded",&buffer)
	if err == nil {
		defer req.Body.Close()
		body, err := ioutil.ReadAll(req.Body)
		if err == nil {
			glog.Infoln(string(body))
			return body
		}else{
			glog.Errorf("deduction error2:%s",err)
		}
	}else{
		glog.Errorf("deduction error1:%s",err)
	}
	return nil
}

//玩吧补单
func (router *stWanba)makeUpOrder(w http.ResponseWriter, r *http.Request,back *stBackstage)  {
	playerid := r.FormValue("playerid")
	itemid := r.FormValue("itemid")
	id,_:= strconv.Atoi(itemid)
	value,ok := Instance().getMoney(id)
	if !ok{
		glog.Errorf("没有对应的商品:%s",itemid)
		back.send(w,-100,"no item")
		return
	}
	if(Instance().checkBuyMapKey(playerid)){
		back.send(w,-100,"有订单正在处理")
		return
	}

	Instance().addBuyMap(playerid)

	billno := Instance().routerConnect.createorderId(playerid,"wb")

	ret,errStr := Instance().routerConnect.sendgm("/tm_charge/pay/by/wanba",billno,strconv.Itoa(value),itemid,playerid)
	go Instance().routerConnect.testsend("gm补单充值:"+errStr)
	if ret != 0{
		back.send(w,-100,errStr)
	}else{
		back.send(w,0,errStr)
	}
	Instance().delBuyMap(playerid)
}