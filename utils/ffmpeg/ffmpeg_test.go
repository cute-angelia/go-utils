package ffmpeg

import (
	"testing"
)

func TestFFmpeg(t *testing.T) {
	iffmpeg := Load().Build(WithFfmpegPath("/usr/local/bin/ffmpeg"))
	err := iffmpeg.GetCutPictureOne("/Users/vanilla/Downloads/luli/陆丽.mp4", 20, "/Users/vanilla/Downloads/cut12323.jpg")
	err2 := iffmpeg.GetCutPictureOneRandom("/Users/vanilla/Downloads/luli/陆丽.mp4", "/Users/vanilla/Downloads/cut12323_random.jpg")
	t.Log(err)
	t.Log(err2)
}
