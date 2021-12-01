package ifile

import (
	"fmt"
	"github.com/cute-angelia/go-utils/syntax/istrings"
	"github.com/cute-angelia/go-utils/utils/generator/snowflake"
	"path"
	"strings"
	"time"
)

type fileName struct {
	uri    string
	suffix string
	prefix string
}

func NewFileName(uri string) *fileName {
	return &fileName{
		uri: uri,
	}
}

func (f *fileName) SetSuffix(suffix string) *fileName {
	f.suffix = suffix
	return f
}

func (f *fileName) SetPrefix(prefix string) *fileName {
	f.prefix = prefix
	return f
}

func (f fileName) GetDir() string {
	return path.Dir(f.uri)
}

func (f fileName) CleanUrl() string {
	newName := f.uri
	if strings.Contains(f.uri, "?") {
		newName = strings.Split(newName, "?")[0]
	}
	return newName
}

// name 保持原有名字
func (f fileName) GetNameOrigin() string {
	uri := f.CleanUrl()
	return fmt.Sprintf("%s%s%s%s", f.prefix, NameNoExt(uri), f.suffix, path.Ext(uri))
}

// name 按时间戳
func (f fileName) GetNameTimeline() string {
	uri := f.CleanUrl()
	ext := path.Ext(uri)
	return fmt.Sprintf("%s%d%s%s", f.prefix, time.Now().UnixNano(), f.suffix, ext)
}

// name 按时间戳 反序：minio 文件获取按文件名排序，需要反序时间戳
// 算法：未来时间减去当前时间， 为了防止串号，增加 nano 长度
func (f fileName) GetNameTimelineReverse(withDate bool) string {
	newName := f.uri
	if strings.Contains(newName, "?") {
		newName = strings.Split(newName, "?")[0]
	}
	ext := path.Ext(newName)

	// 3021-01-01 01:01:01
	etime := int64(31529441953)

	randomstr := fmt.Sprintf("%d", time.Now().UnixNano())

	respName := ""
	respName = fmt.Sprintf("%s%d%s", f.prefix, etime-time.Now().Unix(), randomstr[10:len(randomstr)])

	if withDate {
		respName = fmt.Sprintf("%s%s", respName, time.Now().Format("20060102150405"))
	}

	return fmt.Sprintf("%s%s%s", respName, f.suffix, ext)
}

// name 按时间格式
func (f fileName) GetNameTimeDate() string {
	newName := f.uri
	if strings.Contains(newName, "?") {
		newName = strings.Split(newName, "?")[0]
	}
	ext := path.Ext(newName)
	dname := time.Now().Format("20060102-150405") + "-" + istrings.RandomChars(10)
	return fmt.Sprintf("%s%s%s%s", f.prefix, dname, f.suffix, ext)
}

// name 按雪花算法
func (f fileName) GetNameSnowFlow() string {
	newName := ""
	if strings.Contains(f.uri, "?") {
		newName = strings.Split(f.uri, "?")[0]
	}
	ext := path.Ext(newName)
	n, _ := snowflake.NewSnowId(1)
	newName = fmt.Sprintf("%s%s%s%s", f.prefix, n.String(), f.suffix, ext)
	return newName
}
