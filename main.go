package main

import (
	"fmt"

	"./misaki"
)

func main() {
	m, _ := misaki.New("misaki_gothic.png", "misaki_4x8_jisx0201.png")
	m.ConvertString("あいうえおこんにちは～パルフーズで～す", true)
	fmt.Println("done")
}
