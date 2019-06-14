package main

import (
	"fmt"
	"quark"
	"quark/haproxy"
	"quark/registry"
	"quark/rest"
	"time"
	"zcloud-go/privatezone"
)

func main() {
	register, err := registry.NewEtcdRegistry("127.0.0.1", []string{"http://127.0.0.1:2379"})
	if err != nil {
		panic("create registry failed:" + err.Error())
	}

	rrmanagerProxy := haproxy.GetRestProxy(register, "privatezone", privatezone.SupportedResources())
	//rrmanagerProxy := haproxy.GetRestProxy(register, "rrmanager", nil)
	rrTask := quark.NewTask()
	rrTask.User = "admin"
	rrTask.AddCmd(&rest.GetCmd{
		ResourceType: "private_rr",
		Conds:        map[string]interface{}{"zone": "zhanglei.test."},
	})
	var ret string
	rrs := []privatezone.PrivateRr{}
	err = rrmanagerProxy.HandleTask(rrTask, &rrs, &ret)
	if ret != "" {
		fmt.Println("get rr failed: " + ret)
	}
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(rrs)
	fmt.Println("=================")
	<-time.After(2 * time.Second)
}
