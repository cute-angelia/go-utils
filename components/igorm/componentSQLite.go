package igorm

import (
	"errors"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"path"
	"sync"
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
func (c *Component) MustInitSqlite(dir string) *Component {
	if len(dir) == 0 {
		dir = "./"
	}

	pathdb := fmt.Sprintf(fmt.Sprintf(path.Clean(dir)+"/%s_SQLite.db", c.config.DbName))

	if c.config.Debug {
		log.Println("sqlite path:", pathdb)
	}

	// github.com/mattn/go-sqlite3
	if db, err := gorm.Open(sqlite.Open(pathdb), &gorm.Config{}); err != nil {
		panic(err)
	} else {
		if _, ok := gormPool.Load(c.config.DbName); !ok {
			gormPoolSQLite.Store(c.config.DbName, db)
		}
	}
	return c
}
