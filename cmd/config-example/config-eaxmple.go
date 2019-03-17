package config

import "time"

var (
	/**
	MaxActive 最大连接数，即最多的tcp连接数，一般建议往大的配置，但不要超过操作系统文件句柄个数（centos下可以ulimit -n查看）。
	MaxIdle 最大空闲连接数，即会有这么多个连接提前等待着，但过了超时时间也会关闭。
	IdleTimeout 空闲连接超时时间，但应该设置比redis服务器超时时间短。否则服务端超时了，客户端保持着连接也没用。
	Wait 这是个很有用的配置:如果超过最大连接，是报错，还是等待。
	*/
	RedisHost           = "localhost:6379"
	RedisPassword       = ""
	RedisDb             = 0
	RedisConTimeout     = 5 * time.Second
	RedisReadTimeout    = 2 * time.Second
	RedisWriteTimeout   = 1 * time.Second
	RedisMaxIdle        = 1000
	RedisMaxActive      = 3000
	RedisMaxIdleTimeout = 5 * time.Second

	// mysql
	MysqlMasterDns    = "root:1234@/makaplatv4?charset=utf8&parseTime=True&loc=Local"
	MysqlMaxIdleConns = 100
	MysqlMaxOpenConns = 1000

	//可用逗号分割，开启etcd 集群
	EtcDHost   = "http://localhost:2379"
	EtcdPrefix = "/etcd3_naming"
	ETCDEndpoints   = []string{"http://localhost:2379"} // etcd 集群配置
	ETCDDialTimeout = 2 * time.Second                   // etcd 超时时间
)
