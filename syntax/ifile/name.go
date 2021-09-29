package ifile

import (
	"fmt"
	"github.com/cute-angelia/go-utils/syntax/istrings"
	"github.com/cute-angelia/go-utils/utils/generator/snowflake"
	"net/url"
	"path"
	"strings"
	"time"
)

type fileName struct {
	uri string
}

func NewFileName(uri string) *fileName {
	return &fileName{
		uri: uri,
	}
}

// name 保持原有名字
func (f fileName) GetNameOrigin(prefix string) string {
	if z, err := url.Parse(f.uri); err != nil {
		return f.GetNameTimeline(prefix)
	} else {
		ts := strings.Split(z.Path, "/")
		return prefix + ts[len(ts)-1]
	}
}

// name 按时间错
func (f fileName) GetNameTimeline(prefix string) string {
	newName := f.uri
	if strings.Contains(newName, "?") {
		newName = strings.Split(newName, "?")[0]
	}
	ext := path.Ext(newName)

	if len(prefix) == 0 {
		return fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	} else {
		return fmt.Sprintf("%s%d%s", prefix, time.Now().UnixNano(), ext)
	}
}

// name 按时间格式
func (f fileName) GetNameTimeDate(prefix string) string {
	newName := f.uri
	if strings.Contains(newName, "?") {
		newName = strings.Split(newName, "?")[0]
	}
	ext := path.Ext(newName)
	dname := time.Now().Format("20060102-150405") + "-" + istrings.RandomChars(10)
	if len(prefix) == 0 {
		return fmt.Sprintf("%s%s", dname, ext)
	} else {
		return fmt.Sprintf("%s%s%s", prefix, dname, ext)
	}
}

// name 按雪花算法
func (f fileName) GetNameSnowFlow(prefix string) string {
	newName := ""
	if strings.Contains(f.uri, "?") {
		newName = strings.Split(f.uri, "?")[0]
	}
	ext := path.Ext(newName)
	n, _ := snowflake.NewSnowId(1)
	if len(prefix) > 0 {
		newName = fmt.Sprintf("%s%s%s", prefix, n.String(), ext)
	} else {
		newName = fmt.Sprintf("%s%s", n.String(), ext)
	}
	return newName
}
