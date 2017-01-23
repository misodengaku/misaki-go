package main

import (
	"fmt"

	"./emoji"
	"./misaki"
)

func main() {
	m, _ := misaki.New("misaki_gothic.png", "misaki_4x8_jisx0201.png")
	// e := emoji.New(":innocent:", ":rabbit:", 128)
	e := emoji.New(":fusai:", ":t_r_s:", 128)

	// 7x7, 3x7で切り出す（一部の記号以外に影響なし）
	m.Use7x7Font = false

	querystring := "こんにちは～パルフーズでーす ｱｲｳｴｵ１②Ⅲ　秒速5㌢㍍"

	// querystring := "まみむめもマミムメモ円園"

	// ShiftJISにたぶん含まれてないやつ（エラーになる）
	// imgs, err := m.ConvertString("䨺", true)

	// JIS X0213（エラーになる）
	// imgs, err := m.ConvertString("倻", true)

	// JIS第二水準
	// imgs, err := m.ConvertString("椚", true)
	_, err := m.ConvertString(querystring, true)

	// 適当にヤバそうな文字を含む文字列
	// imgs, err := m.ConvertString("こんにちは～パルフーズでーす ｱｲｳｴｵ１②Ⅲ　秒速5㌢㍍", true)
	if err != nil {
		panic(err)
	}
	// fmt.Printf("%#v\r\n", imgs)
	// e.AAFromImageFile("ア.png")
	// e.AAFromImages(imgs)

	for _, v := range e.ConvertToVerticalAA(e.AAFromImageFileString(querystring)) {
		fmt.Println(v)
	}
	fmt.Println("done")
}
