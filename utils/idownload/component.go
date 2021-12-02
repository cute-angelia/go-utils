package idownload

import (
	"bytes"
	"fmt"
	"github.com/cute-angelia/go-utils/syntax/ifile"
	"github.com/gotomicro/ego/core/elog"
	"github.com/guonaihong/gout"
	"github.com/guonaihong/gout/dataflow"
	"image"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

const PackageName = "component.download.file"

type FileInfo struct {
	Width     int
	Height    int
	SourceUrl string
	Path      string
	Sha1      string
}

type Component struct {
	config *config
	logger *elog.Component
	locker sync.Mutex
}

// newComponent ...
func newComponent(compName string, config *config, logger *elog.Component) *Component {
	return &Component{
		config: config,
		logger: logger,
	}
}

// 请求文件
// return file body & file sha1 & error
func (c *Component) RequestFile(src string) ([]byte, string, error) {
	var body []byte
	igout := gout.GET(src).SetTimeout(c.config.Timeout)

	if c.config.Debug {
		// igout.Debug(true)
		log.Println(PackageName, "配置信息：", fmt.Sprintf("%+v", c.config))
	}

	if len(c.config.ProxySocks5) > 0 {
		igout = igout.SetSOCKS5(c.config.ProxySocks5)
	}
	if len(c.config.ProxyHttp) > 0 {
		igout = igout.SetProxy(c.config.ProxyHttp)
	}

	err := igout.SetHeader(gout.H{
		"cookie":     c.config.Cookie,
		"user-agent": c.config.UserAgent,
		"referer":    c.config.Referer,
	}).Callback(func(c *dataflow.Context) error {
		switch c.Code {
		case 200:
			c.BindBody(&body)
			return nil
		case 404: //http code为404时，服务端返回是html 字符串
			return fmt.Errorf(src + " 404")
		default:
			return fmt.Errorf(src+" error: %d", c.Code)
		}
	}).F().Retry().Attempt(3).WaitTime(time.Second * 2).MaxWaitTime(time.Second * 30).Do()

	if err != nil {
		log.Println(PackageName, "request file error -> ", src, err)
	}

	return body, ifile.FileHashSHA1(bytes.NewReader(body)), err
}

// 下载文件
func (c *Component) Download(imgurl string) (FileInfo, error) {
	var fi FileInfo
	name := ifile.NewFileName(imgurl).SetPrefix(c.config.NamePrefix).GetNameOrigin()
	if c.config.Rename {
		name = ifile.NewFileName(imgurl).SetPrefix(c.config.NamePrefix).GetNameTimeline()
	}

	if body, sha1, err := c.RequestFile(imgurl); err != nil {
		log.Println("error:", err)
		return fi, err
	} else {
		if ifileDownload, err := c.saveFile(body, name); err != nil {
			c.print("下载文件", err.Error(), "error")
			return fi, err
		} else {
			// 过滤图片
			if c.config.Width > 0 || c.config.Height > 0 {
				// 限制图片大小
				if errlimit := c.limitWidthHeightUseIsNot(ifileDownload); errlimit != nil {
					return fi, errlimit
				}
			}

			// 图片大小
			f := ifile.OpenLocalFile(ifileDownload)
			if imgconfig, _, err := image.DecodeConfig(f); err == nil {
				fi.Width = imgconfig.Width
				fi.Height = imgconfig.Height
			}
			fi.Path = name
			fi.SourceUrl = imgurl
			fi.Sha1 = sha1

			c.print("下载文件", "成功"+ifileDownload, "")
			return fi, nil
		}
	}
}

// 删除图片
func (c *Component) RemoveFile(filepath string) error {
	return os.Remove(filepath)
}

// 限制图片大小
func (c *Component) limitWidthHeightUseIsNot(localFile string) error {
	f := ifile.OpenLocalFile(localFile)
	if i, _, err := image.DecodeConfig(f); err != nil {
		c.print("限制图片大小", fmt.Sprintf("图片不存在 %s", localFile), "error")
		return err
	} else {
		if c.config.Width > 0 {
			if i.Width < c.config.Width {
				os.Remove(localFile)
				return fmt.Errorf(fmt.Sprintf("限制图片大小:小于规定宽度:%d，移除图片", c.config.Width))
			}
		}
		if c.config.Height > 0 {
			if i.Height < c.config.Height {
				os.Remove(localFile)
				return fmt.Errorf(fmt.Sprintf("限制图片大小:小于规定高度:%d，移除图片", c.config.Height))
			}
		}
	}

	return nil
}

// 保存文件
func (c *Component) saveFile(body []byte, filenamewithext string) (string, error) {
	dir := c.config.Dest

	r := bytes.NewReader(body)
	//dir
	// saveDir := path.Dir(dir)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	// 有些文件没有扩展名称
	if len(c.config.DefaultExt) > 0 {
		filenamewithext = fmt.Sprintf("%s%s", filenamewithext, c.config.DefaultExt)
	}

	// Create the file
	ifile := fmt.Sprintf("%s/%s", dir, filenamewithext)
	out, err := os.Create(ifile)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Write the body to file
	if _, err = io.Copy(out, r); err != nil {
		return "", err
	} else {
		return ifile, err
	}
}

func (c *Component) print(topic string, msg string, errtype string) {
	if errtype == "error" {
		c.logger.With(elog.FieldName(topic)).Error(msg)
	} else {
		c.logger.With(elog.FieldName(topic)).Info(msg)
	}
}
