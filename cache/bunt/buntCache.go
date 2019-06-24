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
	"github.com/tidwall/buntdb"
	"time"
	"fmt"
	"log"
	"os"
)

var BuntCaches map[string]*buntdb.DB

type BuntCache struct {
}

/**
	初始化
 */
func (self *BuntCache) InitDb(name string, db string) error {
	if BuntCaches[name] != nil {
		return nil
	} else {
		// 无法创建
		if cache, err := buntdb.Open(db); err != nil {
			// 移除异常 db
			if _, e := os.Stat(db); e == nil {
				os.Remove(db)

				// 重新打开一次
				if cache, err := buntdb.Open(db); err != nil {
					return err
				} else {
					self.SetDb(name, cache)
				}
			}
		} else {
			self.SetDb(name, cache)
		}

		return nil
	}
}

func (self *BuntCache) GetDb(name string) *buntdb.DB {
	if BuntCaches[name] == nil {
		log.Println("buntdb.未初始化:" + name)
		return nil
	} else {
		return BuntCaches[name]
	}
}

func (self *BuntCache) SetDb(name string, db *buntdb.DB) {
	if BuntCaches == nil {
		BuntCaches = map[string]*buntdb.DB{}
	}
	BuntCaches[name] = db
}

/**
	设置
 */
func (self *BuntCache) Set(dbname string, key string, val string, ttl time.Duration) error {
	if db := self.GetDb(dbname); db != nil {
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
func (self *BuntCache) Get(dbname string, key string) (string, error) {
	if db := self.GetDb(dbname); db != nil {
		val := ""
		db.View(func(tx *buntdb.Tx) error {
			val, _ = tx.Get(key)
			return nil
		})
		return val, nil
	} else {
		return "", fmt.Errorf("无法找到 db")
	}
}

/**
	查询是否锁定,
	true => 我被锁住了
	false => 没有锁
 */
func (self *BuntCache) IsLocked(dbname string, key string, val string, ttl time.Duration) (bool, error) {
	value, _ := self.Get(dbname, key)
	if len(value) > 0 {
		return true, nil
	} else {
		self.Set(dbname, key, val, ttl)
		return false, nil
	}
}
