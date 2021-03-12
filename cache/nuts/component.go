package nuts

import (
	"fmt"
	"github.com/xujiajun/nutsdb"
	"log"
	"strconv"
	"sync"
	"time"
)

var componentPool map[string]*Component

func init() {
	componentPool = make(map[string]*Component)
}

type Component struct {
	config *config
	db     *nutsdb.DB
	locker sync.Mutex
}

func newComponent(cfg *config) *Component {
	if comp, ok := componentPool[cfg.Key]; ok {
		com := Component{
			config: cfg,
		}
		com.db = comp.db
		return &com
	} else {
		db, err := nutsdb.Open(cfg.Options)
		if err != nil {
			log.Fatal(err)
		}
		return &Component{
			config: cfg,
			db:     db,
		}
	}
}

/**
	true => 非锁定状态，处理正常逻辑
	false => 锁定状态，处理错误逻辑
    if !bunt.IsNotLockedInLimit("cache", "SIGNPRE_REPEAT_"+nonce, time.Minute*60, bunt.NewLockerOpt(bunt.WithToday(true))) {
*/
func (self Component) IsNotLockedInLimit(bucket string, key string, ttl uint32, opts LockerOpts) bool {
	if opts.Today {
		key = fmt.Sprintf("%s_%s", key, time.Now().Format("2006-01-02"))
	}
	if opts.Uid > 0 {
		key = fmt.Sprintf("%s_%d", key, opts.Uid)
	}

	value := self.Get(bucket, key)
	if len(value) > 0 {
		n, _ := strconv.Atoi(value)
		if n >= opts.Limit {
			return false
		} else {
			self.Set(bucket, key, fmt.Sprintf("%d", n+1), ttl)
			return true
		}
	} else {
		self.Set(bucket, key, "1", ttl)
		return true
	}
}

func (self Component) Set(bucket string, key string, val string, ttl uint32) {
	if err := self.db.Update(
		func(tx *nutsdb.Tx) error {
			key := []byte(key)
			val := []byte(val)

			// 如果设置 ttl = 0 or Persistent, 这个key就会永久不删除
			// 这边 ttl = 60 , 60s之后就会过期。
			if err := tx.Put(bucket, key, val, ttl); err != nil {
				return err
			}
			return nil
		}); err != nil {

		logError("Set", err)
	}
}

func (self Component) Incr(bucket string, key string, val string, ttl uint32) int {
	self.locker.Lock()
	defer self.locker.Unlock()

	vcache := self.Get(bucket, key)
	vint := 0
	if vcacheint, err := strconv.Atoi(vcache); err == nil {
		vint = vcacheint
	}
	valint, _ := strconv.Atoi(val)
	self.Set(bucket, key, fmt.Sprintf("%d", vint+valint), ttl)

	return vint + valint
}

func (self Component) Get(bucket string, key string) string {
	res := ""
	if err := self.db.View(
		func(tx *nutsdb.Tx) error {
			key := []byte(key)
			if e, err := tx.Get(bucket, key); err != nil {
				return err
			} else {
				res = string(e.Value)
			}
			return nil
		}); err != nil {
		// log
		// logError("Get", err)
	}

	return res
}

func (self Component) SetByte(bucket string, key string, val []byte, ttl uint32) {
	if err := self.db.Update(
		func(tx *nutsdb.Tx) error {
			key := []byte(key)

			// 如果设置 ttl = 0 or Persistent, 这个key就会永久不删除
			// 这边 ttl = 60 , 60s之后就会过期。
			if err := tx.Put(bucket, key, val, ttl); err != nil {
				return err
			}
			return nil
		}); err != nil {

		logError("Set", err)
	}
}

func (self Component) GetByte(bucket string, key string) []byte {
	res := []byte{}
	if err := self.db.View(
		func(tx *nutsdb.Tx) error {
			key := []byte(key)
			if e, err := tx.Get(bucket, key); err != nil {
				return err
			} else {
				res = e.Value
			}
			return nil
		}); err != nil {

		// log
		logError("Get", err)
	}

	return res
}

func (self Component) Delete(bucket string, key string) error {
	return self.db.Update(
		func(tx *nutsdb.Tx) error {
			key := []byte(key)
			if err := tx.Delete(bucket, key); err != nil {
				return err
			}
			return nil
		})
}

func (self Component) PrefixScan(bucket string, prefix string, limit int) (nutsdb.Entries, error) {
	resentries := nutsdb.Entries{}

	if err := self.db.View(
		func(tx *nutsdb.Tx) error {
			prefix := []byte(prefix)
			// 从offset=0开始 ，限制 100 entries 返回
			if entries, err := tx.PrefixScan(bucket, prefix, limit); err != nil {
				return err
			} else {
				resentries = entries
				//for _, entry := range entries {
				//	fmt.Println(string(entry.Key), string(entry.Value))
				//}
			}
			return nil
		}); err != nil {

		logError("PrefixScan", err)

		return nil, err
	} else {
		return resentries, nil
	}
}

func (self Component) Merge() {
	err := self.db.Merge()
	if err != nil {
		logError("Merge", err)
	} else {
		log.Println(time.Now().Format("2006-01-02 03:04:05"), "Merge")
	}
}

func (self Component) Close() {
	err := self.db.Close()
	if err != nil {
		logError("Close", err)
	}
}

func logError(key string, err error) {
	log.Println(PackageName, ":error:"+key+":", err)
}
