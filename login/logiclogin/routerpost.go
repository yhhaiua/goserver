package logiclogin

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/yhhaiua/goserver/common/glog"

	"github.com/yhhaiua/goserver/common/grouter"
)

type stRetMessage struct {
	Code int
	URL  string
	Msg  string
}

type stZoneMessage struct {
	ZoneData []string
}

type stAccountMessage struct {
	Account string
	Onlyid  string
	Name    string
	Time    string
	Sign    string
}

const (
	routerRetSucc    = 0  //成功
	routerRetFail    = 1  //数据读取失败
	routerRetNoPname = 2  //没有对应平台版本
	routerRetVersion = 3  //版本错误更新新版本有URL内容
	routerRetNoJSON  = 4  //服务端解析失败
	routerMd5Error   = 5  //MD5错误
	routerRepeatAcc  = 6  //重复创建账号
	routerRetNoZone  = 7  //没有区服数据
	routerRetNoOnly  = 8  //OnlyId获取失败
	routerRetSave    = 9  //保存数据失败
	routerRetLoad    = 10 //正在加载数据
)

const (
	versionData    = "VersionData_"   //版本表
	onlyidData     = "GlobalOnlyid_0" //玩家唯一id
	gameData       = "game_"          //区服表
	cAccountData   = "AccountData"    //帐号数据
	writetablename = "updatetable"    //写入改变列表
)

const (
	md5Key = "yhhaiua" //MD5
)

//stRedisVersion redis version data key:pname,value:stRedisVersion
type stRedisVersion struct {
	Pid     int
	Version string
	URL     string
}
type stAccountData struct {
	Onlyid string
	Name   string
	Zoneid string
}

type stRouterPost struct {
}

//getMyOnlyid 获取玩家唯一id
func (myrouter *stRouterPost) getMyOnlyid() int64 {
	var MyOnlyid int64
	if Instance().redisdb() != nil {
		MyOnlyid = Instance().redisdb().Incr(onlyidData)
		if MyOnlyid > 0 {
			err := Instance().redisdb().Zadd(writetablename, onlyidData)
			if err != nil {
				glog.Errorf("Zadd 压人失败 %s,%s", writetablename, onlyidData)
			}
		}
	}
	return MyOnlyid
}

//getPlatformZone 获取游戏区服 http://localhost/public/getPlatformZone?pname=7cool&version=v1.5.6
func (myrouter *stRouterPost) getPlatformZone(w http.ResponseWriter, r *http.Request, _ grouter.Params) {

	var MessageCode stRetMessage
	MessageCode.getMessage(r)

	if MessageCode.Code == routerRetSucc {
		if Instance().redisdb() != nil {
			pname := r.FormValue("pname")
			key := gameData + pname
			panmedata, err := Instance().redisdb().Keys(key)

			if err == nil {
				var RedisData stZoneMessage
				for i := 0; i < len(panmedata); i++ {
					value, err := Instance().redisdb().Get(panmedata[i])
					if err == nil {
						RedisData.ZoneData = append(RedisData.ZoneData, value)
					}
				}
				Message, err := json.Marshal(RedisData)
				if err == nil {
					fmt.Fprintf(w, "%s", Message)
				}

			} else {
				MessageCode.Code = routerRetNoZone
			}
		}
	}
	if MessageCode.Code != routerRetSucc {
		Message, err := json.Marshal(MessageCode)

		if err == nil {
			fmt.Fprintf(w, "%s", Message)
		}
	}
}

