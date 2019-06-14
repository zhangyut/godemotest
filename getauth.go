package main

import (
	"fmt"

	"flag"
	"quark"
	"quark/haproxy"
	"quark/registry"
	"quark/rest"
	"strings"
	"time"
	"zcloud-go/authorize/resource"
)

var (
	etcd  string
	local string
)

func init() {
	flag.StringVar(&etcd, "e", "", "etcd address")
	flag.StringVar(&local, "i", "", "local address")
}

func main() {
	flag.Parse()
	if etcd == "" || local == "" {
		fmt.Println("parameter error.")
		return
	}
	register, err := registry.NewEtcdRegistry(local, strings.Split(etcd, ","))
	if err != nil {
		panic("create registry failed:" + err.Error())
	}

	proxy := haproxy.GetRestProxy(register, "authorize", resource.SupportedResources())

	task := quark.NewTask()
	task.User = "admin"
	am := []resource.AuthMessage{}
	task.AddCmd(&rest.GetCmd{ResourceType: "auth_message"})

	var ret string
	err = proxy.HandleTask(task, &am, &ret)
	if ret != "" {
		fmt.Println("add authorize failed: " + ret)
	}
	if err != nil {
		fmt.Println("add authorize error: " + err.Error())
	}
	fmt.Println(am[0])
	<-time.After(2 * time.Second)
}
