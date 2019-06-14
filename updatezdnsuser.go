package main

import (
	"fmt"
	"time"

	"quark"
	"quark/haproxy"
	"quark/registry"
	"quark/rest"
	"zcloud-go/usermanager"
)

func main() {
	register, err := registry.NewEtcdRegistry("202.173.9.22", []string{"http://202.173.9.22:2379"})
	if err != nil {
		panic("create registry failed:" + err.Error())
	}

	proxy := haproxy.GetRestProxy(register, "usermanager", usermanager.SupportedResources())

	get := quark.NewTask()
	get.User = "admin"
	get.AddCmd(&rest.GetCmd{
		ResourceType: "zdnsuser",
		Conds:        map[string]interface{}{"name": "zhanglei"},
	})
	var ret string
	zdnsusers := []usermanager.Zdnsuser{}
	err = proxy.HandleTask(get, &zdnsusers, &ret)
	if ret != "" {
		return
	}
	if err != nil {
		return
	}

	if len(zdnsusers) != 1 {
		return

	}
	ret = ""

	zdnsusers[0].Ending = time.Now().AddDate(1, 0, 0)

	t := quark.NewTask()
	t.User = zdnsusers[0].Id
	t.AddCmd(&rest.PutCmd{
		NewResource: &zdnsusers[0],
	})
	err = proxy.HandleTask(t, nil, &ret)
	if ret != "" {
		fmt.Println(ret)
		return
	}
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	<-time.After(2 * time.Second)
}
