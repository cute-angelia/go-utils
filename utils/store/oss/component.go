package oss

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/cute-angelia/go-utils/syntax/ifile"
	"github.com/cute-angelia/go-utils/utils/idownload"
	"github.com/gotomicro/ego/core/elog"
	"io"
	"log"
	"net/url"
	"path"
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

// 从完整链接中，获取 objKey
func (e Component) GetObjectKeyByUrl(uri string) string {
	if up, err := url.Parse(uri); err != nil {
		return uri
	} else {
		if path.IsAbs(up.Path) {
			return up.Path[1:]
		} else {
			return up.Path
		}
	}
}

// 拼接完整路径
func (e Component) JoinUrl(objectKey string) string {
	return fmt.Sprintf("%s/%s", e.config.BucketHost, objectKey)
}

func (e Component) CleanObjKey(objectKey string) string {
	if path.IsAbs(objectKey) {
		return objectKey[1:]
	} else {
		return objectKey
	}
}

// 上传文件
func (e Component) PutObject(objectNameIn string, reader io.Reader) (string, error) {
	return e.CleanObjKey(objectNameIn), e.Client.PutObject(objectNameIn, reader)
}

// 按文件上传
func (e Component) FPutObject(objectNameIn string, filePath string) (string, error) {
	return e.CleanObjKey(objectNameIn), e.Client.PutObjectFromFile(objectNameIn, filePath)
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
	if filebyte, _, err := idown.RequestFile(uri); err != nil {
		log.Println(PackageName, "获取图片失败：❌", err)
		return "", "", errors.New("获取图片失败：❌：" + uri)
	} else {
		// 打印日志
		if e.config.Debug {
			log.Printf(PackageName+"获取图片: %s, 代理：%s", uri, e.config.ProxySocks5)
		}

		if p, err := e.PutObject(objectName, bytes.NewReader(filebyte)); err != nil {
			log.Println(PackageName, "上传失败：❌", err, uri)
			return "", "", errors.New("获取图片失败：❌：" + uri)
		} else {
			log.Printf(PackageName+"上传成功：✅ %s => %s", uri, objectName)
			return p, ifile.FileHashSHA1(bytes.NewReader(filebyte)), nil
		}
	}
}

// 删除文件
func (e Component) DeleteObject(filePath string) error {
	return e.Client.DeleteObject(filePath)
}
