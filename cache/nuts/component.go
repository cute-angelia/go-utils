package nuts

import (
	"github.com/xujiajun/nutsdb"
	"log"
	"time"
)

var componentPool map[string]*Component

func init() {
	componentPool = make(map[string]*Component)
}

type Component struct {
	config *config
	db     *nutsdb.DB
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
		logError("Get", err)
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
