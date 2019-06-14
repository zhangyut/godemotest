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
	view := &rrmanager.View{
		Id:       "123aaaaaaaaaaaa",
		Cn:       "自定义视图测试",
		Acls:     []string{"c_ct_henan"},
		Priority: 0,
	}
	view2 := &rrmanager.View{
		Id:       "123bbbbbbbbbbbb",
		Cn:       "自定义视图测试1",
		Acls:     []string{"c_cu_hubei"},
		Priority: 10,
	}
	view3 := &rrmanager.View{
		Id:       "123ccccccccccccc",
		Cn:       "自定义视图测试2",
		Acls:     []string{"c_japan"},
		Priority: 20,
	}
	view4 := &rrmanager.View{
		Id:       "123ddddddddddddd",
		Cn:       "自定义视图测试3",
		Acls:     []string{"c_ct_chongqing"},
		Priority: 2,
	}
	viewTask.AddCmd(&rest.PostCmd{NewResource: view})
	viewTask.AddCmd(&rest.PostCmd{NewResource: view2})
	viewTask.AddCmd(&rest.PostCmd{NewResource: view3})
	viewTask.AddCmd(&rest.PostCmd{NewResource: view4})
	var ret string
	err = rrmanagerProxy.HandleTask(viewTask, nil, &ret)
	if ret != "" {
		fmt.Println("add view failed: " + ret)
	}
	<-time.After(2 * time.Second)
}
