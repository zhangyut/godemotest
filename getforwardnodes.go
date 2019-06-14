package main

import (
	"fmt"

	"flag"
	"quark"
	"quark/haproxy"
	"quark/registry"
	"strings"
	"time"
	"zcloud-go/forwardzone"
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

	proxy := haproxy.GetHttpCmdProxy(register, "forwardzone_cmd", forwardzone.SupportedCmds())

	task := quark.NewTask()
	task.User = "admin"
	gfn := &forwardzone.GetNodes{}
	task.AddCmd(gfn)

	var ret string
	out := []map[string][]map[string]string{}
	err = proxy.HandleTask(task, &out, &ret)
	if ret != "" {
		fmt.Println("get forward nodes failed: " + ret)
	}
	if err != nil {
		fmt.Println("get forward nodes error: " + err.Error())
	}

	fmt.Println(out)
	<-time.After(2 * time.Second)
}
