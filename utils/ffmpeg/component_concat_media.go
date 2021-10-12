package ffmpeg

import (
	"errors"
	"fmt"
	"github.com/cute-angelia/go-utils/syntax/icmd"
	"github.com/cute-angelia/go-utils/syntax/ifile"
	"log"
	"os"
	"strings"
)

const (
	TempConcatTxt = "ConcatMovFiles.txt"
)

func (c Component) getTempText() string {
	return fmt.Sprintf("%s/%s", ifile.GetHomeDir(), TempConcatTxt)
}

func (c Component) generateText(dirPath string, ext []string) (tempText string, err error) {
	log.Println(dirPath)
	// 获取文件夹
	if _, files, err := ifile.GetPaths(dirPath, ext...); err != nil {
		log.Println("get path", err)
		return tempText, err
	} else {
		// 生成一个文件用于合并： 格式
		// file ./name.mov
		// log.Println(c.getTempText())
		itempText := ifile.OpenLocalFile(c.getTempText())
		defer itempText.Close()
		for _, file := range files {
			itempText.WriteString(fmt.Sprintf("file %s\n", file))
		}
		return c.getTempText(), nil
	}
}

/*
合并 MOV
ffmpeg -safe 0 -f concat -i files_to_combine -vcodec copy -acodec copy merged.MOV
*/
func (c Component) ConcatMovFiles(dirPath string, ext []string, savePath string) error {
	if temptext, err := c.generateText(dirPath, ext); err != nil {
		log.Println(err)
		return err
	} else {
		status := icmd.Exec(c.config.FfmpegPath, []string{
			"-loglevel", "error",
			"-safe", "0",
			"-f", "concat",
			"-i", temptext,
			"-vcodec", "copy",
			"-acodec", "copy",
			savePath,
		}, c.config.Timeout)

		//log.Println(ijson.Pretty(status))

		if len(status.Stderr) > 0 {
			err := errors.New(strings.Join(status.Stderr, ""))
			return err
		} else {
			os.Remove(temptext)
			return nil
		}
	}
}
