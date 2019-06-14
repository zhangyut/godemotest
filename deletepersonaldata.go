package main

import (
	"time"

	"quark"
	"quark/haproxy"
	"quark/registry"
	"quark/rest"
	"zcloud-go/personalservice/resource"
)

func main() {
	register, err := registry.NewEtcdRegistry("127.0.0.1", []string{"http://127.0.0.1:2379"})
	if err != nil {
		panic("create registry failed:" + err.Error())
	}

	proxy := haproxy.GetRestProxy(register, "personalservice", resource.SupportedResources())

	t := quark.NewTask()
	t.User = "admin"
	t.AddCmd(&rest.DeleteCmd{ResourceType: "personal_data", Id: "admin"})
	var ret string
	err = proxy.HandleTask(t, nil, &ret)
	if err != nil {
		panic(err.Error())

	}
	if ret != "" {
		panic(ret)
	}
	<-time.After(2 * time.Second)
}
