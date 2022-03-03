package bunt

import (
	"fmt"
	"github.com/cute-angelia/go-utils/utils/generator/RandomString"
	"github.com/tidwall/buntdb"
	"log"
	"testing"
	"time"
)

func TestTimeout(t *testing.T) {
	InitBuntCache("cache", "/tmp/test.db")
	cacheKey := "testtimeout"
	Set("cache", cacheKey, "Value is ------->", time.Second*10)
	for {
		data := Get("cache", cacheKey)
		if len(data) == 0 {
			break
		} else {
			log.Println(data)
			time.Sleep(time.Second)
		}
	}
}

func TestLocker(t *testing.T) {

	dbname := "cache_locker"

	if err := InitBuntCache(dbname, fmt.Sprintf("/tmp/test_%d.db", time.Now().Unix())); err != nil {
		log.Println(err)
		return
	}

	opt1 := NewLockerOpt(WithLimit(3), WithToday(true))
	opt2 := NewLockerOpt(WithLimit(1), WithToday(false))
	opt3 := NewLockerOpt(WithLimit(35))

	for i := 1; i <= 30; i++ {
		if IsNotLockedInLimit(dbname, "test1", time.Hour, opt1) {
			log.Println("opt1 > ", i)
			//2020/08/13 17:17:17 opt1 >  1
			//2020/08/13 17:17:17 opt1 >  2
			//2020/08/13 17:17:17 opt1 >  3
		}
	}

	for i := 1; i <= 30; i++ {
		if IsNotLockedInLimit(dbname, "test2", time.Hour, opt2) {
			log.Println("opt2 > ", i)
			// 2020/08/13 17:17:17 opt2 >  1
		}
	}

	for i := 1; i <= 30; i++ {
		if IsNotLockedInLimit(dbname, "test3", time.Hour, opt3) {
			log.Println("opt3 > ", i)
		}
	}

}

func TestShrink(t *testing.T) {
	if err := InitBuntCache("cache", "/tmp/test_1.db"); err != nil {
		log.Println(err)
		return
	}
	db := GetDb("cache")
	db.SetConfig(buntdb.Config{
		AutoShrinkDisabled:   true,
		AutoShrinkMinSize:    10,
		AutoShrinkPercentage: 10,
	})
	db.Shrink()
}

func TestAutoShrink(t *testing.T) {

	if err := InitBuntCache("cache", "/tmp/test_1.db"); err != nil {
		log.Println(err)
		return
	}

	db := GetDb("cache")

	db.SetConfig(buntdb.Config{
		AutoShrinkDisabled:   true,
		AutoShrinkMinSize:    1,
		AutoShrinkPercentage: 0,
	})

	for i := 0; i < 20; i++ {
		writeCache()
		time.Sleep(time.Second * 10)
	}

	// defer db.Shrink()

	log.Println("end")
}

func writeCache() {
	for i := 0; i < 10000; i++ {
		log.Println("--> ", i)
		Set("cache", fmt.Sprintf("%s%d", "test", i), RandomString.RandomString(332), time.Minute)
	}
}

func TestAutoShrink2(t *testing.T) {
	dbpath := "/tmp/test_2.db"
	// 无法创建
	if db, err := buntdb.Open(dbpath); err == nil {
		db.SetConfig(buntdb.Config{
			AutoShrinkDisabled:   true,
			AutoShrinkMinSize:    10,
			AutoShrinkPercentage: 30,
		})
		for i := 0; i < 20; i++ {
			writeCache2(db)
			time.Sleep(time.Second * 5)
		}
	}
}

func writeCache2(db *buntdb.DB) {
	for i := 0; i < 10000; i++ {
		log.Println("--> ", i)
		val := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		db.Update(func(tx *buntdb.Tx) error {
			tx.Set(fmt.Sprintf("%s%d", "test", i), val, &buntdb.SetOptions{Expires: true, TTL: time.Hour})
			return nil
		})
	}
}
