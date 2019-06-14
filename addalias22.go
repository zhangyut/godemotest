package main

import (
	"fmt"
	"time"

	"quark"
	"quark/haproxy"
	"quark/registry"
	"quark/rest"
	"zcloud-go/limitedresources/resources"
)

func main() {
	register, err := registry.NewEtcdRegistry("192.168.79.12", []string{"http://202.173.9.22:2379"})
	if err != nil {
		panic("create registry failed:" + err.Error())
	}

	proxy := haproxy.GetRestProxy(register, "limitedresources", resources.SupportedResources())

	task := quark.NewTask()
	task.User = "admin"
	alias1 := &resources.UserAlias{
		Id:       "alias1",
		Name:     "alias1",
		Password: "alias1p",
	}
	alias2 := &resources.UserAlias{
		Id:       "alias2",
		Name:     "alias2",
		Password: "alias2p",
	}
	task.AddCmd(&rest.PostCmd{NewResource: alias1})
	task.AddCmd(&rest.PostCmd{NewResource: alias2})
	var ret string
	err = proxy.HandleTask(task, nil, &ret)
	if ret != "" {
		fmt.Println("add alias failed: " + ret)
	}
	if err != nil {
		fmt.Println(err)
	}
	<-time.After(2 * time.Second)
}
