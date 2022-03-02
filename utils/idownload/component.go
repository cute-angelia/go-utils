package idownload

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/cute-angelia/go-utils/syntax/ifile"
	"github.com/gotomicro/ego/core/elog"
	"github.com/guonaihong/gout"
	"github.com/guonaihong/gout/dataflow"
	"github.com/schollz/progressbar/v3"
	"image"
	"io"
	"log"
	"os"
	"path"
	"sync"
	"time"
)

const PackageName = "component.download.file"

var (
	ErrorNotFound = errors.New("404 file not found")
)

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
	if config.Debug {
		log.Println(PackageName, "配置信息：", fmt.Sprintf("%+v", config))
	}

	return &Component{
		config: config,
		logger: logger,
	}
}

// RequestFile 请求文件，返回 字节 & sha1 & error
// return file body & file sha1 & error
func (c *Component) RequestFile(src string) ([]byte, string, error) {
	bar := progressbar.NewOptions(100,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetDescription("Download "+src),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))
	for i := 0; i < 100; i++ {
		bar.Add(1)
		time.Sleep(5 * time.Millisecond)
	}

	var body []byte
	igout := gout.GET(src).SetTimeout(c.config.Timeout)

	if len(c.config.ProxySocks5) > 0 {
		igout = igout.SetSOCKS5(c.config.ProxySocks5)
	}
	if len(c.config.ProxyHttp) > 0 {
		igout = igout.SetProxy(c.config.ProxyHttp)
	}

	// 用于解析 服务端 返回的http header
	type RspHeader struct {
		ContentLength int64 `header:"content-length"`
	}
	var head RspHeader
	err := igout.SetHeader(gout.H{
		"cookie":     c.config.Cookie,
		"user-agent": c.config.UserAgent,
		"referer":    c.config.Referer,
	}).BindHeader(&head).Callback(func(c *dataflow.Context) error {
		// 进度条
		switch c.Code {
		case 200:
			c.BindBody(&body)
			bar.Finish()
			return nil
		case 404: //http code为404时，服务端返回是html 字符串
			return ErrorNotFound
		default:
			return fmt.Errorf(src+" error: %d", c.Code)
		}
	}).F().Retry().Attempt(5).WaitTime(time.Second * 3).MaxWaitTime(time.Second * 30).Do()

	if err != nil {
		log.Println(PackageName, "request file error -> ", src, err)
	}

	return body, ifile.FileHashSHA1(bytes.NewReader(body)), err
}

// Download 下载文件
// 参数
//	uri 下载文件路径
//	saveName 保存路径带后缀
func (c *Component) Download(uri string, name string) (FileInfo, error) {
	var fi FileInfo

	if body, sha1, err := c.RequestFile(uri); err != nil {
		log.Println("error:", err)
		return fi, err
	} else {
		if ifileDownload, err := c.saveFile(body, name); err != nil {
			log.Println("下载文件❌", uri, err.Error())
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
			if f, err := ifile.OpenLocalFile(ifileDownload); err != nil {
				return fi, errors.New("文件获取失败" + ifileDownload)
			} else {

				defer f.Close()
				if imgconfig, _, err := image.DecodeConfig(f); err == nil {
					fi.Width = imgconfig.Width
					fi.Height = imgconfig.Height
				}
				fi.Path = name
				fi.SourceUrl = uri
				fi.Sha1 = sha1

				log.Println("下载文件✅", uri, "成功:"+ifileDownload)
				return fi, nil
			}
		}
	}
}

// DownloadWithProgressbar 下载文件带进度条
func (c *Component) DownloadWithProgressbar(uri string, name string) (FileInfo, error) {
	// 设置字节
	bar := progressbar.NewOptions(-1,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(10),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() {
			fmt.Printf("\n")
		}),
		progressbar.OptionSetDescription("Download "+uri),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	var fi FileInfo
	var body []byte
	igout := gout.GET(uri).SetTimeout(c.config.Timeout)

	if len(c.config.ProxySocks5) > 0 {
		igout = igout.SetSOCKS5(c.config.ProxySocks5)
	}
	if len(c.config.ProxyHttp) > 0 {
		igout = igout.SetProxy(c.config.ProxyHttp)
	}

	// 用于解析 服务端 返回的http header
	type RspHeader struct {
		ContentLength int `header:"content-length"`
	}
	var head RspHeader
	err := igout.SetHeader(gout.H{
		"cookie":     c.config.Cookie,
		"user-agent": c.config.UserAgent,
		"referer":    c.config.Referer,
	}).BindHeader(&head).Callback(func(c *dataflow.Context) error {
		// 进度条
		switch c.Code {
		case 200:
			c.BindBody(&body)
			return nil
		case 404: //http code为404时，服务端返回是html 字符串
			return ErrorNotFound
		default:
			return fmt.Errorf(uri+" error: %d", c.Code)
		}
	}).F().Retry().Attempt(5).WaitTime(time.Second * 3).MaxWaitTime(time.Second * 30).Do()

	if err == nil {
		// 设置字节
		f, _ := os.OpenFile(name, os.O_CREATE|os.O_WRONLY, 0644)
		defer f.Close()
		io.Copy(io.MultiWriter(f, bar), bytes.NewReader(body))
		bar.Finish()

		fi.Path = name
		fi.SourceUrl = uri
	}

	return fi, err
}

// RemoveFile 删除图片
func (c *Component) RemoveFile(filepath string) error {
	return os.Remove(filepath)
}

// 限制图片大小
func (c *Component) limitWidthHeightUseIsNot(localFile string) error {
	if f, err := ifile.OpenLocalFile(localFile); err != nil {
		return err
	} else {
		defer f.Close()
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
}

// 保存文件
func (c *Component) saveFile(body []byte, filenamewithext string) (string, error) {
	dir := ""

	// 空目标文件
	if len(c.config.Dest) == 0 {
		dir = path.Dir(filenamewithext) + "/"
	} else {
		if !path.IsAbs(filenamewithext) {
			dir = c.config.Dest + "/" + path.Dir(filenamewithext) + "/"
		} else {
			dir = c.config.Dest + path.Dir(filenamewithext) + "/"
		}
	}

	// 清理
	dir = path.Clean(dir)
	if dir == "." {
		dir = dir + "/"
	}

	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			//dir
			// saveDir := path.Dir(dir)
			err := os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				panic(err)
			}
		}
	}

	r := bytes.NewReader(body)

	// 有些文件没有扩展名称
	if len(c.config.DefaultExt) > 0 {
		filenamewithext = fmt.Sprintf("%s%s", filenamewithext, c.config.DefaultExt)
	}

	// Create the file
	out, err := os.Create(filenamewithext)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Write the body to file
	if _, err = io.Copy(out, r); err != nil {
		return "", err
	} else {
		return filenamewithext, err
	}
}

func (c *Component) print(topic string, msg string, errtype string) {
	if errtype == "error" {
		c.logger.With(elog.FieldName(topic)).Error(msg)
	} else {
		c.logger.With(elog.FieldName(topic)).Info(msg)
	}
}
