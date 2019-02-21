package tool

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

var redisAddr string

type RedisPool struct {
}

func init() {
	redisAddr = "localhost"
}

func NewPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", redisAddr) },
	}
}

func main() {
	conn := NewPool()
	_, _ = conn.Get().Do("set", "key", "12")
}
