package idownload

/*
	图片帮助类
*/

import (
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

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

// PosterCover 将指定图片进行裁剪，并返回错误信息
//
// scrPhoto 字符串，要裁剪的图片路径，
// newPhoto 字符串，裁剪后的图片保存路径，
func PosterCover(srcPhoto, newPhoto string, fixwidth int) error {
	// 定义各项变量
	var width, height, x int
	// 检查错误
	// 载入图片
	img, errLoad := loadCover(srcPhoto)
	// 检查错误
	if errLoad != nil {
		return errLoad
	}

	// 获取图片边界
	b := img.Bounds()
	// 获取图片宽度并设置为一半
	width = b.Max.X/2 + fixwidth
	// 获取图片高度
	height = b.Max.Y
	// 将x坐标设置为0
	x = width

	// 生成封面图片
	err := clipCover(srcPhoto, newPhoto, x, 0, width, height)

	return err
}

// 剪切图片
func clipCover(srcFile, newFile string, x, y, w, h int) error {
	// 载入图片
	src, err := loadCover(srcFile)
	// 检查错误
	if err != nil {
		return err
	}

	// 剪切图片
	img := src.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(image.Rect(x, y, x+w, y+h))

	// 保存图片
	saveErr := saveCover(newFile, img)
	// 检查错误
	if saveErr != nil {
		return err
	}

	return nil
}

// 载入图片
func loadCover(photo string) (img image.Image, err error) {
	// 打开图片文件
	file, err := os.Open(photo)
	// 检查错误
	if err != nil {
		return
	}
	// 关闭
	defer file.Close()

	// 图片解码
	img, _, err = image.Decode(file)

	return
}

// 保存图片
func saveCover(path string, img image.Image) error {
	// 新建并打开文件
	f, err := os.OpenFile(path, os.O_SYNC|os.O_RDWR|os.O_CREATE, 0666)
	// 检查错误
	if err != nil {
		return err
	}
	// 关闭
	defer f.Close()

	// 获取文件后缀
	ext := filepath.Ext(path)

	// 如果是jpeg类型
	if strings.EqualFold(ext, ".jpg") || strings.EqualFold(ext, ".jpeg") {
		// jpeg图片编码
		err = jpeg.Encode(f, img, &jpeg.Options{Quality: 100})
	} else if strings.EqualFold(ext, ".png") { // png类型
		// png图片编码
		err = png.Encode(f, img)
	}

	return err
}
