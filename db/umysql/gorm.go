package umysql

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/smacker/opentracing-gorm"
	"log"
	"time"
)

// gorm 连接池, gorm 本身也提供了pool, 提供全局对象, 这里初始化全局对象
var Gorm = make(map[string]*gorm.DB)

// 配置, 为了断开重连
var GormOpts = make(map[string]GormOptions)

// 初始化 GORM
func InitGorm(opts GormOptions) {
	orm, ok := Gorm[opts.Dbname]
	if !ok {
		orm = initDB(opts)
		Gorm[opts.Dbname] = orm

		GormOpts[opts.Dbname] = opts
	}
}

// 获取 Gorm
func GetGorm(dbname string) *gorm.DB {
	db, ok := Gorm[dbname]
	if !ok {
		log.Println("没有初始化" + dbname)
	} else {
		err := db.DB().Ping()
		if err != nil {
			Gorm[dbname] = initDB(GormOpts[dbname])
		}
	}
	return Gorm[dbname]
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
		log.Println("数据库连接异常", dbName, err)
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
