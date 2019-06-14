package main

import (
	"fmt"
	idn "github.com/miekg/dns/idn"
	"strings"
	"unicode"
)

func main() {
	name := "你好.我好"
	fmt.Println(name)
	a := strings.FieldsFunc(name, unicode.IsSpace)
	s := strings.Join(a, " ")

	arr := strings.Split(s, " ")
	if len(arr) != 4 {
		panic("aaaaaaaaaaaa")
	}
	arr[3] = idn.ToPunycode(arr[3])
	fmt.Print(strings.Join(arr, " "))
}
