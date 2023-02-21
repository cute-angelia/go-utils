package ffmpeg

import (
	"testing"
)

func TestFFmpegCutPicture(t *testing.T) {

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

// go test -v -test.run TestFFmpegConcat
func TestFFmpegConcatMov(t *testing.T) {
	// iffmpeg := Load().Build(WithFfmpegPath("/usr/local/bin/ffmpeg"))
	// err2 := iffmpeg.ConcatMovFiles(`/Users/vanilla/Documents/AnyTrans-Export-2021-10-10/PhotoVideos/1`, []string{}, "/Users/vanilla/Documents/AnyTrans-Export-2021-10-10/PhotoVideos/1/ok.mov")
	// t.Log(err2)
}

// go test -v --run TestFFmpegConcatMp4
func TestFFmpegConcatMp4(t *testing.T) {
	// iffmpeg := Load().Build(WithFfmpegPath("/usr/local/bin/ffmpeg"))
	// err2 := iffmpeg.ConcatMovFiles(`/Users/vanilla/Downloads/test`, []string{".mp4"}, "")
	// t.Log(err2)
}

// go test -v --run TestFFmpegConvertToMp4
func TestFFmpegConvertToMp4(t *testing.T) {
	iffmpeg := Load().Build(WithFfmpegPath("/usr/local/bin/ffmpeg"))
	err2 := iffmpeg.Convert(`/Users/vanilla/Downloads/v1.avi`, "/Users/vanilla/Downloads/v1.mp4")
	t.Log(err2)
}
