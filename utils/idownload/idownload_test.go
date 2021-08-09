package idownload

import (
	"log"
	"testing"
)

func TestDownload(t *testing.T) {
	fileuri := "https://telegra.ph/file/956de9b5ca3c41703eb52.jpg"
	idown := Load("").Build(WithDebug(true), WithProxySocks5("socks5://host-bwg-new.aaqq.in:8096"))
	if filebyte, err := idown.RequestFile(fileuri); err != nil {
		log.Println("获取图片失败：❌", err)
	} else {
		log.Println(len(filebyte))
	}
}
