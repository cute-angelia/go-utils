package iminio

import (
	"log"
	"path"
	"strings"
	"testing"
)

func TestGetBucketName(t *testing.T) {
	bucket := ""
	object := ""
	objectNameWithBucket := "http://home.xxxxx.cn:xxxxx/paopao-station/public/image/2024/470501184742042175_20240704175041.jpg"
	if len(objectNameWithBucket) > 0 {

		if strings.Contains(objectNameWithBucket, "http:") || strings.Contains(objectNameWithBucket, "https:") {
			temparray := strings.Split(objectNameWithBucket, "/")
			if len(temparray) >= 3 {
				objectNameWithBucket = strings.Join(temparray[3:], "/")
			}
		}

		objectNameWithBucket = strings.TrimLeft(objectNameWithBucket, "/")
		temp := strings.Split(objectNameWithBucket, "/")
		if len(temp) > 1 {
			objkey := temp[1:len(temp)]

			bucket = temp[0]
			object = strings.Join(objkey, "/")
		}
	}

	log.Println("bucket:", bucket)
	log.Println("object:", object)
}

func TestMinioUpload(t *testing.T) {
	t.Log("good ->")
	iminio := New(
		WithEndpoint("x.x.cn:38103"),
		WithAccesskeyId("x"),
		WithSecretaccessKey("x"),
		WithUseSSL(false),
	)

	bucket := "public"

	//endPoint = "home.shixinyi.cn:38103"
	//accessKeyId = "admin"
	//secretAccessKey = "lovetwins"

	log.Println(iminio.SignUrlPublic("public/th.jpeg"))

	fileinput := "/Users/vanilla/Downloads/th.jpeg"
	if info, err := iminio.FPutObject(
		bucket,
		"th2.jpeg",
		fileinput,
		iminio.GetPutObjectOptionByExt(path.Ext(fileinput)),
	); err != nil {
		log.Println(err)
	} else {
		log.Println(info)
	}

	log.Println(iminio.SignUrlPublic("public/th2.jpeg"))
}

//
//func TestMinioUploadBase64(t *testing.T) {
//	iminio := Load("").Build(
//		WithEndpoint(" "),
//		WithAccesskeyId(" "),
//		WithSecretaccessKey(" "),
//		WithUseSSL(false),
//	)
//	base64str := ``
//
//	b64data := base64str[strings.IndexByte(base64str, ',')+1:]
//	if info, err := iminio.PutObjectBase64(
//		"photo-station",
//		"upload/test.jpg",
//		b64data,
//		iminio.GetPutObjectOptionByExt(""),
//	); err != nil {
//		log.Println(err)
//	} else {
//		log.Println(info)
//	}
//}
