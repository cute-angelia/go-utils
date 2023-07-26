package igorm

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"regexp"
	"time"
)

const PackageNameMysql = "component.igorm.mysql"

// GetGormMysql 获取 gorm.DB 对象
func GetGormMysql(dbName string) (*gorm.DB, error) {
	if v, ok := gormPool.Load(dbName); ok {
		return v.(*gorm.DB), nil
	} else {
		return nil, errors.New(packageName + " 获取失败:" + dbName + " 未初始化")
	}
}

// MustInitMysql 初始化
func (c *Component) MustInitMysql() *Component {
	// 配置必须信息
	if len(c.config.Dsn) == 0 || len(c.config.DbName) == 0 {
		panic(fmt.Sprintf("❌数据库配置不正确 dbName=%s dsn=%s", c.config.DbName, c.config.Dsn))
	}
	// 初始化 db
	if _, ok := gormPool.Load(c.config.DbName); !ok {
		gormPool.Store(c.config.DbName, c.initMysqlDb())
	}

	// 初始化日志
	log.Println(fmt.Sprintf("[%s] Name:%s 初始化",
		PackageNameMysql,
		c.config.DbName,
	))

	return c
}

func (c *Component) initMysqlDb() *gorm.DB {
	// log.Println(packageName, "初始化数据库", c.config.DbName)
	var db *gorm.DB
	var err error

	var vlog = new(log.Logger)
	if c.config.LoggerWriter == nil {
		vlog = log.New(os.Stdout, "\r\n", log.LstdFlags|log.Lshortfile)
	} else {
		vlog = log.New(c.config.LoggerWriter, "", 0)
	}

	newLogger := logger.New(
		vlog, // io writer
		logger.Config{
			SlowThreshold: time.Second,       // Slow SQL threshold
			LogLevel:      c.config.LogLevel, // Log level
			Colorful:      true,
		},
	)

	gconfig := gorm.Config{
		Logger: newLogger,
	}

	for db, err = gorm.Open(mysql.Open(c.config.Dsn), &gconfig); err != nil; {
		re := regexp.MustCompile(`tcp\((.*?)\)`)
		match := re.FindStringSubmatch(c.config.Dsn)
		if len(match) > 1 {
			log.Println(packageName, "❌数据库连接异常", match[1], c.config.DbName, err)
		}
		time.Sleep(5 * time.Second)
		db, err = gorm.Open(mysql.Open(c.config.Dsn), &gconfig)
	}

	if idb, err := db.DB(); err != nil {
		log.Println(packageName, "❌数据库获取异常", c.config.DbName, err)
		return nil
	} else {
		// ==>  用于设置连接池中空闲连接的最大数量(10)
		idb.SetMaxIdleConns(c.config.MaxIdleConns)

		// ==>  设置打开数据库连接的最大数量(100)
		idb.SetMaxOpenConns(c.config.MaxOpenConns)

		// 最大时间
		idb.SetConnMaxLifetime(c.config.MaxLifetime)

		// 设置 callback
		// otgorm.AddGormCallbacks(db)

		return db
	}
}
