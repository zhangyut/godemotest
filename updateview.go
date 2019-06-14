package main

import (
	"time"

	"quark"
	"quark/haproxy"
	"quark/registry"
	"quark/rest"
	"zcloud-go/rrmanager"
	"zcloud-go/usermanager"
)

func main() {
	register, err := registry.NewEtcdRegistry("192.168.142.132", []string{"http://202.173.9.11:2379"})
	if err != nil {
		panic("create registry failed:" + err.Error())
	}

	rrmanagerProxy := haproxy.GetRestProxy(register, "rrmanager", rrmanager.SupportedResources())

	view3 := &rrmanager.View{
		Id:       "123ccccccccccccc",
		Cn:       "我是张雷",
		Priority: 1,
		Childs:   []string{"c_japan"},
	}

	viewTask := quark.NewTask()
	viewTask.User = usermanager.Admin
	viewTask.AddCmd(&rest.PutCmd{NewResource: view3})
	var errInfo string
	rrmanagerProxy.HandleTask(viewTask, nil, &errInfo)
	if errInfo != "" {
		panic(errInfo)
	}
	<-time.After(2 * time.Second)
}
