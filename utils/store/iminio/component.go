package iminio

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/cute-angelia/go-utils/cache/bunt"
	"github.com/cute-angelia/go-utils/utils/encrypt/hash"
	"github.com/cute-angelia/go-utils/utils/idownload"
	"github.com/gotomicro/ego/core/elog"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"log"
	"math/rand"
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

// newComponent ...
func newComponent(compName string, config *config, logger *elog.Component) *Component {
	minioClient, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccesskeyId, config.SecretaccessKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		logger.Error("发生错误" + err.Error())
	}
	return &Component{
		name:   compName,
		config: config,
		logger: logger,
		Client: minioClient,
	}
}

// 获取链接 不带bucket
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
		log.Println(PackageName, "Successfully URL: ", presignedURL)
		bunt.Set("cache", hashkey, presignedURL.String(), t)
		return presignedURL.String(), nil
	}
}

// 获取链接 链接带 bucket
func (e *Component) SignCoverWithCache(cover string, t time.Duration) string {
	if strings.Contains(cover, "http") {
		return cover
	}

	if len(cover) > 0 {
		cover = strings.TrimLeft(cover, "/")
		temp := strings.Split(cover, "/")
		objkey := temp[1:len(temp)]
		icover, _ := e.SignUrlWithCache(temp[0], strings.Join(objkey, "/"), t)
		return icover
	} else {
		return ""
	}
}

// 生成缓存
func (e *Component) GenerateHashKey(bucketType int32, bucket string, prefix string) string {
	return hash.NewEncodeMD5(fmt.Sprintf("%d%s%s", bucketType, bucket, prefix))
}

// Objects 获取 分页
func (e *Component) GetObjectsByPage(bucket string, prefix string, page int32, perpage int32) (objs []string, notall bool) {
	// 控制流程
	count := int32(0)
	offset := (page - 1) * perpage

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	opt := minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	}
	objectCh := e.Client.ListObjects(ctx, bucket, opt)

	//maxSize := perpage * page
	//for i := int32(0); i < maxSize; i++ {
	//}

	for object := range objectCh {
		if object.Err == nil {
			// log.Printf("---->1 count:%d, offset:%d, perpage:%d, %v", count, offset, perpage, count >= offset)
			// 小于当前游标
			if count >= offset {
				// 当前计数 - 游标
				// log.Printf("<---- count:%d, offset:%d, perpage:%d, %v false:继续", count, offset, perpage, count-offset >= perpage)
				if count-offset >= perpage {
					notall = true
					cancel()
					break
				}
				//if img_url, err := e.SignUrlWithCache(bucket, object.Key, time.Hour*24*6); err == nil {
				objs = append(objs, bucket+"/"+object.Key)
				count++
				//}
			} else {
				count++
			}
		} else {
			log.Println("object.Err", object.Err)
		}
	}
	return objs, notall
}

// Objects 状态
func (e Component) GetObjectStat(bucket string, objectName string) (minio.ObjectInfo, error) {
	objInfo, err := e.Client.StatObject(context.Background(), bucket, objectName, minio.StatObjectOptions{})
	if err != nil {
		log.Println(err)
	}
	return objInfo, err
}

func (e Component) CheckMode(objectName string) (newObjectName string, canupload bool) {
	// 跳过
	if e.config.ReplaceMode == ReplaceModeIgnore {
		canupload = false
		newObjectName = objectName
	}
	if e.config.ReplaceMode == ReplaceModeReplace {
		canupload = true
		newObjectName = objectName
	}
	if e.config.ReplaceMode == ReplaceModeTwo {
		rand.Seed(time.Now().Unix())
		canupload = true
		newObjectName = fmt.Sprintf("bak_%d_%s", rand.Intn(100), objectName)
	}
	return
}

