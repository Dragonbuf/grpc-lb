package config

var Conf *Config

type Config struct {
	Mysql
	Redis
}

type Mysql struct {
	MysqlMasterDns string
	MysqlSlaverDns string
}

type Redis struct {
	RedisMasterDns string
	RedisSlaverDns string
}

func init() {
	Conf = &Config{}
	Conf.MysqlMasterDns = "mysql:root"
	Conf.MysqlSlaverDns = "mysql"
	Conf.RedisMasterDns = "redis:"
	Conf.RedisSlaverDns = "redis:s"
}
