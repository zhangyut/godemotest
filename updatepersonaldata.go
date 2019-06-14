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
	data := &resource.PersonalData{
		Id:       "admin",
		Zdnsuser: "admin",
		Username: "admin",
		Addrs:    []string{"0.0.0.0/0"},
		Enable:   true,
	}
	t.AddCmd(&rest.PutCmd{NewResource: data})
	var ret string
	proxy.HandleTask(t, nil, &ret)
	if ret != "" {
		panic(ret)
	}
	<-time.After(2 * time.Second)
}