//regAccount 获取玩家帐号信息 http://172.16.2.36:2044/public/regAccount?pname=7cool&version=v1.5.6&time=636401136976819662&sign=95c30fa566a0d19248c31cf1d449b20e&account=yhhhh1&serverid=1
func (myrouter *stRouterPost) regAccount(w http.ResponseWriter, r *http.Request, _ grouter.Params) {

	var MessageCode stRetMessage
	MessageCode.getMessage(r)

	if MessageCode.Code == routerRetSucc {
		if Instance().redisdb() != nil {
			pname := r.FormValue("pname")
			Mtime := r.FormValue("time")
			sign := r.FormValue("sign")
			account := r.FormValue("account")
			serverid := r.FormValue("serverid")

			md5str := pname + Mtime + account + serverid + md5Key
			md5Ctx := md5.New()
			md5Ctx.Write([]byte(md5str))
			cipherStr := md5Ctx.Sum(nil)
			mysigon := hex.EncodeToString(cipherStr)

			if mysigon == sign {
				playeraccount := pname + account + serverid
				TempDataName := cAccountData + pname + serverid + "_" + playeraccount
				value, err := Instance().redisdb().Get(TempDataName)

				var RedisData stAccountData
				if err != nil {
					gamekey := gameData + pname + serverid
					if Instance().redisdb().IsExist(gamekey) {
						check := Instance().redisdb().Incr(playeraccount)
						if check == 1 {
							getOnlyid := myrouter.getMyOnlyid()
							if getOnlyid > 0 {
								RedisData.Onlyid = strconv.FormatInt(getOnlyid, 10)
								RedisData.Name = RedisData.Onlyid
								RedisData.Zoneid = serverid
								RedisMessage, err := json.Marshal(RedisData)
								if err == nil {

									glog.Infof("保存玩家数据%s,%s", TempDataName, RedisMessage)

									err = Instance().redisdb().Set(TempDataName, RedisMessage)
									if err == nil {
										err = Instance().redisdb().Zadd(writetablename, TempDataName)
										if err != nil {
											glog.Errorf("Zadd 压人失败 %s,%s", writetablename, TempDataName)
										}
									} else {
										MessageCode.Code = routerRetSave
									}
								} else {
									MessageCode.Code = routerRetNoJSON
								}
							} else {
								MessageCode.Code = routerRetNoOnly
							}

						} else {
							MessageCode.Code = routerRepeatAcc
						}
					} else {
						MessageCode.Code = routerRetNoZone
					}
				} else {
					if err = json.Unmarshal([]byte(value), &RedisData); err != nil {
						MessageCode.Code = routerRetNoJSON
					}
				}

				if MessageCode.Code == routerRetSucc {

					var MessageStruct stAccountMessage
					MessageStruct.Account = playeraccount
					MessageStruct.Name = RedisData.Name
					MessageStruct.Onlyid = RedisData.Onlyid
					MessageStruct.Time = strconv.FormatInt(time.Now().Unix(), 10)
					md5getstr := MessageStruct.Account + MessageStruct.Onlyid + MessageStruct.Name + MessageStruct.Time + md5Key
					md5Send := md5.New()
					md5Send.Write([]byte(md5getstr))
					md5SendCtr := md5Send.Sum(nil)
					MessageStruct.Sign = hex.EncodeToString(md5SendCtr)
					Message, err := json.Marshal(MessageStruct)
					if err == nil {
						fmt.Fprintf(w, "%s", Message)
					}
				}

			} else {
				MessageCode.Code = routerMd5Error
			}
		}
	}
	if MessageCode.Code != routerRetSucc {
		Message, err := json.Marshal(MessageCode)

		if err == nil {
			fmt.Fprintf(w, "%s", Message)
		}
	}
}
func (MessageCode *stRetMessage) getMessage(r *http.Request) {
	pname := r.FormValue("pname")
	version := r.FormValue("version")
	MessageCode.Code = routerRetSucc

	if Instance().redisdb() != nil {
		if Instance().redismsg().boCon() {
			TempDataName := versionData + pname
			value, err := Instance().redisdb().Get(TempDataName)

			if err == nil {
				var RedisData stRedisVersion
				if err = json.Unmarshal([]byte(value), &RedisData); err == nil {
					if RedisData.Version == version {

					} else {
						MessageCode.Code = routerRetVersion
						MessageCode.URL = RedisData.URL
					}
				} else {
					MessageCode.Code = routerRetNoJSON
				}

			} else {
				MessageCode.Code = routerRetNoPname
			}
		} else {
			MessageCode.Code = routerRetLoad
		}

	} else {
		MessageCode.Code = routerRetFail
	}
}
