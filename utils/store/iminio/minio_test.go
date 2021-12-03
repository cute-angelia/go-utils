package iminio

import (
	"log"
	"path"
	"testing"
)

func TestMinioUpload(t *testing.T) {
	t.Log("good ->")
	iminio := Load("").Build(
		WithEndpoint(" "),
		WithAccesskeyId(" "),
		WithSecretaccessKey(" "),
		WithUseSSL(false),
	)

	fileinput := "/Users/vanilla/gopath/src/github.com/cute-angelia/go-utils/Makefile"
	if info, err := iminio.FPutObject(
		"comic",
		"test.txt",
		fileinput,
		iminio.GetPutObjectOptionByExt(path.Ext(fileinput)),
	); err != nil {
		log.Println(err)
	} else {
		log.Println(info)
	}
}
