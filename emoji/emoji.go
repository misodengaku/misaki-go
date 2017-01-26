package emoji

import (
	"fmt"
	"image"
	"os"
	"strings"
)

type EmojiAA struct {
	Background string
	Foreground string
	Threshold  uint8
}

func New(foreground, background string, th uint8) *EmojiAA {
	e := EmojiAA{Foreground: foreground, Background: background, Threshold: th}
	return &e
}

func (e *EmojiAA) ConvertToHorizontalAA(src [][]string) []string {
	var dst []string
	for y := 0; y < 8; y++ {
		line := ""
		for _, moji := range src {
			line += moji[y]
		}
		dst = append(dst, line)
	}
	return dst
}

func (e *EmojiAA) ConvertToVerticalAA(src [][]string) []string {
	var dst []string
	for _, moji := range src {
		for _, v := range moji {
			slice := strings.Split(strings.Trim(v, " \r\n"), " ")
			for i := len(slice) - 1; i > 0; i-- {
				if slice[i] != e.Background {
					break
				}
				slice[i] = ""
			}
			v = strings.Join(slice, " ")
			dst = append(dst, v)
		}
		dst = append(dst, "\n")
	}
	return dst
}

func (e *EmojiAA) AAFromImageFileString(imgFilePrefix string) [][]string {
	imgFiles := make([]string, len([]rune(imgFilePrefix)))
	for i, v := range []rune(imgFilePrefix) {
		imgFiles[i] = string(v) + ".png"
	}
	return e.AAFromImageFiles(imgFiles)
}

func (e *EmojiAA) AAFromImageFiles(imgFiles []string) [][]string {
	ch := make(chan struct{})
	emoji := make([][]string, len(imgFiles))
	for i, file := range imgFiles {
		go func(index int, filename string) {
			emoji[index] = e.AAFromImageFile(filename)
			ch <- struct{}{}
		}(i, file)
	}
	for i := 0; i < len(imgFiles); i++ {
		<-ch
	}
	return emoji
}

func (e *EmojiAA) AAFromImageFile(imgFile string) []string {
	file, err := os.Open(imgFile)
	defer file.Close()
	if err != nil {
		return nil
	}
	img, _, err := image.Decode(file)
	if err != nil {
		return nil
	}
	fmt.Printf("%#v\r\n", img)
	buf := e.AAFromImage(&img)
	// for _, v := range buf {
	// 	fmt.Println(v)

	// }
	return buf
}

func (e *EmojiAA) AAFromImage(img *image.Image) []string {
	xm, ym := getImageSize(img)
	b := getPixelFromGrayScaleImage(img)
	buf := make([]string, ym)

	for y := 0; y < ym; y++ {
		for x := 0; x < xm; x++ {
			d := (b[y][x])
			if d >= e.Threshold {
				buf[y] += e.Background + " "
			} else {
				buf[y] += e.Foreground + " "
			}
		}
	}
	return buf
}

func (e *EmojiAA) AAFromImages(imgs []*image.Image) []string {

	for _, img := range imgs {
		v := e.AAFromImage(img)
		for _, c := range v {
			fmt.Println(c)

		}
	}
	return nil
}

func getImageSize(img *image.Image) (int, int) {
	return (*img).Bounds().Size().X, (*img).Bounds().Size().Y
}

func getPixelFromGrayScaleImage(img *image.Image) [][]uint8 {
	xm, ym := getImageSize(img)
	buf := make([][]uint8, ym)
	fmt.Printf("%#v\r\n", *img)
	// i := 0
	y := 0

	for y = 0; y < ym; y++ {
		buf[y] = make([]uint8, xm)
		for x := 0; x < xm; x++ {

			r, g, b, _ := (*img).At(x, y).RGBA()
			buf[y][x] = uint8(r + g + b)
			// i++
		}
	}
	return buf
}
