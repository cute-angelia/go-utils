package ifile

import (
	"errors"
	"io/ioutil"
	"path"
	"path/filepath"
)



//获取所有文件和文件夹的路径
// @param ext 过滤文件，只获取匹配后缀名的文件，示例：.go
func GetPaths(dirPath string, exts ...string) (dirPaths []string, filePaths []string, err error) {
	// 处理要过滤的后缀名
	var ext string
	if len(exts) > 0 {
		ext = path.Ext(exts[0])
		if ext == "" {
			err = errors.New("ext format incorrect, ext:" + exts[0])
			return
		}
	}

	// 读取文件和文件夹
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return
	}
	for _, file := range files {
		if file.IsDir() {
			dirPaths = append(dirPaths, filepath.Join(dirPath, file.Name()))
		} else {
			if ext != "" && path.Ext(file.Name()) != ext {
				continue
			}
			filePaths = append(filePaths, filepath.Join(dirPath, file.Name()))
		}
	}
	return
}

//获取所有文件和文件夹的路径，包含子文件夹下的文件和文件夹
// @param ext 过滤文件，只获取匹配后缀名的文件，示例：.go
func GetAllPaths(dirPath string, exts ...string) (dirPaths []string, filePaths []string, err error) {
	dirPaths, filePaths, err = GetPaths(dirPath, exts...)
	if err != nil {
		return
	}

	// 读取子文件夹下文件和文件夹
	for _, dirPath2 := range dirPaths {
		dirPaths2, filePaths2, err := GetAllPaths(dirPath2, exts...)
		if err != nil {
			return nil, nil, err
		}
		dirPaths = append(dirPaths, dirPaths2...)
		filePaths = append(filePaths, filePaths2...)
	}
	return
}
