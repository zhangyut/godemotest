package main

import (
	"errors"
	"fmt"
	"net"
)

func checkAddr(addr string) error {
	_, _, err := net.ParseCIDR(addr)
	if err == nil {
		return nil
	}
	ip := net.ParseIP(addr)
	if ip != nil {
		return nil
	}
	return errors.New("addr format error.")
}

func main() {
	a := "129.2.3.4"
	b := "129.2.3.1/24"
	c := "129.2.3"
	err := checkAddr(a)
	if err != nil {
		fmt.Println("111")
		fmt.Println(err.Error())
	}
	err = checkAddr(b)
	if err != nil {
		fmt.Println("222")
		fmt.Println(err.Error())
	}
	err = checkAddr(c)
	if err != nil {
		fmt.Println("333")
		fmt.Println(err.Error())
	}

}
