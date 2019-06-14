package main

import (
	"fmt"

	"quark"
	"quark/haproxy"
	"quark/registry"
	"quark/rest"
	"zcloud-go/rrmanager"
)

func main() {
	register, err := registry.NewEtcdRegistry("127.0.0.1", []string{"http://127.0.0.1:2379"})
	if err != nil {
		panic("create registry failed:" + err.Error())
	}

	rrmanagerProxy := haproxy.GetRestProxy(register, "rrmanager", rrmanager.SupportedResources())

	viewTask := quark.NewTask()
	viewTask.User = "admin"
	viewTask.AddCmd(&rest.GetCmd{
		ResourceType: "view",
	})
	var ret string
	views := []*rrmanager.View{}
	err = rrmanagerProxy.HandleTask(viewTask, &views, &ret)
	if ret != "" {
		fmt.Println("get view failed: " + ret)
	}
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(*views[0])
	fmt.Println("=================")
}
