package iimage

import (
	"errors"
	"github.com/nfnt/resize"
	"golang.org/x/image/bmp"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
)

/*
* 缩略图生成
* 入参:
* 规则: 如果width 或 hight其中有一个为0，则大小不变 如果精度quality为0则精度保持不变
* 矩形坐标系起点是左上
* 返回:error
 */
func Thumb(in io.Reader, out io.Writer, width, height, quality int) error {
	origin, typ, err := image.Decode(in)
	if err != nil {
		return err
	}
	if width == 0 || height == 0 {
		width = origin.Bounds().Max.X
		height = origin.Bounds().Max.Y
	}
	if quality == 0 {
		quality = 100
	}
	canvas := resize.Thumbnail(uint(width), uint(height), origin, resize.Lanczos3)

	switch typ {
	case "jpeg":
		return jpeg.Encode(out, canvas, &jpeg.Options{quality})
	case "png":
		return png.Encode(out, canvas)
	case "gif":
		return gif.Encode(out, canvas, &gif.Options{})
	case "bmp":
		return bmp.Encode(out, canvas)
	default:
		return errors.New("Thumb IMG ERROR FORMAT")
	}
}

/*
* 缩略图生成
* 入参:
* 规则: 如果width 或 hight其中有一个为0，则大小不变 如果精度quality为0则精度保持不变
* 矩形坐标系起点是左上
* 返回:error
 */
func ThumbPath(inPath string, outPath string, width, height, quality int) error {
	in, _ := os.Open(inPath)
	defer in.Close()
	out, _ := os.Create(outPath)
	defer out.Close()

	origin, typ, err := image.Decode(in)
	if err != nil {
		return err
	}
	if width == 0 || height == 0 {
		width = origin.Bounds().Max.X
		height = origin.Bounds().Max.Y
	}
	if quality == 0 {
		quality = 100
	}
	canvas := resize.Thumbnail(uint(width), uint(height), origin, resize.Lanczos3)

	switch typ {
	case "jpeg":
		return jpeg.Encode(out, canvas, &jpeg.Options{quality})
	case "png":
		return png.Encode(out, canvas)
	case "gif":
		return gif.Encode(out, canvas, &gif.Options{})
	case "bmp":
		return bmp.Encode(out, canvas)
	default:
		return errors.New("Thumb IMG ERROR FORMAT")
	}
}
