package common

import (
	"database/sql"

	"github.com/yhhaiua/goserver/common/glog"
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
	db             *sql.DB
	bodbconnection bool
}

//NewMysql mysql连接创建
func NewMysql(Config *MysqlConfig) (mydb *MysqlDB, err error) {
	mydb = newMysql()
	sconmysql := Config.Suser + ":" + Config.Spassword + "@tcp(" + Config.Shost + ")/" + Config.Sdbname + "?charset=utf8mb4"
	mydb.db, err = sql.Open("mysql", sconmysql)
	if err == nil {

		mydb.db.SetMaxOpenConns(Config.Maxopen)
		mydb.db.SetMaxIdleConns(Config.Maxidle)
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
