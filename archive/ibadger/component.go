package ibadger

import (
	"errors"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"log"
	"strings"
	"sync"
	"time"
)

var badgerComponent *Component
var ErrClosed = errors.New("datastore closed")

const PackageName = "component.ibadger"

type Component struct {
	config *config
	DB     *badger.DB

	closeLk   sync.RWMutex
	closed    bool
	closeOnce sync.Once
	closing   chan struct{}

	syncWrites bool
}

// newComponent ...
func newComponent(config *config) *Component {
	kv, err := badger.Open(config.Options)
	if err != nil {
		if strings.HasPrefix(err.Error(), "manifest has unsupported version:") {
			err = fmt.Errorf("unsupported badger version %s", err.Error())
		}
		panic(err)
	}
	comp := &Component{
		DB:         kv,
		closing:    make(chan struct{}),
		syncWrites: config.Options.SyncWrites,
	}
	comp.config = config

	if config.GcInterval > 0 {
		go comp.periodicGC()
	}

	if err != nil {
		panic(fmt.Sprintf("[%s] 初始化失败: %s", PackageName, err))
	}

	// 赋值
	badgerComponent = comp

	return comp
}

// Keep scheduling GC's AFTER `gcInterval` has passed since the previous GC
func (d *Component) periodicGC() {
	gcTimeout := time.NewTimer(d.config.GcInterval)
	defer gcTimeout.Stop()

	for {
		select {
		case <-gcTimeout.C:
			switch err := d.gcOnce(); err {
			case badger.ErrNoRewrite, badger.ErrRejected:
				// No rewrite means we've fully garbage collected.
				// Rejected means someone else is running a GC
				// or we're closing.
				gcTimeout.Reset(d.config.GcInterval)
			case nil:
				gcTimeout.Reset(d.config.GcSleep)
			case ErrClosed:
				return
			default:
				log.Printf("error during a GC cycle: %s", err)
				// Not much we can do on a random error but log it and continue.
				gcTimeout.Reset(d.config.GcInterval)
			}
		case <-d.closing:
			return
		}
	}
}

func (d *Component) gcOnce() error {
	d.closeLk.RLock()
	defer d.closeLk.RUnlock()
	if d.closed {
		return ErrClosed
	}
	log.Println("Running GC round")
	defer log.Println("Finished running GC round")
	return d.DB.RunValueLogGC(d.config.GcDiscardRatio)
}

// newImplicitTransaction creates a transaction marked as 'implicit'.
// Implicit transactions are created by Datastore methods performing single operations.
func (d *Component) newImplicitTransaction(readOnly bool) *txn {
	return &txn{d, d.DB.NewTransaction(!readOnly), true}
}

// Implements the datastore.Txn interface, enabling transaction support for
// the badger Datastore.
type txn struct {
	ds  *Component
	txn *badger.Txn

	// Whether this transaction has been implicitly created as a result of a direct Datastore
	// method invocation.
	implicit bool
}

func (t *txn) get(key []byte) ([]byte, error) {
	item, err := t.txn.Get(key)
	if err != nil {
		return nil, err
	}
	return item.ValueCopy(nil)
}

func (t *txn) set(key []byte, value []byte, ttl time.Duration) error {
	return t.putWithTTL(key, value, ttl)
}

func (t *txn) putWithTTL(key []byte, value []byte, ttl time.Duration) error {
	return t.txn.SetEntry(badger.NewEntry(key, value).WithTTL(ttl))
}

func (t *txn) getExpiration(key []byte) (time.Time, error) {
	item, err := t.txn.Get(key)
	if err == badger.ErrKeyNotFound {
		return time.Time{}, badger.ErrKeyNotFound
	} else if err != nil {
		return time.Time{}, err
	}
	return time.Unix(int64(item.ExpiresAt()), 0), nil
}

func (t *txn) delete(key []byte) error {
	return t.txn.Delete(key)
}

func (t *txn) commit() error {
	return t.txn.Commit()
}

func (t *txn) discard() {
	t.txn.Discard()
}

// ======= functions =======

// Set 保存数据
func Set(key string, value string, ttl time.Duration) error {
	return badgerComponent.DB.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte(key), []byte(value)).WithTTL(ttl)
		err := txn.SetEntry(e)
		return err
	})
}

// Get 获取数据
func Get(key string) (string, error) {
	txn := badgerComponent.newImplicitTransaction(true)
	defer txn.discard()
	val, err := txn.get([]byte(key))
	return string(val), err
}

func GetExpiration(key string) (time.Time, error) {
	txn := badgerComponent.newImplicitTransaction(false)
	defer txn.discard()
	return txn.getExpiration([]byte(key))
}

func Delete(key string) error {
	txn := badgerComponent.newImplicitTransaction(false)
	defer txn.discard()

	err := txn.delete([]byte(key))
	if err != nil {
		return err
	}
	return txn.commit()
}
