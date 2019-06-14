package main

import (
	"fmt"

	"flag"
	"quark"
	"quark/haproxy"
	"quark/registry"
	"strings"
	"time"
	"zcloud-go/authorize"
	"zcloud-go/authorize/resource"
)

var (
	etcd  string
	local string
	id    string
)

func init() {
	flag.StringVar(&etcd, "e", "", "etcd address")
	flag.StringVar(&local, "i", "", "local address")
	flag.StringVar(&id, "id", "", "message id")
}

func main() {
	flag.Parse()
	if etcd == "" || local == "" || id == "" {
		fmt.Println("parameter error.")
		return
	}
	register, err := registry.NewEtcdRegistry(local, strings.Split(etcd, ","))
	if err != nil {
		panic("create registry failed:" + err.Error())
	}

	proxy := haproxy.GetHttpCmdProxy(register, "authorize_cmd", authorize.SupportedCmds())

	task := quark.NewTask()
	task.User = "admin"
	am := &authorize.Agree{
		Id: id,
	}
	task.AddCmd(am)

	var ret string
	out := resource.AuthMessage{}
	err = proxy.HandleTask(task, &out, &ret)
	if ret != "" {
		fmt.Println("add authorize failed: " + ret)
	}
	if err != nil {
		fmt.Println("add authorize error: " + err.Error())
	}
	fmt.Println(out)
	<-time.After(2 * time.Second)
}
