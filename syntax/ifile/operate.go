package ifile

import (
	"bytes"
	"fmt"
	"github.com/guonaihong/gout"
	"github.com/guonaihong/gout/dataflow"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"
)

/*
0755->即用户具有读/写/执行权限，组用户和其它用户具有读写权限；
0644->即用户具有读写权限，组用户和其它用户具有只读权限；

一般赋予目录0755权限，文件0644权限。
*/

// Mkdir alias of os.MkdirAll()
func Mkdir(dirPath string, perm os.FileMode) error {
	return os.MkdirAll(dirPath, perm)
}

// MkParentDir quick create parent dir
func MkParentDir(fpath string) error {
	dirPath := filepath.Dir(fpath)
	if !IsDir(dirPath) {
		return os.MkdirAll(dirPath, 0775)
	}
	return nil
}

// ************************************************************
//	open files
// ************************************************************

// !!!!! notice: don't forget close file !!!!!

// OpenFile like os.OpenFile, but will auto create dir.
// flag:
//
func OpenFile(filepath string, flag int, perm os.FileMode) (*os.File, error) {
	fileDir := path.Dir(filepath)

	// if err := os.Mkdir(dir, 0775); err != nil {
	if err := os.MkdirAll(fileDir, DefaultDirPerm); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(filepath, flag, perm)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// 打开已经存在的文件， 不存在会新建一个， 返回 *os.File
// open file in read-write mode
// path, os.O_RDWR, 0666) || 0644
func OpenLocalFile(filepath string)(*os.File, error) {
	return  OpenFile(filepath, os.O_RDWR|os.O_CREATE, DefaultFilePerm)
}

// 打开已经存在的文件， 不新建， 返回 *os.File
func OpenLocalFileNoCreate(filepath string) (*os.File, error) {
	return OpenFile(filepath, os.O_RDWR, DefaultFilePerm)
}

// 文件文件-读取本地文件 Local read local file
func GetFileWithLocal(path string) ([]byte, error) {
	imageFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer imageFile.Close()
	return ioutil.ReadAll(imageFile)
}

func CloseFileHandler(f *os.File) {
	f.Close()
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
	}).F().Retry().Attempt(3).WaitTime(time.Second).MaxWaitTime(time.Second * 30).Do()

	if err != nil {
		return nil, err
	}
	return body, nil
}

// 下载文件到磁盘
// dir like /tmp
// Deprecated: utils/idownload 组件 进行下载
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

// ************************************************************
//	remove files
// ************************************************************

// alias methods
var (
	// MustRm  = MustRemove
	QuietRm = DeleteFile
)

func DeleteFile(path string) {
	// delete file
	var err = os.Remove(path)
	if isError(err) {
		return
	}
	fmt.Println("==> done deleting file :" + path)
}

// ************************************************************
//	move files
// ************************************************************

// 移动文件
func Mv(srcPath, destPat string) error {
	return os.Rename(srcPath, destPat)
}

func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}
	return (err != nil)
}
