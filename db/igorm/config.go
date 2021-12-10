package igorm

import (
	"github.com/jinzhu/gorm"
	"time"
)

const packageName = "component.db.gorm"

// config options
type config struct {
	DbName       string
	Dsn          string
	MaxIdleConns int           // SetMaxIdleConns 用于设置连接池中空闲连接的最大数量(10)
	MaxOpenConns int           // SetMaxOpenConns 设置打开数据库连接的最大数量(100)
	MaxLifetime  time.Duration // SetConnMaxLifetime 设置了连接可复用的最大时间。 time.Hour
	LogDebug     bool
	Logger       gorm.Logger
	DbType       string
}

// DefaultConfig 返回默认配置
func DefaultConfig() *config {
	return &config{
		MaxIdleConns: 10,
		MaxOpenConns: 100,
		MaxLifetime:  time.Hour,
		LogDebug:     false,
		DbType:       "mysql",
	}
}
