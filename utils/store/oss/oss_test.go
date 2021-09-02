package oss

import (
	"log"
	"testing"
)

func TestPutFile(t *testing.T) {
	// img1 := "https://img1.baidu.com/it/u=4170534835,2356446070&fm=253&fmt=auto&app=120&f=JPEG?w=349&h=364"
	img2 := "https://images.pexels.com/photos/2583852/pexels-photo-2583852.jpeg?auto=compress&cs=tinysrgb&dpr=2&h=750&w=1260"

	log.Println(Load("").MustBuild(
		WithAccessKeyId("LTAILylfdsfaKYg2NGHsd"),
		WithAccessKeySecret("wUJsSs5J00dfsafdsaOKp2vLMsogPiSU8rNubl"),
		WithEndpoint("https://oss-cn-hangzhou.aliyuncs.com"),
		WithBucketName("wallpapedfasr-douydfasin"),
	).PutObjectWithSrc(img2, "test.jpg"))
}
