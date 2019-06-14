package main

import (
	"fmt"
	"time"

	"quark"
	"quark/haproxy"
	"quark/registry"
	"quark/rest"
	"zcloud-go/rrmanager"
)

func main() {
	register, err := registry.NewEtcdRegistry("192.168.142.132", []string{"http://202.173.9.11:2379"})
	if err != nil {
		panic("create registry failed:" + err.Error())
	}

	rrmanagerProxy := haproxy.GetRestProxy(register, "rrmanager", rrmanager.SupportedResources())

	viewTask := quark.NewTask()
	viewTask.User = "zhanglei"
	zdnsuserview1 := &rrmanager.ZdnsuserView{
		View: "123aaaaaaaaaaaa",
	}
	zdnsuserview2 := &rrmanager.ZdnsuserView{
		View: "123bbbbbbbbbbbb",
	}
	zdnsuserview3 := &rrmanager.ZdnsuserView{
		View: "123ccccccccccccc",
	}
	zdnsuserview4 := &rrmanager.ZdnsuserView{
		View: "123ddddddddddddd",
	}
	viewTask.AddCmd(&rest.PostCmd{NewResource: zdnsuserview1})
	viewTask.AddCmd(&rest.PostCmd{NewResource: zdnsuserview2})
	viewTask.AddCmd(&rest.PostCmd{NewResource: zdnsuserview3})
	viewTask.AddCmd(&rest.PostCmd{NewResource: zdnsuserview4})
	var ret string
	rrmanagerProxy.HandleTask(viewTask, nil, &ret)
	if ret != "" {
		fmt.Println("add view failed: " + ret)
	}
	<-time.After(2 * time.Second)
}
