package igorm

import (
	"errors"
	"github.com/jinzhu/gorm"
	otgorm "github.com/smacker/opentracing-gorm"
	"log"
	"sync"
	"time"
)

type Component struct {
	config   *config
	locker   sync.Mutex
	gormPool sync.Map
}

// newComponent ...
func newComponent(config *config) *Component {
	return &Component{
		config: config,
	}
}

// 初始化
func (c *Component) InitGorm(dsn string) *gorm.DB {
	// 配置必须信息
	c.config.Dsn = dsn

	if v, ok := c.gormPool.Load(c.config.DbName); ok {
		return v.(*gorm.DB)
	} else {
		db := c.initDb()
		c.gormPool.Store(c.config.DbName, db)
		return db
	}
}

func (c *Component) GetGorm() (*gorm.DB, error) {
	if v, ok := c.gormPool.Load(c.config.DbName); ok {
		return v.(*gorm.DB), nil
	} else {
		return nil, errors.New("失败，未能获取:" + c.config.DbName)
	}
}

func (c *Component) initDb() *gorm.DB {
	var db *gorm.DB
	var err error

	for db, err = gorm.Open(c.config.DbType, c.config.Dsn); err != nil; {
		log.Println(packageName+"数据库连接异常", c.config.DbName, err)
		time.Sleep(5 * time.Second)
		db, err = gorm.Open(c.config.DbType, c.config.Dsn)
	}

	//最大闲置打开连接数
	db.DB().SetMaxIdleConns(c.config.MaxIdleConns)

	// ==> 最大连接数
	db.DB().SetMaxOpenConns(c.config.MaxOpenConns)

	// 最大时间
	db.DB().SetConnMaxLifetime(c.config.MaxLifetime)

	// LogMode Print log
	db.LogMode(c.config.LogDebug)

	// with logger
	db.SetLogger(c.config.Logger)

	// 设置 callback
	otgorm.AddGormCallbacks(db)

	return db
}
