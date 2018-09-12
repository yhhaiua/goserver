package main

import (
	"github.com/yhhaiua/goserver/common/glog"
	"github.com/yhhaiua/goserver/timedsave/logictime"
)

func main() {

	sdir := "./logs/timed"
	glog.SetlogDir(sdir)
	if logictime.Instance().LogicInit() {
		select {}
	}
	glog.Flush()
}
