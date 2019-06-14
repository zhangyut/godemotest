package main

import (
	"fmt"
	"net"
)

func main() {
	addrs, err := net.LookupHost("hosttest")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(addrs)
}
