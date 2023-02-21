package ibadger

import (
	"github.com/dgraph-io/badger/v3"
	"log"
	"testing"
	"time"
)

func TestStart(t *testing.T) {
	// disk
	//New()

	// Memory
	// New(WithMemory(16 << 20))

	//t.Log(Set("test", "ok", time.Second*2))
	//t.Log(Get("test"))
	//<-time.After(time.Second * 3)
	//t.Log(Get("test"))
}

func TestOrigin(t *testing.T) {
	// Open the Badger database located in the /tmp/badger directory.
	// It will be created if it doesn't exist.
	opt := badger.DefaultOptions("/tmp/badger")
	//opt.NumVersionsToKeep = 0
	//opt.CompactL0OnClose = true
	//opt.ValueLogFileSize = 1024 * 1024 * 1

	db, err := badger.Open(opt)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// Your code here…

	err = db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte("answer"), []byte("42")).WithTTL(time.Second * 3)
		err := txn.SetEntry(e)
		return err
	})

	log.Println(err)

	db.View(func(txn *badger.Txn) error {
		if item, err := txn.Get([]byte("answer")); err != nil {
			return err
		} else {
			if valCopy, err := item.ValueCopy(nil); err != nil {
				return err
			} else {
				log.Println(string(valCopy))
				return nil
			}
		}
	})

	//ticker := time.NewTicker(5 * time.Second)
	//defer ticker.Stop()
	//for range ticker.C {
	//again:
	//	err := db.RunValueLogGC(0.02)
	//	if err == nil {
	//		goto again
	//	}
	//}

	<-time.After(time.Second * 20)
}
