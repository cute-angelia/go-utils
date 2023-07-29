package igorm

import (
	"errors"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"sync"
	"time"
)

var gormPoolSQLite sync.Map

const PackageNameSQLite = "component.igorm.sqlite"

// GetGormSQLite 获取 gorm.DB 对象
func GetGormSQLite(dbName string) (*gorm.DB, error) {
	if v, ok := gormPoolSQLite.Load(dbName); ok {
		return v.(*gorm.DB), nil
	} else {
		return nil, errors.New(PackageNameSQLite + " 获取失败:" + dbName + " 未初始化")
	}
}

// MustInitSqlite 初始化  dir
func (c *Component) MustInitSqlite() *Component {

	var pathdb string
	if len(c.config.DbFile) > 0 {
		pathdb = c.config.DbFile
	} else {
		pathdb = fmt.Sprintf("/%s_SQLite.db", c.config.DbName)
	}

	if c.config.Debug {
		log.Println("sqlite path:", pathdb)
	}

	// log
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

	// github.com/mattn/go-sqlite3
	if db, err := gorm.Open(sqlite.Open(pathdb), &gconfig); err != nil {
		panic(err)
	} else {
		if _, ok := gormPool.Load(c.config.DbName); !ok {
			gormPoolSQLite.Store(c.config.DbName, db)
		}
	}
	return c
}
