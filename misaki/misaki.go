package misaki

import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type Misaki struct {
	JISX0208   image.Image
	JISX0201   image.Image
	Use7x7Font bool
}

func New(jisx0208, jisx0201 string) (*Misaki, error) {
	m := Misaki{}
	filex0208, err := os.Open(jisx0208)
	defer filex0208.Close()
	if err != nil {
		return nil, errors.New("JIS X0208 source file decode error")
	}
	m.JISX0208, _, err = image.Decode(filex0208)
	if err != nil {
		return nil, errors.New("JIS X0208 source file decode error")
	}

	filex0201, err := os.Open(jisx0201)
	defer filex0201.Close()
	if err != nil {
		return nil, errors.New("JIS X0201 source file decode error")
	}
	m.JISX0201, _, err = image.Decode(filex0201)
	if err != nil {
		return nil, errors.New("JIS X0201 source file decode error")
	}
	return &m, nil
}

func (m *Misaki) SubimageFromRune(r rune, fileSave bool) (*image.Image, error) {

	char := string(r)

	c, err := utf8ToSjis(char)
	if err != nil {
		return nil, err
	}
	kuRaw := c[0]
	ten := byte(0)
	if kuRaw >= 0x81 && kuRaw <= 0x9f {
		ten = c[1]
		ku := (kuRaw - 0x80) * 2
		if ten < 0x9f {
			ku--
			if ten > 0x7e {
				ten--
			}
			ten = ten - 0x40
		} else {
			ten = ten - 0x9f
		}
		fmt.Printf("%s,\t%02d区,%02x,%02x\r\n", char, ku, kuRaw, ten)

		xSrc := int(ten) * 8
		xDst := (int(ten) + 1) * 8
		ySrc := (int(ku) - 1) * 8
		yDst := int(ku) * 8
		if m.Use7x7Font {
			xDst--
			yDst--
		}
		fmt.Printf("%dx%d - %dx%d\r\n", xSrc, ySrc, xDst, yDst)
		rect := image.Rect(xSrc, ySrc, xDst, yDst)

		fontImg := m.JISX0208.(interface {
			SubImage(r image.Rectangle) image.Image
		}).SubImage(rect)

		if fileSave {
			outputFile, err := os.Create(char + ".png")
			if nil != err {
				fmt.Println(err)
			}
			png.Encode(outputFile, fontImg)
			outputFile.Close()
		}
		return &fontImg, nil

	} else if kuRaw < 0x80 || kuRaw >= 0xa0 {
		kuH := int(kuRaw) >> 4 & 0x0f
		kuL := int(kuRaw) & 0x0f
		ySrc := kuH * 8
		yDst := (kuH + 1) * 8
		xSrc := kuL * 4
		xDst := (kuL + 1) * 4
		if m.Use7x7Font {
			xDst--
			yDst--
		}
		fmt.Printf("%s,\t%02x\r\n", char, kuRaw)
		fmt.Printf("%dx%d - %dx%d\r\n", xSrc, ySrc, xDst, yDst)
		rect := image.Rect(xSrc, ySrc, xDst, yDst)
		fontImg := m.JISX0201.(interface {
			SubImage(r image.Rectangle) image.Image
		}).SubImage(rect)
		if fileSave {
			outputFile, err := os.Create(char + ".png")
			if nil != err {
				fmt.Println(err)
			}
			png.Encode(outputFile, fontImg)
			outputFile.Close()
		}
		return &fontImg, nil
	}
	fmt.Printf("%s,\t%02x\r\n", char, kuRaw)
	return nil, errors.New("out of range")
}

func (m *Misaki) ConvertString(s string, fileSave bool) ([]*image.Image, error) {
	c := make(chan struct{})
	e := make(chan error)
	sLen := utf8.RuneCountInString(s)
	imgs := make([]*image.Image, sLen)
	for i, r := range []rune(s) {
		go func(index int, rc rune) {
			var err error
			imgs[index], err = m.SubimageFromRune(rc, fileSave)
			if err != nil {
				e <- errors.New("convert fail: " + string(rc))
			}
			c <- struct{}{}
		}(i, r)
	}
	for i := 0; i < sLen; i++ {
		select {
		case <-c:
			// OK
		case err := <-e:
			// error
			return nil, err
		}
	}
	return imgs, nil
}

// UTF-8 から ShiftJIS
func utf8ToSjis(str string) (string, error) {
	ret, err := ioutil.ReadAll(transform.NewReader(strings.NewReader(str), japanese.ShiftJIS.NewEncoder()))
	if err != nil {
		return "", err
	}
	return string(ret), err
}
