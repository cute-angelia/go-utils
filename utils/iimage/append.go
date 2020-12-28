package iimage

/*
	合并图片， 从底部追加
*/

import (
	"errors"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"os"
)

// outtype "jpg png"
func AppendSingle(imgpath1 string, imgpath2 string, outpath string, outtype string) error {
	img1, _, err := OpenAndDecode(imgpath1)
	if err != nil {
		return err
	}
	img2, _, err := OpenAndDecode(imgpath2)
	if err != nil {
		return err
	}
	// collect pixel data from each image
	pixels1 := DecodePixelsFromImage(img1, 0, 0)
	// the second image has a Y-offset of img1's max Y (appended at bottom)
	pixels2 := DecodePixelsFromImage(img2, 0, img1.Bounds().Max.Y)
	pixelSum := append(pixels1, pixels2...)

	// Set a new size for the new image equal to the max width
	// of bigger image and max height of two images combined
	newRect := image.Rectangle{
		Min: img1.Bounds().Min,
		Max: image.Point{
			X: img2.Bounds().Max.X,
			Y: img2.Bounds().Max.Y + img1.Bounds().Max.Y,
		},
	}
	finImage := image.NewRGBA(newRect)
	// This is the cool part, all you have to do is loop through
	// each Pixel and set the image's color on the go
	for _, px := range pixelSum {
		finImage.Set(
			px.Point.X,
			px.Point.Y,
			px.Color,
		)
	}
	draw.Draw(finImage, finImage.Bounds(), finImage, image.Point{0, 0}, draw.Src)

	// Create a new file and write to it
	// out, err := os.Create("./output.png")
	out, err := os.Create(outpath)
	if err != nil {
		return err
	}

	if outtype == "png" {
		err = png.Encode(out, finImage)
	} else {
		err = jpeg.Encode(out, finImage, &jpeg.Options{100})
	}

	if err != nil {
		return err
	}

	return nil
}

func AppendMulti(multiPaths []string, outpath string, outtype string, deleted bool) error {
	if len(multiPaths) <= 1 {
		return errors.New("合成图片不足2个")
	}

	img1, _, err := OpenAndDecode(multiPaths[0])
	if err != nil {
		return err
	}

	// collect pixel data from each image
	pixels1 := DecodePixelsFromImage(img1, 0, 0)
	pixelSum := []*Pixel{}

	MaxY := img1.Bounds().Max.Y

	// new
	for i := 1; i < len(multiPaths); i++ {
		img2, _, err := OpenAndDecode(multiPaths[i])
		log.Println("追加图片：", multiPaths[i])
		if err != nil {
			panic(err)
		}
		// the second image has a Y-offset of img1's max Y (appended at bottom)
		pixels2 := DecodePixelsFromImage(img2, 0, MaxY)
		pixels1 = append(pixels1, pixels2...)
		pixelSum = pixels1
		MaxY += img2.Bounds().Max.Y
	}

	// Set a new size for the new image equal to the max width
	// of bigger image and max height of two images combined
	newRect := image.Rectangle{
		Min: img1.Bounds().Min,
		Max: image.Point{
			X: img1.Bounds().Max.X,
			Y: MaxY,
		},
	}
	finImage := image.NewRGBA(newRect)
	// This is the cool part, all you have to do is loop through
	// each Pixel and set the image's color on the go
	for _, px := range pixelSum {
		finImage.Set(
			px.Point.X,
			px.Point.Y,
			px.Color,
		)
	}
	draw.Draw(finImage, finImage.Bounds(), finImage, image.Point{0, 0}, draw.Src)

	// Create a new file and write to it
	// out, err := os.Create("./output.png")
	out, err := os.Create(outpath)
	if err != nil {
		return err
	}

	if outtype == "png" {
		err = png.Encode(out, finImage)
	} else {
		err = jpeg.Encode(out, finImage, &jpeg.Options{100})
	}

	if err != nil {
		return err
	} else {
		if deleted {
			for _, i2 := range multiPaths {
				os.RemoveAll(i2)
			}
		}
	}

	return nil
}
