package iurl

import (
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
}
