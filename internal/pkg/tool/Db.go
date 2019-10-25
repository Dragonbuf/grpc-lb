package tool

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"grpc-lb/config"
)

var Mysql *gorm.DB

func init() {
	var err error
	Mysql, err = gorm.Open("mysql", config.MysqlMasterDns)
	if err != nil {
		panic(err)
	}
	Mysql.DB().SetMaxIdleConns(config.MysqlMaxIdleConns)
	Mysql.DB().SetMaxOpenConns(config.MysqlMaxOpenConns)
}
