package idownload

import (
	"context"
	"github.com/cute-angelia/go-utils/syntax/ifile"
	"github.com/cute-angelia/go-utils/syntax/ijson"
	"log"
	"testing"
	"time"
)

type datas struct {
	Uri    string
	Socket string
}

func TestDownloads(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmsgprefix)
	var data []datas
	data = append(data,
		datas{
			Uri:    "https://wx1.sinaimg.cn/large/008oKRrcgy1gv7y8kragcj60u0140goh02.jpg",
			Socket: "",
		},
		datas{
			Uri:    "https://telegra.ph/file/956de9b5ca3c41703eb52.jpg",
			Socket: "socks5://host-bwg-new.aaqq.in:8096",
		},
		datas{
			Uri:    "https://go.dev/dl/go1.17.7.src.tar.gz",
			Socket: "",
		},
	)
	for _, datum := range data {
		idown := Load("").Build()
		if len(datum.Socket) > 0 {
			idown = Load("").Build(WithDebug(true), WithProxySocks5(datum.Socket))
		}

		// 下载文件
		//newName := ifile.NewFileName(datum.Uri).GetNameOrigin()
		//if fileinfo, err := idown.Download(datum.Uri, "/tmp/"+newName); err == nil {
		//	log.Println(ijson.Pretty(fileinfo))
		//} else {
		//	log.Println("获取图片失败：❌", err)
		//}

		if filebyte, err := idown.DownloadToByte(datum.Uri); err != nil {
			log.Println("获取图片失败：❌", err)
		} else {
			log.Println(len(filebyte))
		}
	}
}

func TestDownloadHeader(t *testing.T) {
	fileuri := "http://ali2.a.kwimgs.com/ufile/atlas/NTIwMzM0NjQ3NjI4NzQyODkwM18xNjMyMjMyOTczMzMx_0.jpg"
	idown := Load("").Build(
		WithDebug(true),
		WithUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36"),
	)
	if filebyte, err := idown.DownloadToByte(fileuri); err != nil {
		log.Println("获取图片失败：❌", err)
	} else {
		log.Println(len(filebyte))
	}
}

// go test -v -run TestDownload3
func TestDownload3(t *testing.T) {
	fileuri := "https://cdn.v2ph.com/photos/P0DcKbgkeL39x5Ir.jpg"
	idown := Load("").Build(
		WithDebug(true),
		WithReferer("https://www.v2ph.com/"),
		WithUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36"),
	)
	if filebyte, err := idown.DownloadToByte(fileuri); err != nil {
		log.Println("获取图片失败：❌", err)
	} else {
		log.Println("获取图片成功", len(filebyte))
	}
}

// go test -v -run TestDownLargeFile
func TestDownLargeFile(t *testing.T) {
	fileuri := "https://ttq.jinyemimi.com/2021/0782/202201010831.zip"
	idown := Load("").Build(
		WithDebug(true),
	)
	newname := ifile.NewFileName(fileuri).GetNameSnowFlow()
	if filebyte, err := idown.Download(fileuri, "/tmp/"+newname); err == nil {
		log.Println(ijson.Pretty(filebyte))
	} else {
		log.Println(err)
	}
}

// 定时器测试， 必须 reset
func TestRetryTimer(t *testing.T) {
	timer := time.NewTimer(time.Second * 5)
	for i := 0; i < 12; i++ {
		timer.Reset(time.Second * 2)
		select {
		case <-timer.C:
			log.Println("z...", i)
		}
	}
	timer.Stop()
}

// 重试方法
func TestRetryFuc(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	NewRetry(3, time.Second*3).Func(func() error {
		log.Println("i, ok")
		return ErrRetry
	}).Do(ctx)
}
