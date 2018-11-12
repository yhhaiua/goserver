package gmysql

import (
	"database/sql"
)

//MysqlConfig 连接配置
type MysqlConfig struct {
	Shost     string //ipport
	Sdbname   string //数据库名
	Suser     string //用户名
	Spassword string //密码
	Maxopen   int    //最大连接数
	Maxidle   int    //最大空闲数
}

//MysqlDB Mysql连接结构
type MysqlDB struct {
	*sql.DB
	bodbconnection bool
}

//NewMysql mysql连接创建
func NewMysql(Config *MysqlConfig) (mydb *MysqlDB, err error) {
	mydb = newMysql()
	sconmysql := Config.Suser + ":" + Config.Spassword + "@tcp(" + Config.Shost + ")/" + Config.Sdbname + "?charset=utf8mb4"
	mydb.DB, err = sql.Open("mysql", sconmysql)
	if err == nil {

		mydb.SetMaxOpenConns(Config.Maxopen)
		mydb.SetMaxIdleConns(Config.Maxidle)
	} else {
		mydb = nil
	}
	return
}

//CheckPing 检测连接
func (mydb *MysqlDB) CheckPing() error {
	err := mydb.Ping()
	if err == nil {
		mydb.bodbconnection = true
	}
	return err
}

func newMysql() *MysqlDB {
	newMysql := new(MysqlDB)

	return newMysql
}

//HaveConnect 判断是否已经连接
func (mydb *MysqlDB) HaveConnect() bool {
	if mydb.DB != nil && mydb.bodbconnection {
		return true
	}
	return false
}

//Create 创建mysql通用表
func (mydb *MysqlDB) Create(dbname string) error {

	statement := "CREATE TABLE IF NOT EXISTS `" + dbname + "` (`key` varchar(48) NOT NULL,`value` longtext,PRIMARY KEY (`key`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;"

	_, err := mydb.Exec(statement)

	return err
}

////SavetoRedis mysql转存数据到redis
//func (mydb *MysqlDB) SavetoRedis(redisConnect *gredis.RedisPacket, tablename string, endSync *sync.WaitGroup) {
//
//	go mydb.mysqltoredis(redisConnect, tablename, endSync)
//}
//
//func (mydb *MysqlDB) mysqltoredis(redisConnect *gredis.RedisPacket, tablename string, endSync *sync.WaitGroup) {
//
//	sQuery := "SELECT * FROM " + tablename
//	rows, err := mydb.Query(sQuery)
//
//	if err == nil {
//		defer rows.Close()
//		var Key, Value string
//
//		for rows.Next() {
//
//			err = rows.Scan(&Key, &Value)
//			if err != nil {
//				glog.Errorf("mysql读取错误 %s 2 %s", tablename, err)
//			} else {
//
//				//保存到redis
//				if redisConnect != nil {
//					redisKey := tablename + "_" + Key
//					redisConnect.Set(redisKey, Value)
//					glog.Infof("redis添加数据 %s, %s", redisKey, Value)
//				}
//			}
//		}
//
//	} else {
//
//		glog.Errorf("mysql读取错误 %s 1 %s", tablename, err)
//	}
//	if endSync != nil {
//		endSync.Done()
//	}
//}

//Update mysql更新通用数据
func (mydb *MysqlDB) Update(skey, svale, tablename string) (err error) {

	updatement := "update " + tablename + " set `value`=? where `key`=?"
	stmt, err := mydb.Prepare(updatement)
	if err == nil {
		defer stmt.Close()
		ret, err := stmt.Exec(svale, skey)
		if err == nil {
			if hs, _ := ret.RowsAffected(); hs == 0 {
				insertment := "INSERT " + tablename + " set `key`=?,`value`=?"
				stmtinstet, err := mydb.Prepare(insertment)
				if err == nil {
					defer stmtinstet.Close()
					ret, err = stmtinstet.Exec(skey, svale)
				}
			}

		}
	}
	return
}
