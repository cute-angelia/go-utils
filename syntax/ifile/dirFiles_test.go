package ifile

import (
	"testing"
)

func TestDirFiles(t *testing.T) {
	// rootPath := "/Users/vanilla/Downloads"
	// t.Log(GetAllPaths(rootPath))

	//
	t.Log(Name("/Users/vanilla/Downloads/财务管理学/课件/财务管理学课件精讲四王天娇.pdf"))
	t.Log(Name2("/Users/vanilla/Downloads/财务管理学/课件/财务管理学课件精讲四王天娇.pdf"))
}
