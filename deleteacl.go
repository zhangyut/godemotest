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
	register, err := registry.NewEtcdRegistry("192.168.142.132", []string{"http://202.173.9.11:2379"})
	if err != nil {
		panic("create registry failed:" + err.Error())
	}

	rrmanagerProxy := haproxy.GetRestProxy(register, "rrmanager", rrmanager.SupportedResources())

	zoneTask := quark.NewTask()
	zoneTask.User = "zhanglei"
	zoneTask.AddCmd(&rest.DeleteCmd{ResourceType: string(rrmanager.TypeZone), Id: "zhanglei.test."})
	var ret string
	rrmanagerProxy.HandleTask(zoneTask, nil, &ret)
	if ret != "" {
		panic(ret)
	}
	<-time.After(2 * time.Second)
}
