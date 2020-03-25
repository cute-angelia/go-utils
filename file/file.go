package file

import (
	"bytes"
	"fmt"
	"github.com/cute-angelia/go-utils/generator/snowflake"
	"github.com/guonaihong/gout"
	"github.com/guonaihong/gout/dataflow"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// 打开已经存在的文件， 不存在会新建一个， 返回 *os.File
// open an existed file or create a file if not exists
// 读写覆盖、0644 其他用户只读
func OpenLocalFile(path string) *os.File {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		f, err := os.Create(path)
		if err != nil {
			panic(err)
		}
		return f
	}
	// open file in read-write mode
	// path, os.O_RDWR, 0666) || 0644
	f, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	return f
}

// 打开已经存在的文件， 不存在会新建一个， 返回 *os.File
// open an existed file or create a file if not exists
// 读写覆盖 || 读写追加、0666 全读写， 0644 其他用户只读
func OpenLocalFileWithFlagPerm(path string, flag int, perm os.FileMode) *os.File {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		f, err := os.Create(path)
		if err != nil {
			panic(err)
		}
		return f
	}
	// open file in read-write mode
	// path, os.O_RDWR|os.O_APPEND, 0666) || 0644
	f, err := os.OpenFile(path, flag, perm)
	if err != nil {
		panic(err)
	}
	return f
}

// 遍历文件夹
func GetFilelist(searchDir string) []string {
	fileList := []string{}
	filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() && f.Name() != ".DS_Store" {
			fileList = append(fileList, path)
		}
		return nil
	})

	return fileList
}

// 文件文件-读取本地文件 Local read local file
func GetFileWithLocal(path string) ([]byte, error) {
	imageFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(imageFile)
}

// 文件-读取网络文件 Net read net file
//func GetFileWithSrc(src string) ([]byte, error) {
//	tr := &http.Transport{
//		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
//	}
//
//	client := &http.Client{Transport: tr, Timeout: time.Second * 6}
//	// set request
//	req, err := http.NewRequest("GET", src, nil)
//	if err != nil {
//		return nil, err
//	}
//	req.Close = true
//	// req.Header = args.Header
//	// get response
//	resp, err := client.Do(req)
//	if err != nil {
//		return nil, err
//	}
//	defer resp.Body.Close()
//	return ioutil.ReadAll(resp.Body)
//}

// 文件-读取网络文件 Net read net file
func GetFileWithSrcWithGout(src string) ([]byte, error) {
	var body []byte
	err := gout.GET(src).Callback(func(c *dataflow.Context) error {
		switch c.Code {
		case 200:
			c.BindBody(&body)
			return nil
		case 404: //http code为404时，服务端返回是html 字符串
			return fmt.Errorf(src + " 404")
		default:
			return fmt.Errorf(src + " error")
		}
	}).Do()

	if err != nil {
		return nil, err
	}
	return body, nil
}

// 下载文件
// DownloadFileWithSrc
// src,
// savePath like /tmp/222.jpg
func DownloadFileWithSrc(src string, dir string, filenamewithext string) error {
	if body, err := GetFileWithSrcWithGout(src); err != nil {
		return err
	} else {
		r := bytes.NewReader(body)

		//dir
		// saveDir := path.Dir(dir)
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			panic(err)
		}

		// Create the file
		filepath := fmt.Sprintf("%s/%s", dir, filenamewithext)
		out, err := os.Create(filepath)
		if err != nil {
			return err
		}
		defer out.Close()

		// Write the body to file
		_, err = io.Copy(out, r)
		return err
	}
}

// 对连接生成名字
// random is use snowflake id
// prefix is like parent dir like "tmp" example: xxxxx.jpg will => tmp/xxxxx.jpg
func MakeNameByUrl(random bool, src string, prefix string) string {
	if strings.Contains(src, "?") {
		src = strings.Split(src, "?")[0]
	}
	filesavepath := ""
	if random {
		ext := path.Ext(src)
		n, _ := snowflake.NewSnowId(1)
		if len(prefix) > 0 {
			filesavepath = fmt.Sprintf("%s_%s%s", prefix, n.String(), ext)
		} else {
			filesavepath = fmt.Sprintf("%s%s", n.String(), ext)
		}
	} else {
		z := path.Base(src)
		if len(prefix) > 0 {
			filesavepath = fmt.Sprintf("%s_%s", prefix, z)
		} else {
			filesavepath = fmt.Sprintf("%s", z)
		}
	}
	return filesavepath
}

// name timeline
func MakeNameByTimeline(src string, prefix string) string {
	if strings.Contains(src, "?") {
		src = strings.Split(src, "?")[0]
	}
	ext := path.Ext(src)
	return fmt.Sprintf("%s/%d%s", prefix, time.Now().UnixNano(), ext)
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

func DeleteFile(path string) {
	// delete file
	var err = os.Remove(path)
	if isError(err) {
		return
	}

	fmt.Println("==> done deleting file")
}

func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}

	return (err != nil)
}
