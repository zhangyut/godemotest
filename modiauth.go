package main

import (
	"fmt"

	"cement/uuid"
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
	id    string
)

func init() {
	flag.StringVar(&etcd, "e", "", "etcd address")
	flag.StringVar(&local, "i", "", "local address")
	flag.StringVar(&id, "id", "", "id")
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

	proxy := haproxy.GetRestProxy(register, "authorize", resource.SupportedResources())

	task := quark.NewTask()
	task.User = "admin"
	am := &resource.AuthMessage{
		Id:             id,
		IdentityImages: []string{"11111111111111111"},
		DomainImages:   []string{"11111111111111111"},
		CompanyUscc:    "11111111111111111",
		CompanyImages:  []string{"11111111111111111"},
		Status:         0,
		Resource:       "zhanglei.test.",
	}
	task.AddCmd(&rest.PostCmd{NewResource: am})

	var ret string
	err = proxy.HandleTask(task, nil, &ret)
	if ret != "" {
		fmt.Println("add authorize failed: " + ret)
	}
	if err != nil {
		fmt.Println("add authorize error: " + err.Error())
	}
	<-time.After(2 * time.Second)
}
