package store

import (
	"bytes"
	"context"
	"fmt"

	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	"github.com/cute-angelia/go-utils/file"
)

var (
	QiniuSdk *Qiniu
)

/**
	qiniu
 */
func InitQiniu(ak, sk, bucket, zone string, prefix string) *Qiniu {
	z := Qiniu{}
	z.mac = qbox.NewMac(ak, sk)
	z.bucket = bucket

	switch zone {
	case "as0":
		z.zone = &storage.Zone_as0
	case "na0":
		z.zone = &storage.Zone_na0
	case "z0":
		z.zone = &storage.Zone_z0
	case "z1":
		z.zone = &storage.Zone_z1
	case "z2":
		z.zone = &storage.Zone_z2
	}

	z.prefix = prefix

	z.bucketManager = storage.NewBucketManager(z.mac, &storage.Config{
		Zone:          z.zone,
		UseHTTPS:      false,
		UseCdnDomains: false,
	})

	return &z
}

// Qiniu qiniu
type Qiniu struct {
	prefix        string
	bucket        string
	zone          *storage.Region
	mac           *qbox.Mac
	bucketManager *storage.BucketManager
}

// NewQiniu new qiniu
func NewQiniu() *Qiniu {
	return QiniuSdk
}

// Net net upload
func (q *Qiniu) Net(url string) (string, error) {
	body, err := file.GetFileWithSrcWithGout(url)
	if err != nil {
		return "", err
	}
	return q.Local(body, url)
}

// Local local upload
func (q *Qiniu) Local(body []byte, src string) (string, error) {
	putPolicy := storage.PutPolicy{
		Scope: q.bucket,
	}
	upToken := putPolicy.UploadToken(q.mac)
	formUploader := storage.NewFormUploader(&storage.Config{
		Zone:          q.zone,
		UseHTTPS:      false,
		UseCdnDomains: false,
	})
	key := file.MakeNameByUrl(true, src, q.prefix)
	data := bytes.NewReader(body)
	dataLen := int64(len(body))
	err := formUploader.Put(context.Background(), &storage.PutRet{}, upToken, key, data, dataLen, &storage.PutExtra{}) // 上传
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", key), nil
}

// Local local upload
func (q *Qiniu) LocalDiy(body []byte, key string, src string) (string, error) {
	putPolicy := storage.PutPolicy{
		Scope: q.bucket,
	}
	upToken := putPolicy.UploadToken(q.mac)
	formUploader := storage.NewFormUploader(&storage.Config{
		Zone:          q.zone,
		UseHTTPS:      false,
		UseCdnDomains: false,
	})

	if len(key) == 0 {
		key = file.MakeNameByUrl(true, src, "")
	} else {
		key = key + "/" + src
	}

	data := bytes.NewReader(body)
	dataLen := int64(len(body))
	err := formUploader.Put(context.Background(), &storage.PutRet{}, upToken, key, data, dataLen, &storage.PutExtra{}) // 上传
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", key), nil
}

// delete
// https://developer.qiniu.com/kodo/sdk/1238/go#rs-delete
func (q *Qiniu) Delete(key string) (error) {
	return q.bucketManager.Delete(q.bucket, key)
}
