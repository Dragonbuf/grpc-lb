package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"grpc-lb/configs"
)

var Mysql *gorm.DB

func init() {
	var err error
	Mysql, err = gorm.Open("mysql", configs.MysqlMasterDns)
	if err != nil {
		panic(err)
	}
	Mysql.LogMode(true)
	Mysql.DB().SetMaxIdleConns(configs.MysqlMaxIdleConns)
	Mysql.DB().SetMaxOpenConns(configs.MysqlMaxOpenConns)
}
