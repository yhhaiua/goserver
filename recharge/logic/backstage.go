package logic

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/yhhaiua/goserver/common/glog"
	"github.com/yhhaiua/goserver/common/grouter"
	"net/http"
)

type stBackstage struct {

}
//停止充值
func (back *stBackstage) stopCharge(w http.ResponseWriter, r *http.Request, _ grouter.Params)  {

	operatorid := r.FormValue("operatorid")
	sign := r.FormValue("sign")
	stime := r.FormValue("time")
	stype := r.FormValue("type")

	if operatorid != Instance().mstJSONConfig.operatorid{
		glog.Errorf("stopCharge operatorid error : me:%s,client:%s,操作人:%s",Instance().mstJSONConfig.operatorid,operatorid,Instance().getUserIp(r))
		back.send(w,-100,"operatorid error")
		return
	}
	//operatorid+stime+stype+key
	md5str := operatorid + stime +stype+Instance().mstJSONConfig.rechargekey
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(md5str))
	cipherStr := md5Ctx.Sum(nil)
	mysigon := hex.EncodeToString(cipherStr)
	if(mysigon != sign){
		glog.Errorf("stopCharge md5 error : me:%s,client:%s,操作人:%s",mysigon,sign,Instance().getUserIp(r))
		back.send(w,-100,"md5 error")
		return
	}
	if stype == "1"{
		Instance().stopRecharge = true
		glog.Infof("stopCharge stop success : 操作人:%s",Instance().getUserIp(r))
		back.send(w,0,"stop success")
	}else if stype == "0"{
		Instance().stopRecharge = false
		glog.Infof("stopCharge open success : 操作人:%s",Instance().getUserIp(r))
		back.send(w,0,"open success")
	}else{
		glog.Errorf("stopCharge stype error :client:%s,操作人:%s",stype,Instance().getUserIp(r))
		back.send(w,-100,"stype error")
	}
}

//补单
func (back *stBackstage)makeUpOrder(w http.ResponseWriter, r *http.Request, _ grouter.Params)  {
	operatorid := r.FormValue("operatorid")
	sign := r.FormValue("sign")
	stime := r.FormValue("time")
	itemid := r.FormValue("itemid")
	playerid := r.FormValue("playerid")

	if operatorid != Instance().mstJSONConfig.operatorid{
		glog.Errorf("makeUpOrder operatorid error : me:%s,client:%s,操作人:%s",Instance().mstJSONConfig.operatorid,operatorid,Instance().getUserIp(r))
		back.send(w,-100,"operatorid error")
		return
	}
	//operatorid+stime+itemid+playerid+key
	md5str := operatorid + stime +itemid+playerid+Instance().mstJSONConfig.rechargekey
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(md5str))
	cipherStr := md5Ctx.Sum(nil)
	mysigon := hex.EncodeToString(cipherStr)
	if(mysigon != sign){
		glog.Errorf("makeUpOrder md5 error : me:%s,client:%s,操作人:%s",mysigon,sign,Instance().getUserIp(r))
		back.send(w,-100,"md5 error")
		return
	}
	glog.Infof("makeUpOrder success : 操作人:%s",Instance().getUserIp(r))
	switch operatorid {
	case "1":
		//玩吧渠道
		Instance().routerConnect.wanbaDeal.makeUpOrder(w,r,back)
	default:
		glog.Errorf("错误渠道请求:%s",operatorid)
		back.send(w,-100," operatorid no have error")
	}
}


//发送返回
func (back *stBackstage)send(w http.ResponseWriter,ret int,msg string)  {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	var info RetInfo
	info.Ret = ret
	info.Msg = msg
	Message, err := json.Marshal(info)

	if err == nil {
		fmt.Fprintf(w, "%s", Message)
	}
}