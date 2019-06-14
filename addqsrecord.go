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
	qsRecord := &qssec.QsRecord{
		Name: "ppp.knet.cn.",
		Type: "a",
		View: "others",
		Zone: "knet.cn.",
		Ip:   "11.22.31.4",
		Port: 8090,
	}
	t.AddCmd(&rest.PostCmd{NewResource: qsRecord})
	qsRecord = &qssec.QsRecord{
		Name: "ppp.knet.cn.",
		Type: "txt",
		View: "others",
		Zone: "knet.cn.",
		Ip:   "23.34.45.56",
		Port: 8091,
	}
	t.AddCmd(&rest.PostCmd{NewResource: qsRecord})

	err = qsProxy.HandleTask(t, nil, &ret)
	if ret != "" {
		panic(ret)
	}
	if err != nil {
		panic("add qs record service failed: " + err.Error())
	}
	ret = ""
}
