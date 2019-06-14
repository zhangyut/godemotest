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
	register, err := registry.NewEtcdRegistry("127.0.0.1", []string{"http://127.0.0.1:2379"})
	if err != nil {
		panic("create registry failed:" + err.Error())
	}

	rrmanagerProxy := haproxy.GetRestProxy(register, "rrmanager", rrmanager.SupportedResources())

	zview := &rrmanager.ZdnsuserView{
		Id:       "64b8c0ab40402150802b4577d8bc0606",
		Cn:       "test",
		Priority: 1,
		View:     "zdns10d3809540692b9a80cd1177493a4896",
	}

	viewTask := quark.NewTask()
	viewTask.User = usermanager.Admin
	viewTask.AddCmd(&rest.PutCmd{NewResource: zview})
	var errInfo string
	rrmanagerProxy.HandleTask(viewTask, nil, &errInfo)
	if errInfo != "" {
		panic(errInfo)
	}
	<-time.After(2 * time.Second)
}
