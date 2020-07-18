package bunt

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestLocker(t *testing.T) {

	if err := InitBuntCache("cache", fmt.Sprintf("/tmp/test_%d.db", time.Now().Unix())); err != nil {
		log.Println(err)
		return
	}

	opt1 := NewLockerOpt(WithLimit(3), WithToday(true))
	opt2 := NewLockerOpt(WithLimit(1), WithToday(false))
	opt3 := NewLockerOpt(WithLimit(35))

	for i := 1; i <= 30; i++ {
		if IsNotLockedInLimit("cache", "test1", time.Hour, opt1) {
			log.Println("opt1 > ", i)
		}
	}

	for i := 1; i <= 30; i++ {
		if IsNotLockedInLimit("cache", "test2", time.Hour, opt2) {
			log.Println("opt2 > ", i)
		}
	}

	for i := 1; i <= 30; i++ {
		if IsNotLockedInLimit("cache", "test3", time.Hour, opt3) {
			log.Println("opt3 > ", i)
		}
	}

}
