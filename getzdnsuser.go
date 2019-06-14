package main

import (
	"fmt"
	"quark"
	"quark/haproxy"
	"quark/registry"
	"quark/rest"
	"time"
	"zcloud-go/rrmanager"
)

func main() {
	register, err := registry.NewEtcdRegistry("127.0.0.1", []string{"http://127.0.0.1:2379"})
	if err != nil {
		panic("create registry failed:" + err.Error())
	}

	//rrmanagerProxy := haproxy.GetRestProxy(register, "rrmanager", rrmanager.SupportedResources())
	rrmanagerProxy := haproxy.GetRestProxy(register, "rrmanager", nil)
	rrTask := quark.NewTask()
	rrTask.User = "admin"
	rrTask.AddCmd(&rest.GetCmd{
		ResourceType: "rr",
	})
	var ret string
	rrs := []rrmanager.Rr{}
	err = rrmanagerProxy.HandleTask(rrTask, &rrs, &ret)
	if ret != "" {
		fmt.Println("get rr failed: " + ret)
	}
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(rrs)
	fmt.Println("=================")
	<-time.After(2 * time.Second)
}
