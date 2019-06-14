package main

import (
	"cement/shell"
	"flag"
	"fmt"
	"strconv"
)

var (
	domain string
	host   string
	ttl    int
	rtype  string
	rdata  string
	op     string
)

func init() {
	flag.StringVar(&domain, "domain", "", "domain")
	flag.IntVar(&ttl, "ttl", 3600, "ttl")
	flag.StringVar(&host, "host", "", "host")
	flag.StringVar(&rtype, "type", "", "type")
	flag.StringVar(&rdata, "rdata", "", "rdata")
	flag.StringVar(&op, "op", "add", "operator")
}

func main() {
	flag.Parse()
	if domain == "" || rtype == "" || rdata == "" {
		fmt.Println("usage: windns [add/del] domain host ttl type rdata")
		return
	}
	if op == "add" {
		ret, err := shell.Shell("dnscmd", "/RecordAdd", domain, host, strconv.Itoa(ttl), rtype, rdata)
		if err != nil {
			fmt.Println(err.Error)
			return
		}
		fmt.Println(ret)
	}
	if op == "del" {
		ret, err := shell.Shell("dnscmd", "/RecordDelete", domain, host, rtype, rdata)
		if err != nil {
			fmt.Println(err.Error)
			return
		}
		fmt.Println(ret)
	}
	return
}
