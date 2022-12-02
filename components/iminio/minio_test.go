package iminio

import (
	"log"
	"strings"
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

	log.Println(iminio.SignUrlPublic("blog-station/1669965002139224000.jpg"))

	//fileinput := "/Users/vanilla/gopath/src/github.com/cute-angelia/go-utils/Makefile"
	//if info, err := iminio.FPutObject(
	//	"comic",
	//	"test.txt",
	//	fileinput,
	//	iminio.GetPutObjectOptionByExt(path.Ext(fileinput)),
	//); err != nil {
	//	log.Println(err)
	//} else {
	//	log.Println(info)
	//}
}

func TestMinioUploadBase64(t *testing.T) {
	iminio := Load("").Build(
		WithEndpoint(" "),
		WithAccesskeyId(" "),
		WithSecretaccessKey(" "),
		WithUseSSL(false),
	)
	base64str := ``

	b64data := base64str[strings.IndexByte(base64str, ',')+1:]
	if info, err := iminio.PutObjectBase64(
		"photo-station",
		"upload/test.jpg",
		b64data,
		iminio.GetPutObjectOptionByExt(""),
	); err != nil {
		log.Println(err)
	} else {
		log.Println(info)
	}
}
