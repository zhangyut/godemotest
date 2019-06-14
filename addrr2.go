package main

import (
	"time"

	"quark"
	"quark/haproxy"
	"quark/registry"
	"quark/rest"
	"zcloud-go/rrmanager"
)

func main() {
	register, err := registry.NewEtcdRegistry("192.168.142.132", []string{"http://192.168.249.12:2379"})
	if err != nil {
		panic("create registry failed:" + err.Error())
	}

	rrmanagerProxy := haproxy.GetRestProxy(register, "rrmanager", rrmanager.SupportedResources())

	rrTask := quark.NewTask()
	rrTask.User = "zhanglei"
	rr0 := &rrmanager.Rr{
		Id:    "a.zhanglei.test.id",
		Zone:  "zhanglei.test.",
		Name:  "a.zhanglei.test.",
		View:  "others",
		Type:  "cname",
		Rdata: "www.baidu.com.",
		Ttl:   3600,
	}
	rrTask.AddCmd(&rest.PostCmd{NewResource: rr0})
	var ret string
	rrmanagerProxy.HandleTask(rrTask, nil, &ret)
	if ret != "" {
		panic(ret)
	}
	<-time.After(2 * time.Second)
}
