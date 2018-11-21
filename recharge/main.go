package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/yhhaiua/goserver/common/log4go"
	"github.com/yhhaiua/goserver/recharge/logic"
)

func main() {

	log4go.LoadConfiguration("config/log4j.xml")

	if logic.Instance().LogicInit(){
		log4go.Info("recharge 启动成功")
	}else{
		log4go.Error("recharge 启动失败")
	}
}
