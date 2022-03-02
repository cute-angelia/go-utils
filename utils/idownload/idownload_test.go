package idownload

import (
	"github.com/cute-angelia/go-utils/syntax/ifile"
	"github.com/cute-angelia/go-utils/syntax/ijson"
	"log"
	"testing"
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
			Uri:    "https://best-ecology.oss-cn-hangzhou.aliyuncs.com/common/Clash.for.Windows-0.19.7-win.7z",
			Socket: "",
		},
	)
	for _, datum := range data {
		idown := Load("").Build()
		if len(datum.Socket) > 0 {
			idown = Load("").Build(WithDebug(true), WithProxySocks5(datum.Socket))
		}

		newName := ifile.NewFileName(datum.Uri).GetNameOrigin()
		fileinfo, err := idown.DownloadWithProgressbar(datum.Uri, "/tmp/"+newName)

		log.Println(ijson.Pretty(fileinfo))
		log.Println(err)

		//if filebyte, _, err := idown.RequestFile(datum.Uri); err != nil {
		//	log.Println("获取图片失败：❌", err)
		//} else {
		//	log.Println(len(filebyte))
		//}
	}
}

func TestDownloadHeader(t *testing.T) {
	fileuri := "http://ali2.a.kwimgs.com/ufile/atlas/NTIwMzM0NjQ3NjI4NzQyODkwM18xNjMyMjMyOTczMzMx_0.jpg"
	idown := Load("").Build(
		WithDebug(true),
		WithUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36"),
	)
	if filebyte, key, err := idown.RequestFile(fileuri); err != nil {
		log.Println("获取图片失败：❌", err)
	} else {
		log.Println(len(filebyte), key)
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
	if filebyte, key, err := idown.RequestFile(fileuri); err != nil {
		log.Println("获取图片失败：❌", err)
	} else {
		log.Println("获取图片成功", len(filebyte), key)
	}
}