// 上传文件
func (e Component) PutObject(bucket string, objectNameIn string, reader io.Reader, objectSize int64, objopt minio.PutObjectOptions) (minio.UploadInfo, error) {
	if objectName, ok := e.CheckMode(objectNameIn); ok {
		objectName = strings.Replace(objectName, "//", "/", -1)
		uploadInfo, err := e.Client.PutObject(context.Background(), bucket, objectName, reader, objectSize, objopt)
		if err != nil {
			log.Println(bucket, objectNameIn, err)
			return uploadInfo, err
		}
		if e.config.Debug {
			log.Println("Successfully uploaded bytes: ", uploadInfo)
		}
		return uploadInfo, err
	} else {
		return minio.UploadInfo{}, fmt.Errorf("模式未设置 %s", objectNameIn)
	}
}

// 按文件上传
func (e Component) FPutObject(bucket string, objectNameIn string, filePath string, objopt minio.PutObjectOptions) (minio.UploadInfo, error) {
	if objectName, ok := e.CheckMode(objectNameIn); ok {
		uploadInfo, err := e.Client.FPutObject(context.Background(), bucket, objectName, filePath, objopt)
		if err != nil {
			log.Println(err)
			return uploadInfo, err
		}
		if e.config.Debug {
			log.Println("Successfully uploaded bytes: ", uploadInfo)
		}
		return uploadInfo, err
	} else {
		return minio.UploadInfo{}, fmt.Errorf("模式未设置 %s", objectNameIn)
	}
}

// 提供链接，上传到 minio
// return key & hash sha1 & error
func (e Component) PutObjectWithSrc(dnComponent *idownload.Component, uri string, bucket string, objectName string, objopt minio.PutObjectOptions) (string, string, error) {
	// http 不处理
	if !strings.Contains(uri, "http") {
		return uri, "", errors.New("非链接地址:" + uri)
	}
	// 更换图片到本地
	//idown := idownload.Load("").Build(
	//	idownload.WithProxySocks5(e.config.ProxySocks5),
	//	idownload.WithDebug(e.config.Debug),
	//	idownload.WithTimeout(e.config.Timeout),
	//	idownload.WithReferer(e.config.Referer),
	//)
	if filebyte, sha1, err := dnComponent.RequestFile(uri); err != nil {
		if e.config.Debug {
			log.Println(PackageName, "获取图片失败：❌", uri, err)
		}
		return "", "", errors.New("获取图片失败：❌" + uri + "  " + err.Error())
	} else {
		// 打印日志
		if e.config.Debug {
			log.Printf(PackageName, "获取图片: %s, 代理：%s", uri, e.config.ProxySocks5)
		}
		objectName = strings.Replace(objectName, "//", "/", -1)

		if info, err := e.Client.PutObject(context.TODO(), bucket, objectName, bytes.NewReader(filebyte), int64(len(filebyte)), objopt); err != nil {
			if e.config.Debug {
				log.Println(PackageName, "上传失败：❌", err, bucket, objectName, uri)
			}
			return "", "", fmt.Errorf("上传失败：❌ %v, %s %s %s", err, bucket, objectName, uri)
		} else {
			uri = bucket + "/" + info.Key
			log.Println(PackageName, "上传成功：✅", bucket, objectName, uri)
			return uri, sha1, nil
		}
	}
}

// 删除文件
func (e Component) DeleteObject(objectNameWithBucket string) error {
	opts := minio.RemoveObjectOptions{
	}
	bucket, objectName := e.GetBucketAndObjectName(objectNameWithBucket)

	log.Println(PackageName, "删除对象1", objectNameWithBucket, bucket, objectName)

	err := e.Client.RemoveObject(context.Background(), bucket, objectName, opts)
	if err != nil {
		log.Println(PackageName, "删除对象失败：❌", err, objectNameWithBucket, bucket, objectName)
		return err
	}
	return nil
}

// 根据路径获取bucket 和 object name
func (e Component) GetBucketAndObjectName(objectNameWithBucket string) (string, string) {
	if len(objectNameWithBucket) > 0 {
		objectNameWithBucket = strings.TrimLeft(objectNameWithBucket, "/")
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

func (e Component) GetConfig() {
	log.Println(PackageName, "配置信息：", fmt.Sprintf("%+v", e.config))
}
