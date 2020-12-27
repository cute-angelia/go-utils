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
* 图片裁剪
* 入参:
* 规则:(x0, y0)起点坐标，(x1, y1)终点坐标，如果精度quality为0则精度保持不变
*
* 返回:error
 */
func Clip(in io.Reader, out io.Writer, x0, y0, x1, y1, quality int) error {
	origin, fm, err := image.Decode(in)
	if err != nil {
		return err
	}

	width := origin.Bounds().Max.X
	height := origin.Bounds().Max.Y
	if quality == 0 {
		quality = 100
	}
	canvas := resize.Thumbnail(uint(width), uint(height), origin, resize.Lanczos3)

	switch fm {
	case "jpeg":
		img := origin.(*image.YCbCr)
		subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.YCbCr)
		return jpeg.Encode(out, subImg, &jpeg.Options{quality})
	case "png":
		switch canvas.(type) {
		case *image.NRGBA:
			img := canvas.(*image.NRGBA)
			subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.NRGBA)
			return png.Encode(out, subImg)
		case *image.RGBA:
			img := canvas.(*image.RGBA)
			subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.RGBA)
			return png.Encode(out, subImg)
		}
	case "gif":
		img := origin.(*image.Paletted)
		subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.Paletted)
		return gif.Encode(out, subImg, &gif.Options{})
	case "bmp":
		img := origin.(*image.RGBA)
		subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.RGBA)
		return bmp.Encode(out, subImg)
	default:
		return errors.New("Clip IMG ERROR FORMAT")
	}
	return nil
}

/*
* 图片裁剪
* 入参:
* 规则:(x0, y0)起点坐标，(x1, y1)终点坐标，如果精度quality为0则精度保持不变
*
* 返回:error
 */
func ClipPath(inPath string, outPath string, x0, y0, x1, y1, quality int) error {
	in, _ := os.Open(inPath)
	defer in.Close()
	out, _ := os.Create(outPath)
	defer out.Close()

	origin, fm, err := image.Decode(in)
	if err != nil {
		return err
	}

	width := origin.Bounds().Max.X
	height := origin.Bounds().Max.Y
	if quality == 0 {
		quality = 100
	}
	canvas := resize.Thumbnail(uint(width), uint(height), origin, resize.Lanczos3)

	switch fm {
	case "jpeg":
		img := origin.(*image.YCbCr)
		subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.YCbCr)
		return jpeg.Encode(out, subImg, &jpeg.Options{quality})
	case "png":
		switch canvas.(type) {
		case *image.NRGBA:
			img := canvas.(*image.NRGBA)
			subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.NRGBA)
			return png.Encode(out, subImg)
		case *image.RGBA:
			img := canvas.(*image.RGBA)
			subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.RGBA)
			return png.Encode(out, subImg)
		}
	case "gif":
		img := origin.(*image.Paletted)
		subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.Paletted)
		return gif.Encode(out, subImg, &gif.Options{})
	case "bmp":
		img := origin.(*image.RGBA)
		subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.RGBA)
		return bmp.Encode(out, subImg)
	default:
		return errors.New("Clip IMG ERROR FORMAT")
	}
	return nil
}
