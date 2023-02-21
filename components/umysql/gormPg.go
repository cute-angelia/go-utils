package umysql

import (
	"github.com/jinzhu/gorm"
	"github.com/smacker/opentracing-gorm"
	"log"
	"time"
)

// gorm 连接池, gorm 本身也提供了pool, 提供全局对象, 这里初始化全局对象
var GormPg = make(map[string]*gorm.DB)
// 配置, 为了断开重连
var GormPgOpts = make(map[string]GormOptions)

// 初始化 GORM
func InitGormPg(opts GormOptions) {
	orm, ok := GormPg[opts.Dbname]
	if !ok {
		orm = initPgDB(opts)
		GormPg[opts.Dbname] = orm

		GormPgOpts[opts.Dbname] = opts
	}
}

// 获取 Gorm
func GetGormPg(dbname string) *gorm.DB {
	db, ok := GormPg[dbname]
	if !ok {
		log.Println("没有初始化" + dbname)
	} else {
		err := db.DB().Ping()
		if err != nil {
			GormPg[dbname] = initPgDB(GormPgOpts[dbname])
		}
	}
	return GormPg[dbname]
}

//  内部初始化 DB
func initPgDB(opts GormOptions) *gorm.DB {
	var (
		db  *gorm.DB
		err error
	)

	// config
	maxIdleConns := opts.MaxIdleConns
	maxOpenConns := opts.MaxOpenConns
	dbType := "postgres"

	connectString := ""
	connectString = opts.Dsn

	//db, err := gorm.Open("postgres", "host=myhost port=myport user=gorm dbname=gorm password=mypassword")
	//defer db.Close()

	for db, err = gorm.Open(dbType, connectString); err != nil; {
		log.Println("数据库连接异常", err)
		time.Sleep(5 * time.Second)
		db, err = gorm.Open(dbType, connectString)
	}

	//最大打开连接数
	db.DB().SetMaxIdleConns(maxIdleConns)

	if maxOpenConns == 0 {
		db.DB().SetMaxOpenConns(maxIdleConns * 2)
	} else {
		db.DB().SetMaxOpenConns(maxOpenConns)
	}

	// LogMode Print log
	db.LogMode(opts.LogDebug)

	// 设置 callback
	otgorm.AddGormCallbacks(db)

	return db
}
