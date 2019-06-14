package main

import (
	"fmt"
	"sort"
)

func main() {
	a := sort.StringSlice{}
	a = append(a, "abc123")
	a = append(a, "abc456")
	a = append(a, "abc789")
	a.Sort()
	fmt.Println(a.Search("abc123"))
	fmt.Println(a.Search("abc456"))
	fmt.Println(a.Search("abc789"))
	fmt.Println(a.Search(""))
}
