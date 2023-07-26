package umysql

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/smacker/opentracing-gorm"
	"log"
	"regexp"
	"sync"
	"time"
)

// gorm 连接池, gorm 本身也提供了pool, 提供全局对象, 这里初始化全局对象
// make(map[string]*gorm.DB)
var Gorm sync.Map

// 配置, 为了断开重连 // make(map[string]GormOptions)
var GormOpts sync.Map

// 初始化 GORM
func InitGorm(opts GormOptions) {
	_, ok := Gorm.Load(opts.Dbname)
	if !ok {
		Gorm.Store(opts.Dbname, initDB(opts))
		GormOpts.Store(opts.Dbname, opts)
	}
}

// 获取 Gorm
func GetGorm(dbname string) *gorm.DB {
	v, ok := Gorm.Load(dbname)
	if !ok {
		log.Println("Umysql-DB没有初始化:" + dbname)
		return nil
	} else {
		db := v.(*gorm.DB)
		err := db.DB().Ping()
		if err != nil {
			opts, _ := GormOpts.Load(dbname)
			gormDb := initDB(opts.(GormOptions))
			Gorm.Store(dbname, gormDb)
			return gormDb
		} else {
			return v.(*gorm.DB)
		}
	}
}

// 获取 Gorm Ctx
func GetGormCtx(ctx context.Context, dbName string) *gorm.DB {
	db := otgorm.SetSpanToGorm(ctx, GetGorm(dbName))
	return db
}

//  内部初始化 DB
func initDB(opts GormOptions) *gorm.DB {
	var (
		db  *gorm.DB
		err error
	)

	// config
	dbHost := opts.Host
	dbName := opts.Dbname
	dbUser := opts.Username
	dbPasswd := opts.Password
	dbPort := opts.Port
	maxIdleConns := opts.MaxIdleConns
	maxOpenConns := opts.MaxOpenConns
	dbType := "mysql"

	connectString := ""
	if len(opts.Dsn) > 0 {
		connectString = opts.Dsn
	} else {
		connectString = dbUser + ":" + dbPasswd + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai"
	}

	for db, err = gorm.Open(dbType, connectString); err != nil; {
		re := regexp.MustCompile(`tcp\((.*?)\)`)
		match := re.FindStringSubmatch(connectString)
		if len(match) > 1 {
			log.Println("数据库连接异常", match[1], dbName, err)
		}
		time.Sleep(5 * time.Second)
		db, err = gorm.Open(dbType, connectString)
	}

	//最大闲置打开连接数
	db.DB().SetMaxIdleConns(maxIdleConns)

	// ==> 最大连接数
	db.DB().SetMaxOpenConns(maxOpenConns)

	// 最大时间
	if opts.MaxLifetime > 0 {
		db.DB().SetConnMaxLifetime(opts.MaxLifetime)
	}

	// LogMode Print log
	db.LogMode(opts.LogDebug)

	// with logger
	if opts.Logger.LogWriter != nil {
		db.SetLogger(opts.Logger)
	}

	// 设置 callback
	otgorm.AddGormCallbacks(db)

	return db
}
