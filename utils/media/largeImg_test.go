package media

import (
	"log"
	"testing"
)

func TestAll(t *testing.T) {
	uris := []string{
		"https://wx3.sinaimg.cn/orj360/8f9c1340ly1hhlner7nzaj22yo4g0kjr.jpg",
		"https://pbs.twimg.com/media/F58o4SjaQAAbCnr?format=jpg&name=medium",
	}

	for _, s := range uris {
		log.Println(s, "==>", GetLargeImg(s))
	}
}
