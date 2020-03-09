package file

import (
	"testing"
	"github.com/cute-angelia/go-utils/file"
	"bufio"
	"os"
)

func TestOpenLocalFile(t *testing.T) {
	file.OpenLocalFile("OpenLocalFile.txt")
}

func TestDownloadFileByLocalFile(t *testing.T) {
	files := file.GetFilelist("./download")
	for _, path := range files {
		f, err := os.Open(path)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			lineText := scanner.Text()

			// 设置保存路径
			savepath := file.MakeSavePathWithUrl(false, lineText, "/tmp")

			// 下载文件
			if err := file.DownloadFileWithSrc(lineText, savepath); err != nil {
				t.Fatal(err)
			}

			t.Log(lineText, savepath)
		}
	}
}
