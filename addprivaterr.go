package main

import (
	"time"

	"quark"
	"quark/haproxy"
	"quark/registry"
	"quark/rest"
	"zcloud-go/privatezone"
)

func main() {
	// 如果命令行设置了 cpuprofile
	register, err := registry.NewEtcdRegistry("127.0.0.1", []string{"http://127.0.0.1:2379"})
	if err != nil {
		panic("create registry failed:" + err.Error())
	}

	rrmanagerProxy := haproxy.GetRestProxy(register, "privatezone", privatezone.SupportedResources())

	rrTask := quark.NewTask()
	rrTask.User = "admin"
	rr := &privatezone.PrivateRr{
		Id:          "iddddddddddddd.zhanglei2.test.",
		PrivateZone: "zhanglei2.test.",
		Name:        "www.zhanglei2.test.",
		View:        "others",
		Type:        "a",
		Rdata:       "127.0.0.2",
		Ttl:         3600,
	}
	rrTask.AddCmd(&rest.PostCmd{NewResource: rr})
	var ret string
	err = rrmanagerProxy.HandleTask(rrTask, nil, &ret)
	if err != nil {
		panic(err.Error())
	}
	if ret != "" {
		panic(ret)
	}
	<-time.After(2 * time.Second)
}
