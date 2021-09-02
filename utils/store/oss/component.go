package oss

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/cute-angelia/go-utils/utils/idownload"
	"github.com/gotomicro/ego/core/elog"
	"io"
	"log"
	"strings"
	"time"
)

type Component struct {
	name   string
	config *config
	logger *elog.Component
	Client *oss.Bucket
}

// newComponent ...
func newComponent(compName string, config *config, logger *elog.Component) *Component {
	client, err := oss.New(config.Endpoint, config.AccessKeyId, config.AccessKeySecret)
	if err != nil {
		log.Println(err.Error())
	}
	bucket, err := client.Bucket(config.BucketName)
	if err != nil {
		logger.Error("发生错误" + err.Error())
	}
	return &Component{
		name:   compName,
		config: config,
		logger: logger,
		Client: bucket,
	}
}

// 获取简短路径
func (e Component) GetPathByObjectName(objectName string) string {
	return objectName
}

// 上传文件
func (e Component) PutObject(objectNameIn string, reader io.Reader) (string, error) {
	return e.GetPathByObjectName(objectNameIn), e.Client.PutObject(objectNameIn, reader)
}

// 按文件上传
func (e Component) FPutObject(objectNameIn string, filePath string) (string, error) {
	return e.GetPathByObjectName(objectNameIn), e.Client.PutObjectFromFile(objectNameIn, filePath)
}

// 提供链接，上传, 返回key， filehash, error
func (e Component) PutObjectWithSrc(uri string, objectName string) (string, string, error) {
	// http 不处理
	if !strings.Contains(uri, "http") {
		return "", "", errors.New("链接提供不正确：" + uri)
	}
	// 更换图片到本地
	idown := idownload.Load("").Build(
		idownload.WithProxySocks5(e.config.ProxySocks5),
		idownload.WithDebug(e.config.Debug),
		idownload.WithTimeout(time.Second*20),
	)
	if filebyte, err := idown.RequestFile(uri); err != nil {
		log.Println(PackageName, "获取图片失败：❌", err)
		return "", "", errors.New("获取图片失败：❌：" + uri)
	} else {
		// 打印日志
		if e.config.Debug {
			log.Printf(PackageName, "获取图片: %s, 代理：%s", uri, e.config.ProxySocks5)
		}

		ib := bytes.NewReader(filebyte)

		log.Println("file,==>", objectName, ib.Size())

		if p, err := e.PutObject(objectName, ib); err != nil {
			log.Println(PackageName, "上传失败：❌", err, uri)
			return "", "", errors.New("获取图片失败：❌：" + uri)
		} else {
			log.Println(PackageName, "上传成功：✅", uri)

			log.Println("file222 ==>", objectName, ib.Size())

			h := sha1.New()
			io.Copy(h, ib)
			return p,  hex.EncodeToString(h.Sum(nil)), nil
		}
	}
}

// 删除文件
func (e Component) DeleteObject(filePath string) error {
	return e.Client.DeleteObject(filePath)
}
