package file

import (
	"os"
	"path"
	"fmt"
	"github.com/cute-angelia/go-utils/generator/snowflake"
	"runtime"
	"time"
	"strings"
	"io/ioutil"
	"net/http"
	"crypto/tls"
)

// name snowflake
func MakeNewName(newname bool, src string, prefix string) string {
	if strings.Contains(src, "?") {
		src = strings.Split(src, "?")[0]
	}

	filesavepath := ""
	if newname {
		ext := path.Ext(src)
		n, _ := snowflake.NewSnowId(1)

		filesavepath = fmt.Sprintf("%s/%s%s", prefix, n.String(), ext)
	} else {
		z := path.Base(src)
		filesavepath = fmt.Sprintf("%s/%s", prefix, z)
	}

	return filesavepath
}

// name timeline
func MakeNewNameByTimeline(src string, prefix string) string {
	if strings.Contains(src, "?") {
		src = strings.Split(src, "?")[0]
	}

	ext := path.Ext(src)
	return fmt.Sprintf("%s/%d%s", prefix, time.Now().UnixNano(), ext)
}

// 文件-读取网络文件 Net read net file
func GetFileWithSrc(src string) ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr, Timeout: time.Second * 6}
	// set request
	req, err := http.NewRequest("GET", src, nil)
	if err != nil {
		return nil, err
	}
	req.Close = true
	// req.Header = args.Header
	// get response
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// 文件-读取本地文件 Local read local file
func GetFileWithLocal(path string) ([]byte, error) {
	imageFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(imageFile)
}

/**
	创建文件路径
 */
func GeneratePathBySrc(newname bool, to string, src string, prefix string) string {
	// 处理斜杠
	s := to[len(to)-1: len(to)]
	if s != "/" {
		to = to + "/"
	}

	filesavepath := ""

	if newname {
		ext := path.Ext(src)
		n, _ := snowflake.NewSnowId(1)
		filesavepath = fmt.Sprintf("%s%s%d%s", to, prefix, n, ext)
	} else {
		z := path.Base(src)
		filesavepath = fmt.Sprintf("%s%s%s", to, prefix, z)
	}

	return filesavepath
}

// 创建本地文件，方便写入数据，导出
// OpenCreateFile open an existed file or create a file if not exists
func OpenCreateFile(path string, flag int, perm os.FileMode) *os.File {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		f, err := os.Create(path)
		if err != nil {
			panic(err)
		}
		return f
	}
	// open file in read-write mode
	f, err := os.OpenFile(path, flag, perm)
	if err != nil {
		panic(err)
	}
	return f
}

// 获取用户文件夹
func GetUserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	} else if runtime.GOOS == "linux" {
		home := os.Getenv("XDG_CONFIG_HOME")
		if home != "" {
			return home
		}
	}
	return os.Getenv("HOME")
}
