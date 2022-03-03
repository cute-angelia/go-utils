package idownload

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gotomicro/ego/core/elog"
	"github.com/guonaihong/gout"
	"github.com/guonaihong/gout/dataflow"
	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"sync"
	"time"
)

const PackageName = "component.download.file"

var (
	ErrorNotFound = errors.New("404 file not found")
)

type FileInfo struct {
	SourceUrl string
	Path      string
}

type Component struct {
	config *config
	logger *elog.Component
	// 进度条
	bar *progressbar.ProgressBar
}

// newComponent ...
func newComponent(compName string, config *config, logger *elog.Component) *Component {
	//if config.Debug {
	//	log.Println(PackageName, "配置信息：", fmt.Sprintf("%+v", config))
	//}
	comp := &Component{}
	comp.config = config
	comp.logger = logger

	return comp
}

// getGoHttpClient 业务定制头部
func (d *Component) getGoHttpClient(uri string, method string) *dataflow.DataFlow {
	var igout *dataflow.DataFlow
	switch method {
	case "GET":
		igout = gout.GET(uri)
	case "POST":
		igout = gout.POST(uri)
	case "PUT":
		igout = gout.PUT(uri)
	case "DELETE":
		igout = gout.DELETE(uri)
	case "HEAD":
		igout = gout.HEAD(uri)
	case "OPTIONS":
		igout = gout.OPTIONS(uri)
	default:
		igout = gout.GET(uri)
	}
	// 设置过期时间
	if d.config.Timeout > 0 {
		igout = igout.SetTimeout(d.config.Timeout)
	}
	if len(d.config.ProxySocks5) > 0 {
		igout = igout.SetSOCKS5(d.config.ProxySocks5)
	}
	if len(d.config.ProxyHttp) > 0 {
		igout = igout.SetProxy(d.config.ProxyHttp)
	}
	return igout.SetHeader(gout.H{
		"Cookie":        d.config.Cookie,
		"User-Agent":    d.config.UserAgent,
		"Referer":       d.config.Referer,
		"Authorization": "Bearer " + d.config.Authorization,
	})
}

// Download 下载文件
func (d *Component) Download(strURL, filename string) (FileInfo, error) {
	if filename == "" {
		filename = path.Base(strURL)
	}
	header := http.Header{}
	var statusCode int
	if err := d.getGoHttpClient(strURL, "HEAD").BindHeader(&header).Code(&statusCode).Do(); err == nil {
		if statusCode == http.StatusOK && header.Get("Accept-Ranges") == "bytes" {
			contentLength, _ := strconv.Atoi(header.Get("Content-Length"))
			return d.multiDownload(strURL, filename, contentLength)
		}
		if statusCode == http.StatusNotFound {
			return FileInfo{}, ErrorNotFound
		}
	}
	return d.singleDownload(strURL, filename)
}

// DownloadToByteRetry 请求文件，返回 字节
func (d *Component) DownloadToByteRetry(src string, retry int) ([]byte, error) {
	var body []byte
	err := d.getGoHttpClient(src, "GET").Callback(func(c *dataflow.Context) error {
		// 进度条
		switch c.Code {
		case 200:
			c.BindBody(&body)
			return nil
		case 404: //http code为404时，服务端返回是html 字符串
			return ErrorNotFound
		default:
			return fmt.Errorf(src+" error: %d", c.Code)
		}
	}).F().Retry().Attempt(retry).WaitTime(time.Second * 3).Do()

	if err != nil {
		log.Println(PackageName, "DownloadToByte error -> ", src, err)
	}
	return body, err
}

