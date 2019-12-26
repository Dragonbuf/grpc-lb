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
	cfg := configs.GetConfig()

	RedisPool = &redis.Pool{
		MaxIdle:     cfg.MustInt("redis", "RedisMaxIdle"),
		MaxActive:   cfg.MustInt("redis", "RedisMaxActive"),
		IdleTimeout: time.Duration(cfg.MustInt("redis", "RedisIdleTimeout")) * time.Second,
		Wait:        cfg.MustBool("redis", "Wait"),
		Dial: func() (redis.Conn, error) {
			con, err := redis.Dial("tcp", cfg.MustValue("redis", "RedisHost"),
				redis.DialPassword(cfg.MustValue("redis", "RedisPassword")),
				redis.DialDatabase(cfg.MustInt("redis", "RedisDb")),
				redis.DialConnectTimeout(time.Duration(cfg.MustInt("redis", "RedisConTimeout"))*time.Second),
				redis.DialReadTimeout(time.Duration(cfg.MustInt("redis", "RedisReadTimeout"))*time.Second),
				redis.DialWriteTimeout(time.Duration(cfg.MustInt("redis", "RedisWriteTimeout"))*time.Second))
			if err != nil {
				panic(err)
				return nil, err
			}
			return con, nil
		},
	}
}
