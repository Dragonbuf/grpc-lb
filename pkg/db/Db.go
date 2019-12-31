package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type Config struct {
	Dns          string
	LogMode      bool
	MaxIdleConns int
	OpenConns    int
}

func NewMysql(c *Config) *gorm.DB {

	Mysql, err := gorm.Open("mysql", c.Dns)
	if err != nil {
		panic(err)
	}
	Mysql.LogMode(c.LogMode)
	Mysql.DB().SetMaxIdleConns(c.MaxIdleConns)
	Mysql.DB().SetMaxOpenConns(c.OpenConns)
	return Mysql
}
