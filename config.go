package xsql

import (
	"time"
)

const (
	//最大存活时长
	maxLifeTime = 1800
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

	MaxLifetime  time.Duration
	MaxIdleTime  time.Duration
	MaxOpenConns int
	MaxIdleConns int

	//Logger mysql.Logger //日志记录器
}
