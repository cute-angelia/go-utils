package istrings

import (
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"log"
	"strings"
)

func GbkToUtf8(s string) string {
	reader := transform.NewReader(strings.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		log.Println("GbkToUtf8 error => ", e)
		return ""
	}
	return string(d)
}

func Utf8ToGbk(s string) string {
	reader := transform.NewReader(strings.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		log.Println("GbkToUtf8 error => ", e)
		return ""
	}
	return string(d)
}
