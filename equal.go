package main

import (
	"fmt"
)

type Info struct {
	Name string
	Age  int
}

func main() {
	v1 := Info{"zhanglei", 34}
	v2 := Info{"liujianli", 35}
	v3 := Info{"zhanglei", 36}
	v4 := Info{"zhanglei", 34}
	if v1 == v2 {
		fmt.Println("v1 等于 v2")
	}
	if v1 == v3 {
		fmt.Println("v1 等于 v3")
	}
	if v1 == v4 {
		fmt.Println("v1 等于 v4")
	}
}
