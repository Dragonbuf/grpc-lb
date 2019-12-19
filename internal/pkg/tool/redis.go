package tool

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"grpc-lb/configs"
	"time"
)

var RedisPool *redis.Pool

func init() {
	fmt.Println("redis pool init")
	RedisPool = &redis.Pool{
		MaxIdle:     configs.RedisMaxIdle,
		MaxActive:   configs.RedisMaxActive,
		IdleTimeout: configs.RedisMaxIdleTimeout * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			con, err := redis.Dial("tcp", configs.RedisHost,
				redis.DialPassword(configs.RedisPassword),
				redis.DialDatabase(configs.RedisDb),
				redis.DialConnectTimeout(configs.RedisConTimeout),
				redis.DialReadTimeout(configs.RedisReadTimeout),
				redis.DialWriteTimeout(configs.RedisWriteTimeout))
			if err != nil {
				panic(err)
				return nil, err
			}
			return con, nil
		},
	}
}
