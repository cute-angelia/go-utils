package ibunt

import (
	"fmt"
	"github.com/tidwall/buntdb"
	"log"
	"testing"
	"time"
)

func initDb() {
	New(WithName("cache"), WithDbFile("/tmp/test.db"))
}

type TestA struct {
	A string
}

func TestGetSet(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	initDb()

	testV := TestA{}

	cacheKey := "TestGetSet4"
	err1 := GetAny("cache", cacheKey, &testV)
	log.Println(err1, testV)
	// error:GetAny:TestGetSet4:缓存内容不存在
	// error:nil

	err2 := SetAny("cache", cacheKey, TestA{A: "apple"}, time.Hour)
	log.Println(err2, testV)

	err3 := GetAny("cache", cacheKey, &testV)
	log.Println(err3, testV)

	var testV2 interface{}
	err4 := GetAny("cache", cacheKey, &testV2)
	log.Println(err4, testV2)

	var testV3 string
	err5 := GetAny("cache", cacheKey, &testV3)
	log.Println(err5, testV3)
	// error:json: cannot unmarshal object into Go value of type string
}

func TestTimeout(t *testing.T) {
	initDb()
	cacheKey := "testtimeout"
	// Timeout 测试
	//Set("cache", cacheKey, "Value is ------->", time.Second*10)
	//for {
	//	data := Get("cache", cacheKey)
	//	if len(data) == 0 {
	//		break
	//	} else {
	//		log.Println(data)
	//		time.Sleep(time.Second)
	//	}
	//}

	// 测试 GetOrSet
	v, err := GetOrSet("cache", cacheKey, func() interface{} {
		return "hvzz GetOrSet"
	}, time.Minute)
	log.Println(v)
	log.Println(err)

	log.Println(Get("cache", cacheKey))
}

func TestLocker(t *testing.T) {
	initDb()
	opt1 := NewLockerOpt(WithLimit(3), WithToday(true))
	opt2 := NewLockerOpt(WithLimit(1), WithToday(false))
	opt3 := NewLockerOpt(WithLimit(35))
	for i := 1; i <= 30; i++ {
		if ok, err := IsNotLockedInLimit("cache", "test1", time.Hour, opt1); ok {
			log.Println(err)
			log.Println("opt1 > ", i)
			//2020/08/13 17:17:17 opt1 >  1
			//2020/08/13 17:17:17 opt1 >  2
			//2020/08/13 17:17:17 opt1 >  3
		}
	}

	for i := 1; i <= 30; i++ {
		if ok, err := IsNotLockedInLimit("cache", "test2", time.Hour, opt2); ok {
			log.Println(err)
			log.Println("opt2 > ", i)
			// 2020/08/13 17:17:17 opt2 >  1
		}
	}

	for i := 1; i <= 30; i++ {
		if ok, err := IsNotLockedInLimit("cache", "test3", time.Hour, opt3); ok {
			log.Println(err)
			log.Println("opt3 > ", i)
		}
	}

}

func TestShrink(t *testing.T) {
	initDb()
	db := GetDb("cache")
	db.SetConfig(buntdb.Config{
		AutoShrinkDisabled:   true,
		AutoShrinkMinSize:    10,
		AutoShrinkPercentage: 10,
	})
	db.Shrink()
}

func TestAutoShrink(t *testing.T) {

	initDb()
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
		Set("cache", fmt.Sprintf("%s%d", "test", i), "xx", time.Minute)
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
