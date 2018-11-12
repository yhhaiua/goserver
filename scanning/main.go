package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func main() {

	var clusterinfo = url.Values{}
	clusterinfo.Set("用户名", "chmcqllgameuser")
	clusterinfo.Set("密码", "132f63feb4fc3fa07b15d9a1cd2cf90c")

	urldata := "http://47.98.231.229:81/fight19"

	res, err := http.Get(urldata)
	if err != nil {
		panic(err)

	}
	body, err := ioutil.ReadAll(res.Body)
	if err == nil {
		fmt.Println(string(body))
	}
	// req, err := http.PostForm(urldata, clusterinfo)
	// if err == nil {

	// 	defer req.Body.Close()

	// 	body, err := ioutil.ReadAll(req.Body)
	// 	if err == nil {
	// 		fmt.Println(string(body))
	// 	}
	// }
}
