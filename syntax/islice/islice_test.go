package islice

import (
	"fmt"
	"github.com/cute-angelia/go-utils/syntax/ijson"
	"log"
	"testing"
)

func TestDelete(t *testing.T) {
	var strs []string

	for i := 0; i < 10; i++ {
		strs = append(strs, fmt.Sprintf("xx_%d", i))
	}
	log.Println(ijson.Pretty(strs))
	strs = RemoveStringWithOrder(strs, 3)
	log.Println(ijson.Pretty(strs))

	strs = RemoveStringWithoutOrder(strs, 3)
	log.Println(ijson.Pretty(strs))
}
