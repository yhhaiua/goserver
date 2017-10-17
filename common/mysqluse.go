package common

import (
	"database/sql"

	"github.com/yhhaiua/goserver/common/glog"
)

//MysqlDB Mysql连接结构
type MysqlDB struct {
	db             *sql.DB
	bodbconnection bool
}

//NewMysql mysql连接创建
func NewMysql(suser, spassword, shost, sdbname string, maxopen int, maxidle int) (mydb *MysqlDB, err error) {
	mydb = newMysql()
	sconmysql := suser + ":" + spassword + "@tcp(" + shost + ")/" + sdbname + "?charset=utf8mb4"
	mydb.db, err = sql.Open("mysql", sconmysql)
	if err == nil {

		mydb.db.SetMaxOpenConns(maxopen)
		mydb.db.SetMaxIdleConns(maxidle)
	}
	if err != nil {
		mydb = nil
	}
	return
}

//CheckPing 检测连接
func (mydb *MysqlDB) CheckPing() error {
	err := mydb.db.Ping()
	if err == nil {
		mydb.bodbconnection = true
	}
	return err
}

func newMysql() *MysqlDB {
	newMysql := new(MysqlDB)

	return newMysql
}

//Create 创建mysql通用表
func (mydb *MysqlDB) Create(dbname string) error {

	statement := "CREATE TABLE IF NOT EXISTS `" + dbname + "` (`key` varchar(48) NOT NULL,`value` longtext,PRIMARY KEY (`key`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;"

	_, err := mydb.db.Exec(statement)

	return err
}

//SavetoRedis mysql转存数据到redis
func (mydb *MysqlDB) SavetoRedis(redisConnect *RedisPool, tablename string) {
	sQuery := "SELECT * FROM " + tablename
	rows, err := mydb.db.Query(sQuery)

	if err == nil {
		defer rows.Close()
		var Key, Value string

		for rows.Next() {

			err = rows.Scan(&Key, &Value)
			if err != nil {
				glog.Errorf("mysql读取错误 %s 2 %s", tablename, err)
			} else {

				//保存到redis
				if redisConnect != nil {
					redisKey := tablename + Key
					redisConnect.Set(redisKey, Value)
					glog.Infof("redis添加玩家帐号数据 %s, %s", redisKey, Value)
				}
			}
		}

	} else {

		glog.Errorf("mysql读取错误 %s 1 %s", tablename, err)
	}
}
