package main

import (
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

	qsProxy := haproxy.GetRestProxy(register, "qssec", qssec.SupportedResources())
	var ret string
	t := quark.NewTask()
	t.User = "admin"
	qsZone := &qssec.QsZone{
		Zone:      "knet.cn.",
		QsService: "58",
	}
	t.AddCmd(&rest.PostCmd{NewResource: qsZone})

	err = qsProxy.HandleTask(t, nil, &ret)
	if ret != "" {
		panic("add qs zone service failed: " + ret)
	}
	if err != nil {
		panic("add qs zone service failed: " + err.Error())
	}
	ret = ""
}
