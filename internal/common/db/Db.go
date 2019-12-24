package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"grpc-lb/configs"
)

var Mysql *gorm.DB

func init() {
	cfg := configs.GetConfig()
	var err error
	Mysql, err = gorm.Open("mysql", cfg.MustValue("mysql", "MysqlMasterDns"))
	if err != nil {
		panic(err)
	}
	Mysql.LogMode(cfg.MustBool("mysql", "LogMode"))
	Mysql.DB().SetMaxIdleConns(cfg.MustInt("mysql", "MysqlMaxIdleConns"))
	Mysql.DB().SetMaxOpenConns(cfg.MustInt("mysql", "MysqlMaxOpenConns"))
}
