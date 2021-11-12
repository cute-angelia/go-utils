package idownload

import (
	"log"
	"testing"
)

func TestDownload(t *testing.T) {
	fileuri := "https://telegra.ph/file/956de9b5ca3c41703eb52.jpg"
	idown := Load("").Build(WithDebug(true), WithProxySocks5("socks5://host-bwg-new.aaqq.in:8096"))
	if filebyte, _, err := idown.RequestFile(fileuri); err != nil {
		log.Println("获取图片失败：❌", err)
	} else {
		log.Println(len(filebyte))
	}
}

func TestDownload2(t *testing.T) {
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
		log.Println("获取图片成功",len(filebyte), key)
	}
}
