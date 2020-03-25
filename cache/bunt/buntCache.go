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
	"strings"
)

var BuntCaches map[string]*buntdb.DB

// 初始化缓存  内部使用昵称：保存名称
func InitBuntCache(nickname string, dbname string) error {
	if BuntCaches[nickname] != nil {
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
					SetDb(nickname, cache)
				}
			}
		} else {
			SetDb(nickname, cache)
		}
		return nil
	}
}

func GetDb(name string) *buntdb.DB {
	if BuntCaches[name] == nil {
		log.Println("buntdb.未初始化:" + name)
		return nil
	} else {
		return BuntCaches[name]
	}
}

func SetDb(name string, db *buntdb.DB) {
	if BuntCaches == nil {
		BuntCaches = map[string]*buntdb.DB{}
	}
	BuntCaches[name] = db
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
func Get(dbname string, key string) (string, error) {
	if db := GetDb(dbname); db != nil {
		val := ""
		db.View(func(tx *buntdb.Tx) error {
			val, _ = tx.Get(key)
			return nil
		})

		if len(val) == 0 {
			return "", fmt.Errorf("数据为空")
		} else {
			return val, nil
		}
	} else {
		return "", fmt.Errorf("无法找到 db")
	}
}

/**
	查询是否锁定,
	true => 我被锁住了， 不操作业务
	false => 没有锁， 操作业务
 */
func IsLocked(dbname string, key string, val string, ttl time.Duration) (bool, error) {
	value, _ := Get(dbname, key)
	if len(value) > 0 {
		return true, nil
	} else {
		Set(dbname, key, val, ttl)
		return false, nil
	}
}
