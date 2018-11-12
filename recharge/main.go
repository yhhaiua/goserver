package main

import (
	"github.com/yhhaiua/goserver/common/glog"
	"github.com/yhhaiua/goserver/recharge/logic"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	sdir := "./logs/"
	glog.SetlogDir(sdir)

	if logic.Instance().LogicInit(){
		glog.Infof("recharge 启动成功")
	}else{
		glog.Infof("recharge 启动失败")
	}
	glog.Flush()
}
