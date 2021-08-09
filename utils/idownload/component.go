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
)

const PackageName = "component.download.file"

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
func (c *Component) RequestFile(src string) ([]byte, error) {
	var body []byte
	igout := gout.GET(src).SetTimeout(c.config.Timeout)

	if c.config.Debug {
		log.Println("配置信息：", fmt.Sprintf("%+v", c.config))
	}

	if len(c.config.ProxySocks5) > 0 {
		log.Println("proxy:", c.config.ProxySocks5)
		igout = igout.SetSOCKS5(c.config.ProxySocks5)
	}
	if len(c.config.ProxyHttp) > 0 {
		log.Println("proxy:", c.config.ProxyHttp)
		igout = igout.SetProxy(c.config.ProxyHttp)
	}

	err := igout.SetHeader(gout.H{
		"cookie":     c.config.Cookie,
		"user-agent": c.config.UserAgent,
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
	}).Do()

	if err != nil {
		log.Println("request file error -> ", src, err)
	}
	return body, err
}

// 下载文件
func (c *Component) Download(imgurl string) {
	name := ifile.NewFileName(imgurl).GetNameOrigin(c.config.NamePrefix)
	if c.config.Rename {
		name = ifile.NewFileName(imgurl).GetNameTimeline(c.config.NamePrefix)
	}

	if body, err := c.RequestFile(imgurl); err != nil {
		return
	} else {
		if ifileDownload, err := c.saveFile(body, name); err != nil {
			c.print("下载文件", err.Error(), "error")
		} else {
			// 过滤图片
			if c.config.Width > 0 || c.config.Height > 0 {
				if !c.limit(ifileDownload) {
					c.print("下载文件", "成功"+ifileDownload, "")
				}
			} else {
				c.print("下载文件", "成功"+ifileDownload, "")
			}
		}
	}
}

// 限制图片大小
func (c *Component) limit(localFile string) bool {
	f := ifile.OpenLocalFile(localFile)
	if i, _, err := image.DecodeConfig(f); err != nil {
		c.print("限制图片大小", fmt.Sprintf("图片不存在 %s", localFile), "error")
		return true
	} else {
		if c.config.Width > 0 {
			if i.Width < c.config.Width {
				os.Remove(localFile)
				c.print("限制图片大小", fmt.Sprintf("小于规定宽度:%d，移除图片", c.config.Width), "error")
				return true
			}
		}
		if c.config.Height > 0 {
			if i.Height < c.config.Height {
				os.Remove(localFile)
				c.print("限制图片大小", fmt.Sprintf("小于规定g高度:%d，移除图片", c.config.Width), "error")
				return true
			}
		}
	}

	return false
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
