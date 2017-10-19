package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func main() {

	getPlatformZone()

	for {
		fmt.Println("Please input your full name: ")
		var data string
		fmt.Scanln(&data)
		if len(data) > 0 {
			regAccount(data)
		}
	}

}

func getPlatformZone() {

	var clusterinfo = url.Values{}
	clusterinfo.Set("pname", "7cool")
	clusterinfo.Set("version", "1.0.0")

	urldata := "http://172.16.3.73:8001/public/getPlatformZone"

	req, err := http.PostForm(urldata, clusterinfo)
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

	urldata := "http://172.16.3.73:8001/public/regAccount"

	req, err := http.PostForm(urldata, clusterinfo)
	if err == nil {

		defer req.Body.Close()

		body, err := ioutil.ReadAll(req.Body)
		if err == nil {
			fmt.Println(string(body))
		}
	}
}
