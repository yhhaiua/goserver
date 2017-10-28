package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/yhhaiua/goserver/common/glog"
	"github.com/yhhaiua/goserver/gate/logicgate"
)

func main() {

	//runtime.GOMAXPROCS(runtime.NumCPU())
	args := os.Args
	if args == nil || len(args) < 2 {
		return
	}
	svrsplit := strings.Split(args[1], "=")

	if svrsplit == nil || len(svrsplit) != 2 {
		return
	}
	svrid, error := strconv.Atoi(svrsplit[1])

	if error == nil {
		sdir := "./log/gateServer_" + svrsplit[1]
		glog.SetlogDir(sdir)

		logicgate.Instance().LogicInit(svrid)

		select {}
	}
	glog.Flush()
}
