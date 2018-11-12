package main

import (
	"github.com/yhhaiua/goserver/common/grouter"
	"log"
	"net/http"
)

func main() {
	routerInit()
}

func  routerInit() {

	router := grouter.New()

	//Get service information
	router.GET("/login", serverGet)

	//glog.Infof("http监听开启%s", ":19003")
	err := http.ListenAndServe(":19003", router)
	if err != nil {
		//glog.Errorf("http监听s失败 %s", err)
	}

}
func serverGet(w http.ResponseWriter, r *http.Request, _ grouter.Params)  {

	data := r.FormValue("operatorid")
	log.Println(data)
}