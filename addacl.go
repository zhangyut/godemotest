package main

import (
	"fmt"
	"time"

	"flags"
	"quark"
	"quark/haproxy"
	"quark/registry"
	"quark/rest"
	"strings"
	"zcloud-go/rrmanager"
)

var (
	etcd  string
	local string
)

func init() {
	flags.StringVar(&etcd, "e", "", "etcd address")
	flags.StringVar(&local, "i", "", "local address")
}

func main() {
	flags.Parse()
	if etcd == "" || local == "" {
		fmt.Println("parameter error.")
		return
	}
	register, err := registry.NewEtcdRegistry(local, strings.Split(etcd, ","))
	if err != nil {
		panic("create registry failed:" + err.Error())
	}

	rrmanagerProxy := haproxy.GetRestProxy(register, "rrmanager", rrmanager.SupportedResources())

	aclTask := quark.NewTask()
	aclTask.User = "admin"
	acl := &rrmanager.Acl{
		Id:    "acl2222222",
		Name:  "acl222222",
		Cn:    "自定义ACL22222",
		Addrs: []string{"。*$"},
	}
	aclTask.AddCmd(&rest.PostCmd{NewResource: acl})

	var ret string
	err = rrmanagerProxy.HandleTask(aclTask, nil, &ret)
	if ret != "" {
		fmt.Println("add acl failed: " + ret)
	}
	ret = ""

	<-time.After(2 * time.Second)
}
