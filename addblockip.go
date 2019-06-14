package main

import (
	"quark"
	"quark/haproxy"
	"quark/registry"
	"quark/rest"
	"zcloud-go/yundiapi"
)

func main() {
	register, err := registry.NewEtcdRegistry("127.0.0.1", []string{"http://127.0.0.1:2379"})
	if err != nil {
		panic("create registry failed:" + err.Error())
	}

	proxy := haproxy.GetRestProxy(register, "yundiapi", yundiapi.SupportedResources())
	var ret string
	t := quark.NewTask()
	t.User = "admin"
	blockip := &yundiapi.BlockIp{
		Ip: "192.168.12.12/33",
	}
	t.AddCmd(&rest.PostCmd{NewResource: blockip})

	err = proxy.HandleTask(t, nil, &ret)
	if ret != "" {
		panic(ret)
	}
	if err != nil {
		panic("add block ip failed: " + err.Error())
	}
	ret = ""
}
