package main

import (
	"fmt"

	"quark"
	"quark/haproxy"
	"quark/registry"
	"quark/rest"
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
		ResourceType: "qs_record",
		Conds:        map[string]interface{}{"zone": "zdns.cn."},
	})
	var ret string
	qsrecords := []qssec.QsRecord{}
	err = proxy.HandleTask(get, &qsrecords, &ret)
	if ret != "" {
		fmt.Println(ret)
		return
	}
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	ret = ""

	qsrecords[0].StartDdos = false

	t := quark.NewTask()
	t.User = "admin"
	t.AddCmd(&rest.PutCmd{
		NewResource: &qsrecords[0],
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
