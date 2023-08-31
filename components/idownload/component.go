package idownload

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	humanize "github.com/dustin/go-humanize"
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
	"strings"
	"sync"
	"time"
)

// 全局 iclient
var iHttpClient *gout.Client

func init() {
	iHttpClient = gout.NewWithOpt(gout.WithInsecureSkipVerify(), gout.WithClose3xxJump())
}

var (
	ErrorNotFound = errors.New("404 file not found")
	ErrorHead     = errors.New("error head")
	ErrorUrl      = errors.New("url 不合法")
)

type FileInfo struct {
	SourceUrl string
	Path      string
}

type Component struct {
	config *config
	// 进度条
	bar *progressbar.ProgressBar
}

// newComponent ...
func newComponent(config *config) *Component {
	comp := &Component{}
	comp.config = config
	return comp
}

func (d *Component) getHttpHeader() gout.H {
	gh := gout.H{}
	if len(d.config.UserAgent) > 0 {
		gh["User-Agent"] = d.config.UserAgent
	}
	if len(d.config.Referer) > 0 {
		gh["Referer"] = d.config.Referer
	}
	if len(d.config.Cookie) > 0 {
		gh["Cookie"] = d.config.Cookie
	}
	if len(d.config.Host) > 0 {
		gh["Host"] = d.config.Host
	}
	if len(d.config.Authorization) > 0 {
		gh["Authorization"] = "Bearer " + d.config.Authorization
	}

	return gh
}

