package tool

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

var RedisPool *redis.Pool

/**
MaxActive 最大连接数，即最多的tcp连接数，一般建议往大的配置，但不要超过操作系统文件句柄个数（centos下可以ulimit -n查看）。
MaxIdle 最大空闲连接数，即会有这么多个连接提前等待着，但过了超时时间也会关闭。
IdleTimeout 空闲连接超时时间，但应该设置比redis服务器超时时间短。否则服务端超时了，客户端保持着连接也没用。
Wait 这是个很有用的配置:如果超过最大连接，是报错，还是等待。
*/

var (
	host           = "localhost:6379"
	password       = ""
	db             = 0
	conTimeout     = 5 * time.Second
	readTimeout    = 2 * time.Second
	writeTimeout   = 1 * time.Second
	maxIdle        = 1000
	maxActive      = 3000
	MaxIdleTimeout = 5 * time.Second
)

func init() {
	fmt.Println("redis pool init")
	RedisPool = &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: MaxIdleTimeout * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			con, err := redis.Dial("tcp", host,
				redis.DialPassword(password),
				redis.DialDatabase(db),
				redis.DialConnectTimeout(conTimeout),
				redis.DialReadTimeout(readTimeout),
				redis.DialWriteTimeout(writeTimeout))
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
