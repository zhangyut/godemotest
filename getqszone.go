package main

import (
	"time"

	"fmt"
	"quark"
	"quark/haproxy"
	"quark/registry"
	"quark/rest"

	"zcloud-go/limitedresources"
)

func main() {
	register, err := registry.NewEtcdRegistry("202.173.9.22", []string{"http://202.173.9.22:2379"})
	if err != nil {
		panic("create registry failed:" + err.Error())
	}

	proxy := haproxy.GetRestProxy(register, "limitedresources", limitedresources.SupportedResources())

	task := quark.NewTask()
	task.User = "admin"
	task.AddCmd(&rest.GetCmd{
		ResourceType: "validity_period",
	})
	var ret string
	res := []limitedresources.ValidityPeriod{}
	err = proxy.HandleTask(task, &res, &ret)
	if ret != "" {
		panic(ret)
	}
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(res)
	<-time.After(2 * time.Second)
}
