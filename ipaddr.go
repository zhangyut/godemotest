package main

import (
	"fmt"
	"net"
)

func main() {
	ipv4 := "1.1.1.1"
	ipv6 := "2401:8d00:3::15"

	ipv4Addr, err := net.ResolveIPAddr("ip", ipv4)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ipv4Addr.Zone)
	ipv6Addr, err := net.ResolveIPAddr("ip", ipv6)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ipv6Addr.Zone)
}
