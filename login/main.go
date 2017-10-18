package main

import (
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/yhhaiua/goserver/common/glog"
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
		sdir := "./log/loginServer_" + svrsplit[1]
		glog.SetlogDir(sdir)

		logiclogin.Instance().LogicInit(svrid)
	}
	glog.Flush()
}
