package main

import (
	"fmt"
	"strings"
)

func main() {
	zone := "zhanglei.test."
	domain1 := "www.zhanglei.test."
	domain2 := "zhanglei.test."

	d1 := domain1[:(len(domain1) - len(zone))]
	d2 := domain1[:(len(domain2) - len(zone))]

	if d1 != "" {
		fmt.Println(d1[:len(d1)-1])
	} else {
		fmt.Println("@")
	}

	if d2 != "" {
		fmt.Println(d2[:len(d2)-1])
	} else {
		fmt.Println("@")
	}
	tmp := strings.Split(domain1, zone)
	fmt.Println(tmp)
	fmt.Println(len(tmp))
	fmt.Println(tmp[0][0 : len(tmp[0])-1])
	fmt.Println(tmp[1])
	fmt.Println("end")
}
