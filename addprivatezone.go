package main

import (
	"fmt"
	"time"

	"flag"
	"quark"
	"quark/haproxy"
	"quark/registry"
	"quark/rest"
	"strings"
	"zcloud-go/privatezone"
)

var (
	etcd  string
	local string
	zone  string
)

func init() {
	flag.StringVar(&etcd, "e", "http://127.0.0.1:2379", "etcd")
	flag.StringVar(&local, "i", "127.0.0.1", "local ip")
	flag.StringVar(&zone, "z", "", "zone")
}

func main() {
	flag.Parse()
	if zone == "" {
		panic("zone is nil.")
	}
	register, err := registry.NewEtcdRegistry(local, strings.Split(etcd, ","))
	if err != nil {
		panic("create registry failed:" + err.Error())
	}

	proxy := haproxy.GetRestProxy(register, "privatezone", privatezone.SupportedResources())

	task := quark.NewTask()
	task.User = "admin"
	zone := &privatezone.PrivateZone{
		Id:   zone,
		Name: zone,
	}
	task.AddCmd(&rest.PostCmd{NewResource: zone})
	var ret string
	err = proxy.HandleTask(task, nil, &ret)
	if ret != "" {
		fmt.Println("add private zone failed: " + ret)
	}
	if err != nil {
		fmt.Println(err.Error())
	}
	<-time.After(2 * time.Second)
}
