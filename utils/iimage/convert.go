package iimage

import (
	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	"image"
	"image/jpeg"
	"log"
	"os"
	"strings"
)

// img, _, err = image.Decode(file)
// err = jpeg.Encode(f, img, &jpeg.Options{Quality: 100})
// err = png.Encode(f, img)

// ConvertJPG 将图片转换为jpg格式，并返回错误信息
//
// sourceFile 字符串，传入原始图片路径，
// newFile 字符串，传入转换后图片保存路径。
func ConvertJPG(sourceFile, newFile string) error {
	// 打开原文件
	f, err := os.Open(sourceFile)
	// 检查错误
	if err != nil {
		return err
	}
	// 关闭连接
	defer f.Close()

	// 图片解码
	src, _, err := image.Decode(f)
	// 检查错误
	if err != nil {
		return err
	}

	// 获取图片信息
	b := src.Bounds()

	// YCBCr
	img := src.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(image.Rect(0, 0, b.Max.X, b.Max.Y))

	// 新建并打开新图片
	cf, err := os.OpenFile(newFile, os.O_SYNC|os.O_RDWR|os.O_CREATE, 0666)
	// 检查错误
	if err != nil {
		return err
	}
	// 关闭
	defer cf.Close()

	// 图片编码
	return jpeg.Encode(cf, img, &jpeg.Options{Quality: 100})
}

// PosterCover 将指定封面图片进行裁剪，并返回错误信息
//
// scrPhoto 字符串，要裁剪的图片路径，
// newPhoto 字符串，裁剪后的图片保存路径，
func CropJavCover(srcPhoto, newPhoto string, fixWidth int) error {
	img, err := imgio.Open(srcPhoto)
	if err != nil {
		log.Println(err)
		return err
	}

	// 获取图片边界
	// 获取图片宽度并设置为一半
	b := img.Bounds()
	cropWidth := b.Max.X/2 + fixWidth

	return Crop(srcPhoto, newPhoto, cropWidth, 0, b.Max.X, b.Max.Y)
}

// 剪切图片
func Crop(srcFile, newFile string, x, y, w, h int) error {
	img, err := imgio.Open(srcFile)
	if err != nil {
		log.Println(err)
		return err
	}

	b := img.Bounds()
	if w == 0 {
		w = b.Max.X
	}
	if h == 0 {
		h = b.Max.Y
	}

	// crop 图片
	resized := transform.Crop(img, image.Rect(x, y, w, h))

	// log.Println(resized)
	if strings.Contains(srcFile, ".jpg") || strings.Contains(srcFile, ".jpeg") {
		if err := imgio.Save(newFile, resized, imgio.JPEGEncoder(100)); err != nil {
			log.Println(err)
			return err
		} else {
			// log.Println("th---->")
			return nil
		}
	}

	if strings.Contains(srcFile, ".png") {
		if err := imgio.Save(newFile, resized, imgio.PNGEncoder()); err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}
