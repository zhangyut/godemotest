package main

import (
	"fmt"
	"time"

	"quark"
	"quark/httpcmd"
	"quark/registry"

	"flag"
	"strings"
	"zcloud-go/yundiapi"
)

var (
	domain string
	local  string
	etcd   string
)

func init() {
	flag.StringVar(&domain, "d", "", "domain name")
	flag.StringVar(&local, "l", "127.0.0.1", "local ip")
	flag.StringVar(&etcd, "e", "http://127.0.0.1:2379", "etcd")
}

func main() {
	flag.Parse()
	if domain == "" {
		panic("domain is nil")
	}

	registry, _ := registry.NewEtcdRegistry(local, strings.Split(etcd, ","))
	authproxy, err := httpcmd.GetProxy(registry, "yundiapi_cmd", yundiapi.SupportedCmds())
	if err != nil {
		fmt.Println(err.Error())
	}
	//set login notify
	task := quark.NewTask()
	task.User = "admin"
	task.AddCmd(&yundiapi.Flush{domain})
	errMsg := ""

	var ret int
	err = authproxy.HandleTask(task, &ret, &errMsg)
	if errMsg != "" {
		fmt.Println(errMsg)
		return
	}
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("success.")
	fmt.Println(ret)

	<-time.After(1 * time.Second)
}
