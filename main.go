package main

import (
	"fmt"

	"./misaki"
)

func main() {
	m, _ := misaki.New("misaki_gothic.png", "misaki_4x8_jisx0201.png")
	// ShiftJISにたぶん含まれてないやつ（エラーになる）
	// img, err := m.ConvertString("䨺", true)

	// JIS X0213（エラーになる）
	// img, err := m.ConvertString("倻", true)

	// JIS第二水準
	// img, err := m.ConvertString("椚", true)

	// 適当にヤバそうな文字を含む文字列
	img, err := m.ConvertString("こんにちは～パルフーズでーす ｱｲｳｴｵ１②Ⅲ　秒速5㌢㍍", true)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\r\n", img)
	fmt.Println("done")
}
