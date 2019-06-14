package main

import (
	"fmt"
	"github.com/miekg/dns"
)

func main() {
	name := "ccc_.zhanglei111.test."
	ok := dns.IsFqdn(name)
	fmt.Println(ok)
}
