package main

import (
	"time"

	"quark"
	"quark/haproxy"
	"quark/registry"
	"quark/rest"
	"zcloud-go/prune"
)

func main() {
	register, err := registry.NewEtcdRegistry("127.0.0.1", []string{"http://127.0.0.1:2379"})
	if err != nil {
		panic("create registry failed:" + err.Error())
	}

	proxy := haproxy.GetRestProxy(register, "prune", prune.SupportedResources())

	rrTask := quark.NewTask()
	rrTask.User = "admin"
	rr := &prune.Rr{
		Id:    "111ddddddddddddd.zhanglei.test",
		Zone:  "zhanglei.test.",
		Name:  "www.zhanglei.test.",
		View:  "others",
		Type:  "a",
		Rdata: "127.0.0.3",
		Flags: 1,
		Ttl:   3600,
	}
	rrTask.AddCmd(&rest.PostCmd{NewResource: rr})
	var ret string
	err = proxy.HandleTask(rrTask, nil, &ret)
	if ret != "" {
		panic(ret)
	}
	if err != nil {
		panic(err.Error())
	}
	<-time.After(2 * time.Second)
}