// getGoHttpClient 业务定制头部
func (d *Component) getGoHttpClient(uri string, method string) *dataflow.DataFlow {
	var igout *dataflow.DataFlow
	switch method {
	case "GET":
		igout = iHttpClient.GET(uri)
	case "POST":
		igout = iHttpClient.POST(uri)
	case "PUT":
		igout = iHttpClient.PUT(uri)
	case "DELETE":
		igout = iHttpClient.DELETE(uri)
	case "HEAD":
		igout = iHttpClient.HEAD(uri)
	case "OPTIONS":
		igout = iHttpClient.OPTIONS(uri)
	default:
		igout = iHttpClient.GET(uri)
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

	if d.config.Debug {
		igout = igout.Debug(true)
	}

	return igout.SetHeader(d.getHttpHeader()).Debug(d.config.Debug)
}

func (d *Component) validFileContentLength(strURL string) error {
	if d.config.FileMax != -1 {
		length := d.GetContentLength(strURL)
		if length > d.config.FileMax {
			return errors.New(fmt.Sprintf("链接：%s 未下载，大小：%s, 超过设置大小: %s", strURL, humanize.Bytes(uint64(length)), humanize.Bytes(uint64(d.config.FileMax))))
		}
	}

	return nil
}

// GetContentLength 获取文件长度
func (d *Component) GetContentLength(strURL string) int {
	contentLength := 0
	header := http.Header{}
	var statusCode int
	if err := d.getGoHttpClient(strURL, "HEAD").BindHeader(&header).Code(&statusCode).Do(); err == nil {
		contentLength, _ = strconv.Atoi(header.Get("Content-Length"))
	} else {
		log.Println("err:", err)
	}
	return contentLength
}

// Download 下载文件
func (d *Component) Download(strURL, filename string) (fileInfo FileInfo, errResp error) {
	strURL = strings.TrimSpace(strURL)

	if !strings.Contains(strURL, "http") {
		return fileInfo, errors.New("Url 不合法：" + strURL)
	}

	// debug
	if d.config.Debug {
		log.Println("下载地址：", strURL, "保存地址：", filename)
	}

	// valid
	if err := d.validFileContentLength(strURL); err != nil {
		return fileInfo, err
	}

	if filename == "" {
		filename = path.Base(strURL)
	}
	header := http.Header{}
	var statusCode int

	ctx, cancel := context.WithTimeout(context.Background(), d.config.Timeout)
	defer cancel()

	err := d.getGoHttpClient(strURL, "HEAD").BindHeader(&header).Code(&statusCode).Do()
	if err != nil {
		log.Println("Head", err.Error())
		//return FileInfo{}, errors.New("HEAD 失败：" + strURL + err.Error())
	}

	if statusCode == http.StatusNotFound {
		return FileInfo{}, ErrorNotFound
	}

	// 下载地址切换
	if len(header["Location"]) > 0 {
		strURL = header["Location"][0]
	}

	// 是否分片下载
	if statusCode == http.StatusOK && header.Get("Accept-Ranges") == "bytes" && d.config.Concurrency > 0 {
		contentLength, _ := strconv.Atoi(header.Get("Content-Length"))
		if fileInfo, errResp = d.multiDownload(strURL, filename, contentLength); err != nil {
			//  重试下载
			if d.config.RetryAttempt > 0 {
				NewRetry(d.config.RetryAttempt, d.config.RetryWaitTime).Func(func() error {
					log.Println("NewRetry multiDownload", strURL, filename)
					fileInfo, errResp = d.multiDownload(strURL, filename, contentLength)
					if errResp != nil {
						return ErrRetry
					} else {
						return nil
					}
				}).Do(ctx)
			}
		}
		return fileInfo, errResp
	}

	// 单例下载
	if fileInfo, errResp = d.singleDownload(strURL, filename); errResp != nil {
		log.Println("下载失败：错误：", strURL, errResp)
		//  重试下载
		if d.config.RetryAttempt > 0 {
			NewRetry(d.config.RetryAttempt, d.config.RetryWaitTime).Func(func() error {
				log.Println("NewRetry singleDownload", strURL, filename)
				fileInfo, errResp = d.singleDownload(strURL, filename)
				if errResp != nil {
					log.Println("singleDownload 下载失败：", errResp)
					return ErrRetry
				} else {
					return nil
				}
			}).Do(ctx)
		}
	}
	return fileInfo, errResp
}

// DownloadToByteRetry 请求文件，返回 字节
func (d *Component) DownloadToByteRetry(src string, retry int) ([]byte, error) {
	src = strings.TrimSpace(src)

	// valid
	if err := d.validFileContentLength(src); err != nil {
		return []byte{}, err
	}

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
	strURL = strings.TrimSpace(strURL)

	// valid
	if err := d.validFileContentLength(strURL); err != nil {
		return []byte{}, err
	}

	iClient := d.getGoHttpClient(strURL, "GET").Client()
	req, err := http.NewRequest("GET", strURL, nil)
	if err != nil {
		return nil, err
	}

	// 定制 header
	headers := d.getHttpHeader()
	for key, value := range headers {
		req.Header.Add(key, fmt.Sprintf("%v", value))
	}

	resp, err := iClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if d.config.Progressbar {
		// 进度条
		d.getBar(int(resp.ContentLength), strURL)
	}

	var buf bytes.Buffer
	bufcache := make([]byte, 32*1024)
	// You don't need use "bufio.NewWriter(&b)" to create an io.Writer. &b is an io.Writer itself.
	// bufio.NewWriter(&buf)

	if d.config.Progressbar {
		io.CopyBuffer(io.MultiWriter(&buf, d.bar), resp.Body, bufcache)
	} else {
		io.CopyBuffer(&buf, resp.Body, bufcache)
	}

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

	// defer os.RemoveAll(partDir)

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
				if d.config.Progressbar {
					d.bar.Add(downloaded)
				}
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

// singleDownload 直接下载
func (d *Component) singleDownload(strURL, filename string) (FileInfo, error) {
	var info FileInfo
	// 需要进度条，更要复用 http.client 这里就不使用原生的 http
	// resp, err := http.Get(strURL)
	// Transport: &http.Transport{
	//	MaxIdleConnsPerHost: 10000,
	// },
	//d.getGoHttpClient(strURL, "GET").Do()

	iClient := d.getGoHttpClient(strURL, "GET").Client()

	req, err := http.NewRequest("GET", strURL, nil)
	if err != nil {
		return info, err
	}

	// 定制 header
	headers := d.getHttpHeader()
	for key, value := range headers {
		req.Header.Add(key, fmt.Sprintf("%v", value))
	}

	resp, err := iClient.Do(req)
	if err != nil {
		return info, err
	}
	defer resp.Body.Close()

	if d.config.Progressbar {
		// 进度条
		d.getBar(int(resp.ContentLength), strURL)
	}

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return info, err
	}
	defer f.Close()

	buf := make([]byte, 32*1024)

	if d.config.Progressbar {
		_, err = io.CopyBuffer(io.MultiWriter(f, d.bar), resp.Body, buf)
	} else {
		_, err = io.CopyBuffer(f, resp.Body, buf)
	}

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
		progressbar.OptionUseANSICodes(true),
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
		log.Println(err)
		return
	}

	// 定制 header
	headers := d.getHttpHeader()
	for key, value := range headers {
		req.Header.Add(key, fmt.Sprintf("%v", value))
	}

	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", rangeStart, rangeEnd))
	resp, err := iClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	flags := os.O_CREATE | os.O_WRONLY
	if d.config.Resume {
		flags |= os.O_APPEND
	}

	partFile, err := os.OpenFile(d.getPartFilename(filename, i), flags, 0666)
	if err != nil {
		log.Println(err)
		return
	}
	defer partFile.Close()

	buf := make([]byte, 32*1024)

	if d.config.Progressbar {
		_, err = io.CopyBuffer(io.MultiWriter(partFile, d.bar), resp.Body, buf)
	} else {
		_, err = io.CopyBuffer(partFile, resp.Body, buf)
	}

	if err != nil {
		if err == io.EOF {
			return
		}
		log.Println("分片下载错误", err)
		// 合并文件并删除临时文件
		for iz := 0; iz < d.config.Concurrency; iz++ {
			os.Remove(d.getPartFilename(filename, iz))
		}
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
