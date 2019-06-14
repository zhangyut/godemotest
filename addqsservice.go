package main

import (
	"flag"
	"fmt"
	"quark"
	"quark/haproxy"
	"quark/registry"
	"quark/rest"
	"strings"
	"zcloud-go/qssec"
)

var (
	localIp string
	etcds   string
)

func init() {
	flag.StringVar(&localIp, "i", "127.0.0.1", "local ip")
	flag.StringVar(&etcds, "e", "http://127.0.0.1:2739", "etcd address")
}

func main() {
	flag.Parse()
	register, err := registry.NewEtcdRegistry(localIp, strings.Split(etcds, ","))
	if err != nil {
		panic("create registry failed:" + err.Error())
	}

	qsProxy := haproxy.GetRestProxy(register, "qssec", qssec.SupportedResources())

	t := quark.NewTask()
	t.User = "admin"
	t.AddCmd(&rest.PostCmd{NewResource: &qssec.QsService{}})
	var ret string
	qsProxy.HandleTask(t, nil, &ret)
	if ret != "" {
		fmt.Println("add qsservice failed: " + ret)
	}
}
