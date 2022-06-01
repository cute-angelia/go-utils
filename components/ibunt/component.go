package ibunt

import (
	"fmt"
	"github.com/cute-angelia/go-utils/syntax/ijson"
	"github.com/tidwall/buntdb"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var BuntCaches sync.Map

const PackageName = "component.ibunt"

var RedisPools sync.Map

type Component struct {
	config *config
}

// newComponent ...
func newComponent(config *config) *Component {
	comp := &Component{}
	comp.config = config

	//  自动初始化， 有些人经常忘。。。
	if err := comp.initBuntDb(); err != nil {
		log.Println(fmt.Sprintf("[%s] 初始化失败", PackageName))
	}

	return comp
}

// initBuntDb 初始化
func (c Component) initBuntDb() error {
	if _, ok := BuntCaches.Load(c.config.Name); ok {
		return nil
	} else {
		// dbname
		db := c.config.DbFile
		if !strings.Contains(db, ".db") {
			db = db + ".db"
		}
		// 无法创建
		if cache, err := buntdb.Open(db); err != nil {
			// 移除异常 db
			if _, e := os.Stat(db); e == nil {
				os.Remove(db)

				// retry 重新打开一次
				if cache, err := buntdb.Open(db); err != nil {
					panic(err)
				} else {
					cache.SetConfig(buntdb.Config{
						AutoShrinkDisabled:   true,
						AutoShrinkMinSize:    40,
						AutoShrinkPercentage: 30,
					})
					cache.Shrink()
					SetDb(c.config.Name, cache)
				}
			}
		} else {
			cache.SetConfig(buntdb.Config{
				AutoShrinkDisabled:   true,
				AutoShrinkMinSize:    40,
				AutoShrinkPercentage: 30,
			})
			cache.Shrink()
			SetDb(c.config.Name, cache)
		}
		// 初始化日志
		log.Println(fmt.Sprintf("[%s] Name:%s, dbFile=%s 初始化",
			PackageName,
			c.config.Name,
			c.config.DbFile,
		))
		return nil
	}
}

func SetConfig(db *buntdb.DB, conf buntdb.Config) {
	if conf.AutoShrinkMinSize > 0 {
		conf.AutoShrinkDisabled = true
		err := db.SetConfig(conf)
		if err != nil {
			log.Println("buntdb.DB SetConfig error", err)
		}
	} else {
		err := db.SetConfig(buntdb.Config{
			AutoShrinkDisabled:   true,
			AutoShrinkMinSize:    50,
			AutoShrinkPercentage: 100,
			SyncPolicy:           buntdb.Always,
		})
		if err != nil {
			log.Println("buntdb.DB SetConfig error", err)
		} else {
			var config buntdb.Config
			if err := db.ReadConfig(&config); err != nil {
				log.Println("buntdb.DB GetConfig error", err)
			} else {
				// log.Println("ReadConfig", config)
			}
		}
	}
}

func GetDb(name string) *buntdb.DB {
	if v, ok := BuntCaches.Load(name); ok {
		return v.(*buntdb.DB)
	} else {
		log.Println("buntdb.未初始化:" + name)
		return nil
	}
}

func SetDb(name string, db *buntdb.DB) {
	SetConfig(db, buntdb.Config{})
	BuntCaches.Store(name, db)
}

func Set(dbname string, key string, val string, ttl time.Duration) error {
	if db := GetDb(dbname); db != nil {
		db.Update(func(tx *buntdb.Tx) error {
			tx.Set(key, val, &buntdb.SetOptions{Expires: true, TTL: ttl})
			return nil
		})
		return nil
	} else {
		return fmt.Errorf("无法找到 db")
	}
}

// GetOrSet 获取数据或者存储数据
func GetOrSet(dbname string, key string, function func() interface{}, ttl time.Duration) (string, error) {
	cacheData := Get(dbname, key)
	if len(cacheData) > 10 {
		return cacheData, nil
	} else {
		byteJson, _ := ijson.Marshal(function())
		strJson := string(byteJson)
		err := Set(dbname, key, strJson, ttl)
		return strJson, err
	}
}

func Get(dbname string, key string) string {
	if db := GetDb(dbname); db != nil {
		val := ""
		db.View(func(tx *buntdb.Tx) error {
			val, _ = tx.Get(key)
			return nil
		})

		if len(val) == 0 {
			return ""
		} else {
			return val
		}
	} else {
		log.Println(fmt.Errorf("无法找到 db" + dbname))
		return ""
	}
}

func Delete(dbname string, key string) error {
	if db := GetDb(dbname); db != nil {
		return db.Update(func(tx *buntdb.Tx) error {
			_, err := tx.Delete(key)
			return err
		})
	} else {
		log.Println(fmt.Errorf("无法找到 db" + dbname))
		return fmt.Errorf("无法找到 db" + dbname)
	}
}
