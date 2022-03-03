/**
 * buntCache.go
 * buntdb 缓存类
 *
 * @author : Cyw
 * @email  : rose2099.c@gmail.com
 * @created: 2018/11/27 下午1:34
 * @logs   :
 *
 */
package bunt

import (
	"fmt"
	"github.com/tidwall/buntdb"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var BuntCaches sync.Map

// 初始化缓存  内部使用昵称：保存名称
func InitBuntCache(nickname string, dbname string) error {
	if _, ok := BuntCaches.Load(nickname); ok {
		return nil
	} else {
		// dbname
		db := dbname
		if !strings.Contains(dbname, ".db") {
			db = nickname + ".db"
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
					SetDb(nickname, cache)
				}
			}
		} else {
			cache.SetConfig(buntdb.Config{
				AutoShrinkDisabled:   true,
				AutoShrinkMinSize:    40,
				AutoShrinkPercentage: 30,
			})
			cache.Shrink()
			SetDb(nickname, cache)
		}
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

/**
设置
不要设置为空
*/
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

/**
获取
*/
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
