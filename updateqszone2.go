package main

import (
	"fmt"
	"quark"
	"quark/haproxy"
	"quark/registry"
	"quark/rest"
	"time"
	"zcloud-go/qssec"
)

func main() {
	register, err := registry.NewEtcdRegistry("127.0.0.1", []string{"http://127.0.0.1:2379"})
	if err != nil {
		panic("create registry failed:" + err.Error())
	}

	proxy := haproxy.GetRestProxy(register, "qssec", qssec.SupportedResources())

	get := quark.NewTask()
	get.User = "admin"
	get.AddCmd(&rest.GetCmd{
		ResourceType: "qs_zone",
		Conds:        map[string]interface{}{"scope": "all"},
	})
	var ret string
	qszones := []qssec.QsZone{}
	err = proxy.HandleTask(get, &qszones, &ret)
	if ret != "" {
		fmt.Println(ret)
		return
	}
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	ret = ""

	fmt.Println(qszones[0])
	qszones[0].Expired = time.Now().AddDate(0, 12, 0)

	t := quark.NewTask()
	t.User = "admin"
	t.AddCmd(&rest.PutCmd{
		NewResource: &qszones[0],
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
}
