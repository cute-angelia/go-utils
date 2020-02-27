package main

import (
	"github.com/cute-angelia/go-utils/file"
	"log"
)

func main() {
	if err := file.DownloadFileWithSrc("https://www.baidu.com/img/baidu_resultlogo@2.png", "/tmp/2323.png"); err != nil {
		log.Println(err.Error())
	}
}
