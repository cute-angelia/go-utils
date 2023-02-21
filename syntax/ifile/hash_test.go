package ifile

import (
	"bytes"
	"github.com/guonaihong/gout"
	"log"
	"os"
	"testing"
)

func TestHash(t *testing.T) {
	a := []byte("a1234567688")
	b := []byte("b1234567688")

	t.Log(FileHashSha256(bytes.NewReader(a)))
	t.Log(FileHashSha256(bytes.NewReader(b)))

	t.Log(FileHashMd5(bytes.NewReader(a)))
	t.Log(FileHashMd5(bytes.NewReader(b)))

	// 文件
	file1 := "/Users/vanilla/Downloads/a.jpeg"
	file2 := "/Users/vanilla/Downloads/b.jpeg"

	fileOpen, err := os.Open(file1)
	defer fileOpen.Close()
	if err != nil {
		log.Println(err)
	} else {
		t.Log(FileHashSha256(fileOpen))
	}

	fileOpen2, err2 := os.Open(file2)
	defer fileOpen2.Close()
	if err2 != nil {
		log.Println(err2)
	} else {
		t.Log(FileHashSha256(fileOpen2))

		t.Log(FileHashSHA1(fileOpen2))
		t.Log(len(FileHashSHA1(fileOpen2)))
	}

	t.Log("test img")

	var imgbyte []byte
	img := "https://img1.baidu.com/it/u=4170534835,2356446070&fm=253&fmt=auto&app=120&f=JPEG?w=349&h=364"
	gout.GET(img).BindBody(&imgbyte).Do()

	t.Log(len(imgbyte), FileHashSHA1(bytes.NewReader(imgbyte)))
}
