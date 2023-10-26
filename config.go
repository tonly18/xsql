package xsql

import "time"

const (
	//默认主键
	primary = "id"
	//最大存活时长
	maxLifeTime = 3600
	//最大空闲时长
	maxIdleTime = 600
	//同时最大链接数
	maxOpenConns = 50
	//最大空闲链接数
	maxIdleConns = 10
)

type Config struct {
	Host     string
	Port     int
	UserName string
	Password string
	DBName   string
	Charset  string

	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
	MaxOpenConns    int
	MaxIdleConns    int
}
