package main

import (
	"fmt"
	"unicode"
)

func main() {
	v := "中华人民共和国"
	for _, c := range v {
		fmt.Println(c)
		fmt.Println(unicode.IsLetter(c))
	}
}