// DownloadToByte 请求文件，返回 字节
func (d *Component) DownloadToByte(strURL string) ([]byte, error) {

	iClient := d.getGoHttpClient(strURL, "GET").Client()
	req, err := http.NewRequest("GET", strURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := iClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// 进度条
	d.getBar(int(resp.ContentLength), strURL)

	var buf bytes.Buffer
	bufcache := make([]byte, 32*1024)
	// You don't need use "bufio.NewWriter(&b)" to create an io.Writer. &b is an io.Writer itself.
	// bufio.NewWriter(&buf)
	io.CopyBuffer(io.MultiWriter(&buf, d.bar), resp.Body, bufcache)
	return buf.Bytes(), nil
}

// RemoveFile 删除图片
func (d *Component) RemoveFile(filepath string) error {
	return os.Remove(filepath)
}

// multiDownload 并发下载
func (d *Component) multiDownload(strURL, filename string, contentLen int) (FileInfo, error) {
	var info FileInfo

	d.getBar(contentLen, strURL)

	partSize := contentLen / d.config.Concurrency

	// 创建部分文件的存放目录
	partDir := d.getPartDir(filename)
	os.Mkdir(partDir, 0777)
	defer os.RemoveAll(partDir)

	var wg sync.WaitGroup
	wg.Add(d.config.Concurrency)

	rangeStart := 0

	for i := 0; i < d.config.Concurrency; i++ {
		go func(i, rangeStart int) {
			defer wg.Done()

			rangeEnd := rangeStart + partSize
			// 最后一部分，总长度不能超过 ContentLength
			if i == d.config.Concurrency-1 {
				rangeEnd = contentLen
			}

			downloaded := 0
			if d.config.Resume {
				partFileName := d.getPartFilename(filename, i)
				content, err := os.ReadFile(partFileName)
				if err == nil {
					downloaded = len(content)
				}
				d.bar.Add(downloaded)
			}

			d.downloadPartial(strURL, filename, rangeStart+downloaded, rangeEnd, i)

		}(i, rangeStart)

		rangeStart += partSize + 1
	}

	wg.Wait()

	if err := d.merge(filename); err == nil {
		info.SourceUrl = strURL
		info.Path = filename
		return info, nil
	} else {
		log.Println("合并文件发生错误，", filename, err)
		return info, err
	}
}

//  singleDownload 直接下载
func (d *Component) singleDownload(strURL, filename string) (FileInfo, error) {
	var info FileInfo

	// 需要进度条，更要复用 http.client 这里就不使用原生的 http
	// resp, err := http.Get(strURL)
	// Transport: &http.Transport{
	//	MaxIdleConnsPerHost: 10000,
	// },
	iClient := d.getGoHttpClient(strURL, "GET").Client()
	req, err := http.NewRequest("GET", strURL, nil)
	if err != nil {
		return info, err
	}
	resp, err := iClient.Do(req)
	if err != nil {
		return info, err
	}
	defer resp.Body.Close()

	// 进度条
	d.getBar(int(resp.ContentLength), strURL)

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return info, err
	}
	defer f.Close()

	buf := make([]byte, 32*1024)
	_, err = io.CopyBuffer(io.MultiWriter(f, d.bar), resp.Body, buf)

	info.SourceUrl = strURL
	info.Path = filename

	return info, err
}

// getBar 设置进度条
func (d *Component) getBar(length int, name string) {
	d.bar = progressbar.NewOptions(
		length,
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(50),
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() {
			fmt.Print("\n")
		}),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionSetDescription("downloading "+name),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
}

// 下载分片
func (d *Component) downloadPartial(strURL, filename string, rangeStart, rangeEnd, i int) {
	if rangeStart >= rangeEnd {
		return
	}

	iClient := d.getGoHttpClient(strURL, "GET").Client()

	req, err := http.NewRequest("GET", strURL, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", rangeStart, rangeEnd))
	resp, err := iClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	flags := os.O_CREATE | os.O_WRONLY
	if d.config.Resume {
		flags |= os.O_APPEND
	}

	partFile, err := os.OpenFile(d.getPartFilename(filename, i), flags, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer partFile.Close()

	buf := make([]byte, 32*1024)
	_, err = io.CopyBuffer(io.MultiWriter(partFile, d.bar), resp.Body, buf)
	if err != nil {
		if err == io.EOF {
			return
		}
		log.Fatal(err)
	}
}

func (d *Component) merge(filename string) error {
	destFile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer destFile.Close()

	for i := 0; i < d.config.Concurrency; i++ {
		partFileName := d.getPartFilename(filename, i)
		partFile, err := os.Open(partFileName)
		if err != nil {
			return err
		}
		io.Copy(destFile, partFile)
		partFile.Close()
		os.Remove(partFileName)
	}

	return nil
}

// getPartDir 部分文件存放的目录
func (d *Component) getPartDir(filename string) string {
	return path.Dir(filename)
}

// getPartFilename 构造部分文件的名字
func (d *Component) getPartFilename(filename string, partNum int) string {
	partDir := d.getPartDir(filename)
	return fmt.Sprintf("%s/%s-%d", partDir, path.Base(filename), partNum)
}

func (d *Component) print(topic string, msg string, errtype string) {
	if errtype == "error" {
		d.logger.With(elog.FieldName(topic)).Error(msg)
	} else {
		d.logger.With(elog.FieldName(topic)).Info(msg)
	}
}
