package igorm

import (
	"errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

// 获取 gorm.DB 对象
func GetGormMysql(dbName string) (*gorm.DB, error) {
	if v, ok := gormPool.Load(dbName); ok {
		return v.(*gorm.DB), nil
	} else {
		return nil, errors.New(packageName + " 获取失败:" + dbName + " 未初始化")
	}
}

// 初始化
func (c *Component) MustInitMysql() *Component {
	// 配置必须信息
	if len(c.config.Dsn) == 0 {
		log.Println(packageName, "❌数据库配置不正确，dsn未设置", c.config.DbName)
		panic(c.config.DbName + "❌数据库配置不正确，dsn未设置")
	}
	// 初始化 db
	if _, ok := gormPool.Load(c.config.DbName); !ok {
		gormPool.Store(c.config.DbName, c.initMysqlDb())
	}
	return c
}

func (c *Component) initMysqlDb() *gorm.DB {
	log.Println(packageName, "初始化数据库", c.config.DbName)

	var db *gorm.DB
	var err error
	gconfig := gorm.Config{
		Logger: c.config.Logger,
	}

	for db, err = gorm.Open(mysql.Open(c.config.Dsn), &gconfig); err != nil; {
		log.Println(packageName, "❌数据库连接异常", c.config.DbName, err)
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

		// LogMode Print log
		// idb.LogMode(c.config.LogDebug)

		// with logger
		//if &c.config.Logger != nil {
		//	db.SetLogger(c.config.Logger)
		//}

		// 设置 callback
		// otgorm.AddGormCallbacks(db)

		return db
	}
}
