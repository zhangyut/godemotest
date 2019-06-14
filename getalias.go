package main

import (
	"fmt"
	"quark"
	"quark/haproxy"
	"quark/registry"
	"quark/rest"
	"time"
	"zcloud-go/limitedresources/resources"
)

func main() {
	register, err := registry.NewEtcdRegistry("127.0.0.1", []string{"http://127.0.0.1:2379"})
	if err != nil {
		panic("create registry failed:" + err.Error())
	}

	rrmanagerProxy := haproxy.GetRestProxy(register, "limitedresources", nil)
	t := quark.NewTask()
	t.User = "admin"
	t.AddCmd(&rest.GetCmd{
		ResourceType: "user_alias",
	})
	var ret string
	alias := []resources.UserAlias{}
	err = rrmanagerProxy.HandleTask(t, &alias, &ret)
	if ret != "" {
		fmt.Println("get rr failed: " + ret)
	}
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(alias)
	<-time.After(2 * time.Second)
}
