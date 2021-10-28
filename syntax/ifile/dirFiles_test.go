package ifile

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"testing"
)

func TestDirFiles(t *testing.T) {
	// rootPath := "/Users/vanilla/Downloads"
	// t.Log(GetAllPaths(rootPath))

	//
	t.Log(Name("/Users/vanilla/Downloads/财务管理学/课件/财务管理学课件精讲四王天娇.pdf"))
	t.Log(Name2("/Users/vanilla/Downloads/财务管理学/课件/财务管理学课件精讲四王天娇.pdf"))
}

// go test -v -run TestMapList
func TestMapList(t *testing.T) {
	// inpath := "/Users/vanilla/Downloads/けんけんぱ(けんけ) fantia 2020.03-2021.08"
	basePath := "/Users/vanilla/Downloads/けんけんぱ(けんけ) fantia 2020.03-2021.08"

	data := GetFileMapList(basePath, map[string][]string{})
	// t.Log(ijson.Pretty(data))
	DoPutImage(basePath, data, 1)
}

// 处理图片 保留多少step级子集
func DoPutImage(basePath string, data map[string][]string, step int) {
	baseStepCount := strings.Count(basePath, "/")
	// 其他测试 处理 子集图片，保留一个子集
	for key, val := range data {
		log.Println("===>",key, len(val))
		if len(val) > 0 {
			for _, v := range val {
				// 保留多少级
				baseStepSonCount := strings.Count(key, "/")
				if baseStepSonCount-step <= baseStepCount {
					continue
				}

				filename := path.Base(v)
				//log.Println("---")
				//log.Println(baseStepCount, baseStepSonCount)
				//log.Println(key)
				//log.Println(v)
				//log.Println(filename)

				// 假如是顶级，不处理
				tempv := strings.Replace(v, basePath, "", -1)
				if len(strings.Split(tempv, "/")) == 2 {
					continue
				}
				// 将子集移动到父级
				nfilename := strings.Replace(filename, "fulibl.net", "", -1)
				nfilename = strings.TrimSpace(nfilename)
				newname := fmt.Sprintf("%s_%s", key, nfilename)
				log.Println("--")
				log.Println(v)
				log.Println(newname)
				Mv(v, newname)
			}
		} else {
			log.Println("empty:", key)
		}
	}

	// 清理空文件夹
	dirs,_,_ := GetAllPaths(basePath)

	for _, dir := range dirs {
		if CheckIsEmptyDir(dir) {
			log.Println("empty:", dir)
			os.Remove(dir)
		}
	}

}

// go test -v -run TestEmptyDir
func TestEmptyDir(t *testing.T) {
	inpath := "/Users/vanilla/Downloads/けんけんぱ(けんけ) fantia 2020.03-2021.08/2020.11"
	t.Log(CheckIsEmptyDir(inpath))
}
