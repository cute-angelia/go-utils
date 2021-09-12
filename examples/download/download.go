package main

import (
	"github.com/cute-angelia/go-utils/syntax/ifile"
	"log"
)

func main() {
	if _, err := ifile.DownloadFileWithSrc("https://www.baidu.com/img/baidu_resultlogo@2.png", "/tmp/2323.png", "2.png"); err != nil {
		log.Println(err.Error())
	}
}
