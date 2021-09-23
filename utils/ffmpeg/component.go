package ffmpeg

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/cute-angelia/go-utils/syntax/itime"
	"github.com/gotomicro/ego/core/elog"
	"math/rand"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Component struct {
	config *config
	logger *elog.Component
	locker sync.Mutex
}

// newComponent ...
func newComponent(compName string, config *config, logger *elog.Component) *Component {
	return &Component{
		config: config,
		logger: logger,
	}
}

func (c Component) cutOne(sec string, mp4path string, savePic string) error {
	var outputerror string
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(50000)*time.Millisecond)
	cmd := exec.CommandContext(ctx, c.config.FfmpegPath,
		"-loglevel", "error",
		"-y",
		"-ss", sec,
		"-t", "1",
		"-i", mp4path,
		"-vframes", "1",
		savePic)
	defer cancel()
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		outputerror += fmt.Sprintf("lastframecmderr:%v;", err)
		return errors.New(outputerror)
	}
	if stderr.Len() != 0 {
		outputerror += fmt.Sprintf("lastframestderr:%v;", stderr.String())
		return errors.New(outputerror)
	}
	if ctx.Err() != nil {
		outputerror += fmt.Sprintf("lastframectxerr:%v;", ctx.Err())
		return errors.New(outputerror)
	}
	return nil
}

// 截取视频第几秒图片
func (c *Component) GetCutPictureOneRandom(mp4path string, savePic string) error {
	rand.Seed(time.Now().UnixNano())
	videoLen, _ := c.GetVideoLength(mp4path)
	cutsec := rand.Intn(videoLen)
	sec := strconv.Itoa(cutsec)
	return c.cutOne(sec, mp4path, savePic)
}

func (c *Component) GetCutPictureOne(mp4path string, cutsec int, savePic string) error {
	videoLen, _ := c.GetVideoLength(mp4path)
	if cutsec > videoLen {
		return errors.New(fmt.Sprintf("视频长度(%d)小于剪辑秒数(%d)", videoLen, cutsec))
	}
	sec := strconv.Itoa(cutsec)
	return c.cutOne(sec, mp4path, savePic)
}

// 截取视频图片，按x秒
func (c *Component) GetCutPictures(mp4path string, cutsec int, saveDir string) error {
	var outputerror string
	videoLen, _ := c.GetVideoLength(mp4path)
	for i := 0; i < videoLen; i = i + cutsec {
		sec := strconv.Itoa(i)
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(50000)*time.Millisecond)
		cmd := exec.CommandContext(ctx, c.config.FfmpegPath,
			"-loglevel", "error",
			"-y",
			"-ss", sec,
			"-t", "1",
			"-i", mp4path,
			"-vframes", "1",
			saveDir+"/"+strings.Join(itime.SecToStr(int64(i)), "-")+".jpg")
		defer cancel()
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			outputerror += fmt.Sprintf("lastframecmderr:%v;", err)
			return errors.New(outputerror)
		}
		if stderr.Len() != 0 {
			outputerror += fmt.Sprintf("lastframestderr:%v;", stderr.String())
			return errors.New(outputerror)
		}
		if ctx.Err() != nil {
			outputerror += fmt.Sprintf("lastframectxerr:%v;", ctx.Err())
			return errors.New(outputerror)
		}
	}
	return nil
}

// 获取视频时长 s
func (c *Component) GetVideoLength(mp4path string) (int, error) {
	var length int
	for i := 0; i < 2; i++ {
		//视频处理使用，延长超时时间
		ctx, cancel := context.WithTimeout(context.Background(), c.config.Timeout)
		cmd := exec.CommandContext(ctx, c.config.FfmpegPath,
			"-i", mp4path)
		defer cancel()
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		// always return err exit 1; do not catch
		cmd.Run()

		// Duration: 00:22:08.60
		video_length_regexp, _ := regexp.Compile(`Duration: (\d\d:\d\d:\d\d).\d\d`)
		str := video_length_regexp.FindStringSubmatch(stderr.String())
		// log.Println(stderr.String())
		if len(str) > 0 {
			// str = "2006-01-02" + strings.TrimPrefix(str, "Duration:")
			strtime := "2006-01-02 " + str[1]
			if videotime, err := time.Parse("2006-01-02 15:04:05", strtime); err != nil {
				if ctx.Err() != nil {
					c.logger.Info("GenerateLength Err:" + ctx.Err().Error())
				}
				length = 0
			} else {
				length = videotime.Hour()*3600 + videotime.Minute()*60 + videotime.Second()
				break
			}
		}
	}
	//fmt.Println("---------->>>videolength:", length)
	return length, nil
}
