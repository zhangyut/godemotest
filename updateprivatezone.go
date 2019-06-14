package main

import (
	"fmt"
	"time"

	"quark"
	"quark/haproxy"
	"quark/registry"
	"quark/rest"
	"zcloud-go/privatezone"
)

func main() {
	register, err := registry.NewEtcdRegistry("127.0.0.1", []string{"http://127.0.0.1:2379"})
	if err != nil {
		panic("create registry failed:" + err.Error())
	}

	proxy := haproxy.GetRestProxy(register, "privatezone", privatezone.SupportedResources())

	get := quark.NewTask()
	get.User = "admin"
	get.AddCmd(&rest.GetCmd{
		ResourceType: "private_zone",
		Conds:        map[string]interface{}{"id": "zhanglei.test."},
	})
	var ret string
	privatezones := []privatezone.PrivateZone{}
	err = proxy.HandleTask(get, &privatezones, &ret)
	if ret != "" {
		return
	}
	if err != nil {
		return
	}

	if len(privatezones) != 1 {
		return

	}
	ret = ""

	t := quark.NewTask()
	privatezones[0].Flags = 2
	t.User = privatezones[0].Zdnsuser
	t.AddCmd(&rest.PutCmd{
		NewResource: &privatezones[0],
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
