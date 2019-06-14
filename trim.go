package main

import (
	"fmt"
	"strings"
)

func main() {
	a := strings.TrimSuffix("cloud.zdns.cn", "zdns.cn")
	fmt.Println(a)
	b := strings.TrimLeft("www.baidu.com", "123.www.baidu.com")
	fmt.Println(b)
}
