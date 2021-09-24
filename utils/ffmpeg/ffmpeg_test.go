package ffmpeg

import (
	"testing"
)

func TestFFmpeg(t *testing.T) {

	mp4 := "/Users/vanilla/Downloads/沙丘1.mp4"

	iffmpeg := Load().Build(WithFfmpegPath("/usr/local/bin/ffmpeg"))

	t.Log(iffmpeg.GetVideoLength(mp4))

	err := iffmpeg.GetCutPictureOne(mp4, 5, "/Users/vanilla/Downloads/cut12323.jpg")
	t.Log(err)
	//
	err2 := iffmpeg.GetCutPictureOneRandom(mp4, "/Users/vanilla/Downloads/cut12323_random.jpg")
	t.Log(err2)
	//
	//err3 := iffmpeg.GetCutPictures(mp4, 60, "/Users/vanilla/Downloads")
	//t.Log(err3)
}
