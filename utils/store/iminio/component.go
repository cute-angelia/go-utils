package iminio

import (
	"bytes"
	"context"
	"fmt"
	"github.com/cute-angelia/go-utils/cache/bunt"
	"github.com/cute-angelia/go-utils/utils/encrypt/hash"
	"github.com/cute-angelia/go-utils/utils/idownload"
	"github.com/gotomicro/ego/core/elog"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"log"
	"net/url"
	"strings"
	"sync"
	"time"
)

const PackageName = "component.store.iminio"

type Component struct {
	name   string
	config *config
	logger *elog.Component
	locker sync.Mutex
	Client *minio.Client
}

var comp *Component
var once2 sync.Once

// newComponent ...
func newComponent(compName string, config *config, logger *elog.Component) *Component {
	once2.Do(func() {
		minioClient, err := minio.New(config.Endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(config.AccesskeyId, config.SecretaccessKey, ""),
			Secure: config.UseSSL,
		})
		if err != nil {
			logger.Error("发生错误" + err.Error())
		}
		comp = &Component{
			name:   compName,
			config: config,
			logger: logger,
			Client: minioClient,
		}
	})

	return comp
}

//获取链接 不带bucket
func (e *Component) SignUrlWithCache(bucket string, key string, t time.Duration) (string, error) {
	hashkey := e.GenerateHashKey(1, bucket, key)
	cachedata := bunt.Get("cache", hashkey)

	if len(cachedata) > 3 {
		return cachedata, nil
	}

	reqParams := make(url.Values)
	if presignedURL, err := e.Client.PresignedGetObject(context.Background(), bucket, key, t, reqParams); err != nil {
		e.logger.Info(err.Error())
		return "", err
	} else {
		log.Println("Successfully URL: ", presignedURL)
		bunt.Set("cache", hashkey, presignedURL.String(), t)
		return presignedURL.String(), nil
	}
}

//获取链接 链接带 bucket
func (e *Component) SignCoverWithCache(cover string, t time.Duration) string {
	if strings.Contains(cover, "http") {
		return cover
	}

	if len(cover) > 0 {
		temp := strings.Split(cover, "/")
		objkey := temp[1:len(temp)]
		icover, _ := e.SignUrlWithCache(temp[0], strings.Join(objkey, "/"), t)
		return icover
	} else {
		return ""
	}
}

// 用于缓存
func (e *Component) GenerateHashKey(bucketType int32, bucket string, prefix string) string {
	return hash.NewEncodeMD5(fmt.Sprintf("%d%s%s", bucketType, bucket, prefix))
}

// 上传文件
func (e Component) PutObject(bucket string, objectName string, reader io.Reader, objectSize int64, objopt minio.PutObjectOptions) minio.UploadInfo {
	uploadInfo, err := e.Client.PutObject(context.Background(), bucket, objectName, reader, objectSize, objopt)
	if err != nil {
		fmt.Println(err)
		return uploadInfo
	}
	fmt.Println("Successfully uploaded bytes: ", uploadInfo)
	return uploadInfo
}

// 提供链接，上传到 minio
func (e Component) PutObjectWithSrc(uri string, bucket string, objectName string, objopt minio.PutObjectOptions) string {
	// http 不处理
	if !strings.Contains(uri, "http") {
		return uri
	}
	// 更换图片到本地
	idown := idownload.Load("").Build(idownload.WithProxySocks5(e.config.ProxySocks5))
	if filebyte, err := idown.RequestFile(uri); err != nil {
		log.Println("获取图片失败：❌", err)
		return ""
	} else {
		// 打印日志
		if e.config.Debug {
			log.Printf("获取图片: %s, 代理：%s", uri, e.config.ProxySocks5)
		}

		if info, err := e.Client.PutObject(context.TODO(), bucket, objectName, bytes.NewReader(filebyte), int64(len(filebyte)), objopt); err != nil {
			log.Println("上传失败：❌", err)
		} else {
			uri = bucket + "/" + info.Key
			log.Println("上传成功：✅", uri)
		}
		return uri
	}
}

// 删除文件
func (e Component) DeleteObject(objectNameWithBucket string) error {
	opts := minio.RemoveObjectOptions{
	}
	bucket, objectName := e.GetBucketAndObjectName(objectNameWithBucket)

	log.Println("删除对象1", objectNameWithBucket, bucket, objectName)

	err := e.Client.RemoveObject(context.Background(), bucket, objectName, opts)
	if err != nil {
		log.Println("删除对象失败：❌", err, objectNameWithBucket, bucket, objectName)
		return err
	}
	return nil
}

// 根据路径获取bucket 和 object name
func (e Component) GetBucketAndObjectName(objectNameWithBucket string) (string, string) {
	if len(objectNameWithBucket) > 0 {
		temp := strings.Split(objectNameWithBucket, "/")
		if len(temp) > 1 {
			objkey := temp[1:len(temp)]
			return temp[0], strings.Join(objkey, "/")
		} else {
			return objectNameWithBucket, ""
		}
	} else {
		return "", ""
	}
}

// 获取默认PutObjectOptions
// video/mp4,video/webm,video/ogg
func (e Component) GetPutObjectOptions(contentType string) minio.PutObjectOptions {
	if len(contentType) > 0 {
		return minio.PutObjectOptions{ContentType: contentType}
	}
	return minio.PutObjectOptions{ContentType: "image/jpeg,image/png,image/jpeg"}
}
