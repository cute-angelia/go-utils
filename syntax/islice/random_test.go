package islice

import (
	"log"
	"testing"
)

func TestAll(t *testing.T) {
	log.Println(GetRandomElement([]int{1, 3, 5, 6}))
	log.Println(GetRandomElement([]string{"good", "xxx", "viel"}))

	log.Println(RandomWeightIndex([]int{4150, 4150, 1000, 500, 150, 50}))
}
