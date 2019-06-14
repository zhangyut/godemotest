package main

import (
	"fmt"

	"flag"
	"quark"
	"quark/haproxy"
	"quark/registry"
	"strings"
	"time"
	"zcloud-go/privatezone"
	"zcloud-go/usermanager"
)

var (
	etcd  string
	local string
)

func init() {
	flag.StringVar(&etcd, "e", "http://127.0.0.1:2379", "etcd address")
	flag.StringVar(&local, "i", "127.0.0.1", "local address")
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

	proxy := haproxy.GetHttpCmdProxy(register, "privatezone_cmd", privatezone.SupportedCmds())

	task := quark.NewTask()
	task.User = "admin"
	am := &privatezone.GetZdnsuser{}
	task.AddCmd(am)

	var ret string
	out := []usermanager.Zdnsuser{}
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
