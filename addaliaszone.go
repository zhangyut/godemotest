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
	register, err := registry.NewEtcdRegistry("127.0.0.1", []string{"http://127.0.0.1:2379"})
	if err != nil {
		panic("create registry failed:" + err.Error())
	}

	proxy := haproxy.GetRestProxy(register, "limitedresources", resources.SupportedResources())

	task := quark.NewTask()
	task.User = "admin"
	zone1 := &resources.AliasZone{
		Name:      "zhanglei.test",
		UserAlias: "alias1",
	}
	zone2 := &resources.AliasZone{
		Name:      "zhanglei.test",
		UserAlias: "alias2",
	}
	task.AddCmd(&rest.PostCmd{NewResource: zone1})
	task.AddCmd(&rest.PostCmd{NewResource: zone2})
	var ret string
	err = proxy.HandleTask(task, nil, &ret)
	if ret != "" {
		fmt.Println("add alias zone failed: " + ret)
	}
	if err != nil {
		fmt.Println(err)
	}
	<-time.After(2 * time.Second)
}
