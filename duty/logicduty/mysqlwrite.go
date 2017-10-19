package logicduty

import (
	"strings"
	"time"

	"github.com/yhhaiua/goserver/common/glog"
)

const (
	writetickertime = 60 * time.Second
	writetablename  = "updatetable"
)

type stMysqlWrite struct {
}

func (mydata *stMysqlWrite) Write() {

	if Instance().mysqldb().HaveConnect() {
		go mydata.datewrite()
	}

}

func (mydata *stMysqlWrite) datewrite() {

	//ZRANGE cheset 0 -1 WITHSCORES
	for _ = range time.NewTicker(writetickertime).C {
		if Instance().redisdb() != nil {
			sdata, err := Instance().redisdb().Zrange(writetablename)
			if err == nil {
				datelen := len(sdata)
				if datelen%2 == 0 {
					for i := 0; i < datelen; i = i + 2 {
						mydata.dateUpdateMysql(sdata[i], sdata[i+1])
					}
				}
			} else {
				glog.Errorf("datewrite table:%s,err:%s", writetickertime, err)
			}

		}
	}
}

func (mydata *stMysqlWrite) dateUpdateMysql(table string, mytime string) {
	var skey, svalue, stable string
	svalue, err := Instance().redisdb().Get(table)
	if err == nil {
		sdata := strings.Split(table, "_")
		if len(sdata) == 2 {
			stable = sdata[0]
			skey = sdata[1]
		} else {
			glog.Errorf("dateUpdateMysql redis获取表名长度错误 表:[%s]", table)
			return
		}
	} else {
		glog.Errorf("dateUpdateMysql redis获取数值错误 表:[%s],错误:[%s]", table, err)
		return
	}

	err = Instance().mysqldb().Update(skey, svalue, stable)
	if err == nil {
		glog.Infof("Update数据成功 %s,%s,%s", skey, svalue, stable)
		snewtime, err := Instance().redisdb().Zscore(writetablename, table)
		if err == nil {
			if snewtime == mytime {
				err = Instance().redisdb().Zrem(writetablename, table)
				if err != nil {
					glog.Errorf("dateUpdateMysql redis Zrem错误 表:[%s],错误:[%s]", table, err)
				}
			} else {
				glog.Infof("dateUpdateMysql 表有新的更新，等待下次循环[%s]", table)
			}
		} else {
			glog.Errorf("dateUpdateMysql redis zscore错误 表:[%s],错误:[%s]", table, err)
		}
	} else {
		glog.Errorf("dateUpdateMysql myslq更新错误 表:[%s],错误:[%s]", stable, err)
	}
}
