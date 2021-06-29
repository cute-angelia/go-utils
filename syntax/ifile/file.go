package ifile

import (
	"bytes"
	"fmt"
	"github.com/guonaihong/gout"
	"github.com/guonaihong/gout/dataflow"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
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

func IsExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
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

// 获取文件内容，保存在内存
func GetFileWithSrcWithGout(src string) ([]byte, error) {
	UserAgent := "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1"
	var body []byte
	err := gout.GET(src).SetHeader(gout.H(gout.H{
		"user-agent": UserAgent,
	})).Callback(func(c *dataflow.Context) error {
		switch c.Code {
		case 200:
			c.BindBody(&body)
			return nil
		case 404: //http code为404时，服务端返回是html 字符串
			return fmt.Errorf(src + " 404")
		default:
			return fmt.Errorf(src+" error: %d", c.Code)
		}
	}).Do()

	if err != nil {
		return nil, err
	}
	return body, nil
}

// 下载文件到磁盘
// dir like /tmp
func DownloadFileWithSrc(src string, dir string, filenamewithext string) (string, error) {
	if body, err := GetFileWithSrcWithGout(src); err != nil {
		return "", err
	} else {
		r := bytes.NewReader(body)
		//dir
		// saveDir := path.Dir(dir)
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			panic(err)
		}

		// Create the file
		ifile := fmt.Sprintf("%s/%s", dir, filenamewithext)
		out, err := os.Create(ifile)
		if err != nil {
			return "", err
		}
		defer out.Close()

		// Write the body to file
		if _, err = io.Copy(out, r); err != nil {
			return "", err
		} else {
			return ifile, err
		}
	}
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
