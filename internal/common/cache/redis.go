package cache

import (
	"github.com/gomodule/redigo/redis"
	"grpc-lb/configs"
	"time"
)

var RedisPool *redis.Pool

type Redis struct {
	pool *redis.Pool
}

func init() {
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

func NewRedis() *Redis {
	return &Redis{pool: RedisPool}
}
