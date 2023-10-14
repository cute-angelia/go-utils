package iurl

import (
	"github.com/cute-angelia/go-utils/syntax/ijson"
	"log"
	"net/url"
	"testing"
)

func TestUrl(t *testing.T) {
	uri := "https://twattraction.akamaized.net/www_v1/imgComic/ep_content/1281_36208_1606901412.316.jpg?__token__=exp=1651744127~acl=/*~hmac=3280f2cee55caf511f17387507a40cdd923568ba20029786449cd5d03553922e"

	t.Log(Encode(uri))
	t.Log(EncodeQuery(uri))
	// url_test.go:8: https%3A%2F%2Ftwattraction.akamaized.net%2Fwww_v1%2FimgComic%2Fep_content%2F1281_36208_1606901412.316.jpg%3F__token__%3Dexp%3D1651744127~acl%3D%2F%2A~hmac%3D3280f2cee55caf511f17387507a40cdd923568ba20029786449cd5d03553922e
	// url_test.go:9: https://twattraction.akamaized.net/www_v1/imgComic/ep_content/1281_36208_1606901412.316.jpg?__token__=exp%3D1651744127~acl%3D%2F%2A~hmac%3D3280f2cee55caf511f17387507a40cdd923568ba20029786449cd5d03553922e

	t.Log(Decode(Encode(uri)))

	u, _ := url.Parse(uri)

	t.Log(u.Host)

	t.Log(Decode("http://home.shixinyi.cn:38095/photo-station/douyin/img/stars-Ai%E5%B9%BB%E6%83%B3%E8%80%851667814109703008333.jpeg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=admin%2F20221111%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20221111T070857Z&X-Amz-Expires=3600&X-Amz-SignedHeaders=host&X-Amz-Signature=a7eaefefd77d9c4def9251f688c743df70cb669668264bd1aa113776acf7de46"))

	v, _ := url.Parse("http://home.shixinyi.cn:38095/photo-station/douyin/img/一只天涯涯1667823109304367889.jpeg")

	log.Println(ijson.Pretty(v))

	uri2 := "https://baidu.com/"
	t.Log(CleanUrlWithoutParam(uri2))
	t.Log(GetDomainWithOutSlant(uri2))
}
