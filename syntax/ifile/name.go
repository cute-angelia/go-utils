package ifile

import (
	"fmt"
	"github.com/cute-angelia/go-utils/utils/generator/random"
	"github.com/cute-angelia/go-utils/utils/generator/snowflake"
	"path"
	"strings"
	"time"
)

type fileName struct {
	uri    string
	suffix string
	prefix string
	ext    string
}

func NewFileName(uri string) *fileName {
	return &fileName{
		uri: uri,
	}
}

// 后缀
func (f *fileName) SetSuffix(suffix string) *fileName {
	f.suffix = suffix
	return f
}

// 前缀
func (f *fileName) SetPrefix(prefix string) *fileName {
	f.prefix = prefix
	return f
}

// 自定义 ext
func (f *fileName) SetExt(ext string) *fileName {
	f.ext = ext
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

	// ext
	ext := path.Ext(uri)
	if len(f.ext) > 0 {
		ext = f.ext
	}

	return fmt.Sprintf("%s%s%s%s", f.prefix, NameNoExt(uri), f.suffix, ext)
}

// name 按时间戳
func (f fileName) GetNameTimeline() string {
	uri := f.CleanUrl()

	// ext
	ext := path.Ext(uri)
	if len(f.ext) > 0 {
		ext = f.ext
	}

	return fmt.Sprintf("%s%d%s%s", f.prefix, time.Now().UnixNano(), f.suffix, ext)
}

// name 按时间戳 反序：minio 文件获取按文件名排序，需要反序时间戳
// 算法：未来时间减去当前时间， 为了防止串号，增加 nano 长度
func (f fileName) GetNameTimelineReverse(withDate bool) string {
	newName := f.CleanUrl()
	ext := path.Ext(newName)
	if len(f.ext) > 0 {
		ext = f.ext
	}
	/* 自定义算法，未来时间相减
	// 3021-01-01 01:01:01
	etime := int64(3152944195300000000)
	respName := ""
	respName = fmt.Sprintf("%s%d", f.prefix, etime-time.Now().UnixNano())
	if withDate {
		respName = fmt.Sprintf("%s_%s", respName, strings.Replace(time.Now().Format("20060102150405.9999"), ".", "", 1))
	}
	// 加随机数
	// randomstr := fmt.Sprintf("%d", time.Now().UnixNano())
	// respName = fmt.Sprintf("%s_%s", respName, randomstr[10:len(randomstr)])
	*/

	// 雪花算法
	n, _ := snowflake.NewSnowId(1)
	diffTime := int64(991529441953000000) - n.Int64()

	respName := ""
	if withDate {
		respName = fmt.Sprintf("%d_%s", diffTime, time.Now().Format("20060102150405"))
	}

	return fmt.Sprintf("%s%s%s%s", f.prefix, respName, f.suffix, ext)
}

// name 按时间格式
func (f fileName) GetNameTimeDate() string {
	newName := f.CleanUrl()

	// ext
	ext := path.Ext(newName)
	if len(f.ext) > 0 {
		ext = f.ext
	}

	dname := time.Now().Format("20060102-150405") + "-" + random.RandString(10, random.LetterAbcLower)
	return fmt.Sprintf("%s%s%s%s", f.prefix, dname, f.suffix, ext)
}

// name 按雪花算法
func (f fileName) GetNameSnowFlow() string {
	newName := f.CleanUrl()

	// ext
	ext := path.Ext(newName)
	if len(f.ext) > 0 {
		ext = f.ext
	}

	n, _ := snowflake.NewSnowId(1)
	newName = fmt.Sprintf("%s%s%s%s", f.prefix, n.String(), f.suffix, ext)
	return newName
}
