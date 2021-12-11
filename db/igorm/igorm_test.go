package igorm

import (
	"github.com/cute-angelia/go-utils/syntax/ijson"
	"testing"
)

type MinioPostModel struct {
	Id           int32  `json:"id" gorm:"primary_key"`
	Title        string `json:"title" gorm:"column:title"`
	Desc         string `json:"desc" gorm:"column:desc"`
	Cover        string `json:"cover" gorm:"column:cover"`
	Idols        string `json:"idols" gorm:"column:idols"`
	Tags         string `json:"tags" gorm:"column:tags"`
	Categorys    string `json:"categorys" gorm:"column:categorys"`
	SourcePage   string `json:"source_page" gorm:"column:source_page"`
	BucketHash   string `json:"bucket_hash" gorm:"column:bucket_hash"`
	Bucket       string `json:"bucket" gorm:"column:bucket"`
	BucketGroup  int32  `json:"bucket_group" gorm:"column:bucket_group"`
	BucketPrefix string `json:"bucket_prefix" gorm:"column:bucket_prefix"`
	BucketType   int32  `json:"bucket_type" gorm:"column:bucket_type"`
	Dateline     int64  `json:"dateline" gorm:"column:dateline"`
	CollectNum   int32  `json:"collect_num" gorm:"column:collect_num"`
	ZanNum       int32  `json:"zan_num" gorm:"column:zan_num"`
	Views        int32  `json:"views" gorm:"column:views"`
	Liked        bool   `json:"liked" gorm:"-"`
	CoverUrl     string `json:"cover_url" gorm:"-"`
}

func (MinioPostModel) TableName() string {
	return "minio_post"
}

func TestConnect(t *testing.T) {
	dbName := "dtwk_meter"
	dsn := "root:Luck2018@@tcp(host-tx.aaqq.in:43306)/orange?charset=utf8mb4&parseTime=true&loc=Local"
	Load(dbName).Build(
		//WithLogDebug(true),
		WithMaxIdleConns(1),
		WithMaxOpenConnss(5),
		WithDsn(dsn),
	).MustInitMysql()

	orm, _ := GetGormMysql(dbName)
	model := MinioPostModel{}
	orm.First(&model)
	t.Log(ijson.Pretty(model))

	orm2, _ := GetGormMysql(dbName)
	model2 := MinioPostModel{} // 不新建数据有问题，和上面一样
	orm2.Debug().Order("id desc").First(&model2)
	t.Log(ijson.Pretty(model2))
}
