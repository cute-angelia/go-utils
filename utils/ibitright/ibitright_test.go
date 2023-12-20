package ibitright

import (
	"log"
	"testing"
)

const (
	p1 Position = iota + 1
	p2
	p3
	p4
	p5
	p31 Position = 31
	p32 Position = 32
	p33 Position = 33
)

func TestSet(t *testing.T) {
	bitrightx := NewBitRight()

	log.Println(bitrightx.Set(p1, 0))
	log.Println(bitrightx.Set(p31, 0))
	log.Println(bitrightx.Set(p32, 0))
	log.Println(bitrightx.Set(p32, 1))
	log.Println(bitrightx.Set(p33, 0))
}

func TestUnSet(t *testing.T) {
	bitrightx := NewBitRight()

	log.Println(bitrightx.UnSet(p1, 0))
	log.Println(bitrightx.UnSet(p2, 0))
	log.Println(bitrightx.UnSet(p31, 1073741824))
	log.Println(bitrightx.UnSet(p32, 2147483648))
	log.Println(bitrightx.UnSet(p32, 2147483649))
	log.Println(bitrightx.UnSet(p33, 2147483648))
}

func TestCheck(t *testing.T) {

	bitrightx := NewBitRight()

	log.Println(bitrightx.Check(p1, 0))
	log.Println(bitrightx.Check(p31, 1073741824))
	log.Println(bitrightx.Check(p32, 2147483648))
	log.Println(bitrightx.Check(p32, 2147483648))
	log.Println(bitrightx.Check(p33, 2147483648))
}
