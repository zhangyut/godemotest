package main

import (
	"fmt"
)

func main() {
	var l []byte
	var b [4]byte
	l = b[:]
	fmt.Println(l)
}
