package cache

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

type Config struct {
	MaxIdle      int
	MaxActive    int
	IdleTimeout  time.Duration
	Wait         bool
	Host         string
	Password     string
	Db           int
	ConTimeout   time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func NewRedisPool(c *Config) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     c.MaxIdle,
		MaxActive:   c.MaxActive,
		IdleTimeout: c.IdleTimeout * time.Second,
		Wait:        c.Wait,
		Dial: func() (redis.Conn, error) {
			con, err := redis.Dial("tcp", c.Host,
				redis.DialPassword(c.Password),
				redis.DialDatabase(c.Db),
				redis.DialConnectTimeout(c.ConTimeout*time.Second),
				redis.DialReadTimeout(c.ReadTimeout*time.Second),
				redis.DialWriteTimeout(c.WriteTimeout*time.Second))
			if err != nil {
				panic(err)
			}
			return con, nil
		},
	}
}
