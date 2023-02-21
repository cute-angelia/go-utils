package ibuntV2

import (
	"errors"
	"fmt"
	"github.com/tidwall/buntdb"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var BuntCaches sync.Map
var iComponent *Component

const PackageName = "component.ibunt"

type Component struct {
	config *config
}

func (c *Component) GenerateCacheKey(bucket string, key string) string {
	return fmt.Sprintf("%s:%s", bucket, key)
}

func GetComponent(dbname string) *Component {
	iComponent.config.Name = dbname
	return iComponent
}

// newComponent ...
func newComponent(config *config) *Component {
	comp := &Component{}
	comp.config = config

	//  自动初始化， 有些人经常忘。。。
	if err := comp.initBuntDb(); err != nil {
		log.Println(fmt.Sprintf("[%s] 初始化失败", PackageName))
	}

	// 赋值
	iComponent = comp

	return comp
}

// initBuntDb 初始化
func (c *Component) initBuntDb() error {
	if _, ok := BuntCaches.Load(c.config.Name); ok {
		return nil
	} else {
		// dbname
		dbfile := c.config.DbFile
		if !strings.Contains(dbfile, ".db") {
			dbfile = dbfile + ".db"
		}
		// 无法创建
		if cache, err := buntdb.Open(dbfile); err != nil {
			// 移除异常 db
			if _, e := os.Stat(dbfile); e == nil {
				os.Remove(dbfile)

				// retry 重新打开一次
				if dbopen, err := buntdb.Open(dbfile); err != nil {
					panic(err)
				} else {
					dbopen.SetConfig(buntdb.Config{
						AutoShrinkDisabled:   true,
						AutoShrinkMinSize:    40,
						AutoShrinkPercentage: 30,
					})
					dbopen.Shrink()
					SetDb(c.config.Name, dbopen)
				}
			} else {
				log.Println("error:", dbfile, e)
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
		log.Println(fmt.Sprintf("[%s] 初始化成功 Name:%s, dbFile=%s",
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

func ShowTttl(dbname string, key string) string {
	if db := GetDb(dbname); db != nil {
		val := ""
		db.View(func(tx *buntdb.Tx) error {
			val, _ = tx.Get(key)

			itemTTL, err := tx.TTL(key)
			if err == nil {
				exat := time.Now().Add(itemTTL)
				val = exat.Format(time.RFC3339)
			}
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

func (c *Component) Get(key string) (string, error) {
	if db := GetDb(c.config.Name); db != nil {
		val := ""
		db.View(func(tx *buntdb.Tx) error {
			val, _ = tx.Get(key)
			return nil
		})
		if len(val) == 0 {
			return "", nil
		} else {
			if len(val) <= 10 {
				return "", nil
			} else {
				endTime := val[:10]
				vv := val[len(val)-10:]
				endTimeInt, _ := strconv.Atoi(endTime)
				if time.Now().Before(time.Unix(int64(endTimeInt), 0)) {
					return vv, nil
				} else {
					return "", nil
				}
			}
		}
	} else {
		log.Println(fmt.Errorf("无法找到 db" + c.config.Name))
		return "", nil
	}
}

func (c *Component) GetMulti(keys []string) map[string]string {
	result := make(map[string]string)
	for _, key := range keys {
		result[key], _ = c.Get(key)
	}
	return result
}

func (c *Component) Set(key string, value string, ttl time.Duration) error {
	endTime := time.Now().Add(ttl)
	value = fmt.Sprintf("%d%s", endTime.Unix(), value)
	if db := GetDb(c.config.Name); db != nil {
		db.Update(func(tx *buntdb.Tx) error {
			tx.Set(key, value, &buntdb.SetOptions{Expires: true, TTL: ttl})
			return nil
		})
		return nil
	} else {
		return fmt.Errorf("无法找到 db")
	}
}

func (c *Component) SetWithBucket(bucket string, key string, value string, ttl time.Duration) error {
	endTime := time.Now().Add(ttl)
	value = fmt.Sprintf("%d%s", endTime.Unix(), value)
	if db := GetDb(c.config.Name); db != nil {
		db.CreateIndex(bucket, bucket+":*", buntdb.IndexString)
		db.Update(func(tx *buntdb.Tx) error {
			tx.Set(key, value, &buntdb.SetOptions{Expires: true, TTL: ttl})
			return nil
		})
		return nil
	} else {
		return fmt.Errorf("无法找到 db")
	}
}

func (c *Component) Contains(key string) bool {
	v, _ := c.Get(key)
	return len(v) > 0
}

func (c *Component) Delete(key string) error {
	return Delete(c.config.Name, key)
}

func (c *Component) Flush() error {
	return GetDb(c.config.Name).Shrink()
}

func (c *Component) Scan(bucket string, f func(key string) error) (err error) {
	if db := GetDb(c.config.Name); db != nil {
		return db.View(func(tx *buntdb.Tx) error {
			err := tx.Ascend(bucket, func(key, value string) bool {
				f(key)
				return true
			})
			return err
		})
	} else {
		return errors.New(fmt.Sprintf("无法找到 db" + c.config.Name))
	}
}

func (c *Component) Fold(f func(key string) error) (err error) {
	panic("implement me")
}
