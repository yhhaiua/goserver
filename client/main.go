package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/yhhaiua/goserver/common/glog"
	"golang.org/x/net/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/axgle/mahonia"
)

func main() {

	//	getPlatformZone()

	//	for {
	//		fmt.Println("Please input your full name: ")
	//		var data string
	//		fmt.Scanln(&data)
	//		if len(data) > 0 {
	//			regAccount(data)
	//		}
	//	}
	//getPlatformZone()
	//get()
	//websend()
	//cestext()
	//stopCharge()
	makeuporder()
}
func websend()  {
	var origin = "http://127.0.0.1:18003/"
	var url = "ws://127.0.0.1:18003/echo"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	message := []byte("hello, world!你好")
	_, err = ws.Write(message)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Send: %s\n", message)

	var msg = make([]byte, 512)
	m, err := ws.Read(msg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Receive: %s\n", msg[:m])

	ws.Close()//关闭连接
}
func Echo(ws *websocket.Conn) {

	var err error

	for {

		var reply string

		//websocket接受信息

		if err = websocket.Message.Receive(ws, &reply); err != nil {

			fmt.Println("receive failed:", err)

			break

		}

		fmt.Println("reveived from client: " + reply)

		msg := "received:" + reply

		fmt.Println("send to client:" + msg)

		//这里是发送消息

		if err = websocket.Message.Send(ws, msg); err != nil {

			fmt.Println("send failed:", err)

			break

		}

	}

}

func get() {
	response, _ := http.Get("https://tcc.taobao.com/cc/json/mobile_tel_segment.htm?tel=15850781443")
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println(ConvertToString(string(body), "gbk", "utf-8"))
}

func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

func getPlatformZone() {

	//var clusterinfo = url.Values{}
	//clusterinfo.Set("operatorid", "1")
	//clusterinfo.Set("version", "1.0.0")

	urldata := "http://127.0.0.1:19003/login?operatorid=1&version=2"

	req, err := http.PostForm(urldata, nil)
	if err == nil {

		defer req.Body.Close()

		body, err := ioutil.ReadAll(req.Body)
		if err == nil {
			fmt.Println(string(body))
		}
	}
}

func regAccount(account string) {
	var clusterinfo = url.Values{}
	clusterinfo.Set("pname", "7cool")
	clusterinfo.Set("version", "1.0.0")

	stime := strconv.FormatInt(time.Now().Unix(), 10)

	md5str := "7cool" + stime + account + "1" + "yhhaiua"
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(md5str))
	cipherStr := md5Ctx.Sum(nil)
	mysigon := hex.EncodeToString(cipherStr)

	clusterinfo.Set("time", stime)
	clusterinfo.Set("sign", mysigon)
	clusterinfo.Set("account", account)
	clusterinfo.Set("serverid", "1")

	urldata := "http://127.0.0.1:19003/public/regAccount"

	req, err := http.PostForm(urldata, clusterinfo)
	if err == nil {

		defer req.Body.Close()

		body, err := ioutil.ReadAll(req.Body)
		if err == nil {
			fmt.Println(string(body))
		}
	}
}

func cestext()  {
	timeout := time.Duration(3 * time.Second)//超时时间3s
	client := &http.Client{
		Timeout: timeout,
	}
	glog.Infof("start")
	req, err := client.PostForm("https://font.tmcb.jiulingwan.com/login?ip=shenhe.tmcb.jiulingwan.com&port=19003&operatorid=1&serverid=1&time=1541682796&sign=0c681858d99e9caaaaed548a1746ef0a&account=E24DC90D59ED5D61451863BCD5ED67E0",nil)
	if err == nil {
		defer req.Body.Close()
		body, err := ioutil.ReadAll(req.Body)
		if err == nil {
			valueRet := string(body)
			glog.Infof(valueRet)
		}else{
			glog.Errorf("sendgm error2:%s",err)
		}
	}else{
		glog.Errorf("sendgm error1:%s",err)
	}
}

func stopCharge()  {
	var buffer bytes.Buffer
	buffer.WriteString("http://127.0.0.1:8001/stopcharge?")
	buffer.WriteString("operatorid=")
	buffer.WriteString("1")
	buffer.WriteString("&time=")
	timestamp := time.Now().Unix()
	stime := strconv.FormatInt(timestamp, 10)
	buffer.WriteString(stime)
	buffer.WriteString("&type=")
	buffer.WriteString("1")

	//operatorid+stime+stype+key
	md5str := "1" + stime +"1"+"chm&tianmi&recharge&1122"
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(md5str))
	cipherStr := md5Ctx.Sum(nil)
	mysigon := hex.EncodeToString(cipherStr)
	buffer.WriteString("&sign=")
	buffer.WriteString(mysigon)
	fmt.Println(buffer.String())
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := client.Get(buffer.String())
	if err == nil {
		defer req.Body.Close()
		body, err := ioutil.ReadAll(req.Body)
		if err == nil {
			valueRet := string(body)
			fmt.Println(valueRet)
		}else {
			fmt.Println(err)
		}
	}else {
		fmt.Println(err)
	}
}
func makeuporder()  {

	var buffer bytes.Buffer
	buffer.WriteString("http://127.0.0.1:8001/makeuporder?")
	buffer.WriteString("operatorid=")
	buffer.WriteString("1")
	buffer.WriteString("&time=")
	timestamp := time.Now().Unix()
	stime := strconv.FormatInt(timestamp, 10)
	buffer.WriteString(stime)
	buffer.WriteString("&itemid=11")
	buffer.WriteString("&playerid=8000101")
	//operatorid+stime+itemid+playerid+key
	md5str := "1" + stime +"118000101"+"chm&tianmi&recharge&1122"
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(md5str))
	cipherStr := md5Ctx.Sum(nil)
	mysigon := hex.EncodeToString(cipherStr)
	buffer.WriteString("&sign=")
	buffer.WriteString(mysigon)
	fmt.Println(buffer.String())
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := client.Get(buffer.String())
	if err == nil {
		defer req.Body.Close()
		body, err := ioutil.ReadAll(req.Body)
		if err == nil {
			valueRet := string(body)
			fmt.Println(valueRet)
		}else {
			fmt.Println(err)
		}
	}else {
		fmt.Println(err)
	}
}
