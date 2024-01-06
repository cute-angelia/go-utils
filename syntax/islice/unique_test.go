package islice

import (
	"log"
	"testing"
)

func TestUnique(t *testing.T) {
	log.Println(Unique([]int{1, 1, 2, 3, 4, 5, 3, 4, 5}))
	log.Println(Unique([]string{"1", "001", "good"}))
}
