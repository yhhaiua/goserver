package main

import (
	"github.com/yhhaiua/goserver/common/glog"
	"github.com/yhhaiua/goserver/timemachine/logicmachine"
)

func main() {

	sdir := "./logs/timed"
	glog.SetlogDir(sdir)
	if logicmachine.Instance().LogicInit() {
		select {}
	}
	glog.Flush()

}
