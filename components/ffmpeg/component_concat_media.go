package ffmpeg

import (
	"errors"
	"fmt"
	"github.com/cute-angelia/go-utils/syntax/icmd"
	"github.com/cute-angelia/go-utils/syntax/ifile"
	"log"
	"os"
	"strings"
	"time"
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
		if itempText, err := ifile.OpenLocalFile(c.getTempText()); err != nil {
			return "", err
		} else {
			defer itempText.Close()
			for _, file := range files {
				itempText.WriteString(fmt.Sprintf("file '%s'\n", file))
			}
			return c.getTempText(), nil
		}
	}
}

/*
合并 MOV
ffmpeg -safe 0 -f concat -i files_to_combine -vcodec copy -acodec copy merged.MOV
*/
func (c Component) ConcatMovFiles(dirPath string, ext []string, saveName string) error {
	if len(saveName) == 0 {
		if len(ext) > 0 {
			saveName = time.Now().Format("20060102-150405") + ext[0]
		} else {
			saveName = time.Now().Format("20060102-150405")
		}
	}
	if temptext, err := c.generateText(dirPath, ext); err != nil {
		log.Println(err)
		return err
	} else {
		if !ifile.IsDir(dirPath + "/success/") {
			ifile.Mkdir(dirPath+"/success/", 0755)
		}

		status := icmd.Exec(c.config.FfmpegPath, []string{
			"-loglevel", "error",
			"-safe", "0",
			"-f", "concat",
			"-i", temptext,
			"-vcodec", "copy",
			"-acodec", "copy",
			dirPath + "/success/" + saveName,
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
