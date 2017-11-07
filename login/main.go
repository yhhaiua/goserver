package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/yhhaiua/goserver/common/glog"
	"github.com/yhhaiua/goserver/login/logiclogin"
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
	svrid, err := strconv.Atoi(svrsplit[1])

	if err == nil {
		sdir := "./log/loginServer_" + svrsplit[1]
		glog.SetlogDir(sdir)

		logiclogin.Instance().LogicInit(svrid)
	}
	glog.Flush()
}
