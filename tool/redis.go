package tool

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"grpc-lb/cmd/config"
	"time"
)

var RedisPool *redis.Pool

func init() {
	fmt.Println("redis pool init")
	RedisPool = &redis.Pool{
		MaxIdle:     config.RedisMaxIdle,
		MaxActive:   config.RedisMaxActive,
		IdleTimeout: config.RedisMaxIdleTimeout * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			con, err := redis.Dial("tcp", config.RedisHost,
				redis.DialPassword(config.RedisPassword),
				redis.DialDatabase(config.RedisDb),
				redis.DialConnectTimeout(config.RedisConTimeout),
				redis.DialReadTimeout(config.RedisReadTimeout),
				redis.DialWriteTimeout(config.RedisWriteTimeout))
			if err != nil {
				panic(err)
				return nil, err
			}
			return con, nil
		},
	}
}

func test() {
	pool := RedisPool
	conn := pool.Get()
	defer pool.Close()
	if conn.Err() != nil {
		//TODO
	}

}
