package main

import (
	"time"

	"fmt"
	"quark"
	"quark/haproxy"
	"quark/registry"
	"quark/rest"

	"flag"
	"zcloud-go/yundiapi"
)

var (
	scope string
)

func init() {
	flag.StringVar(&scope, "s", "", "scope for get")
}

func main() {
	flag.Parse()
	conds := map[string]interface{}{}
	if scope != "" {
		conds["scope"] = scope
	}
	register, err := registry.NewEtcdRegistry("127.0.0.1", []string{"http://127.0.0.1:2379"})
	if err != nil {
		panic("create registry failed:" + err.Error())
	}

	proxy := haproxy.GetRestProxy(register, "yundiapi", yundiapi.SupportedResources())

	task := quark.NewTask()
	task.User = "admin"
	task.AddCmd(&rest.GetCmd{
		ResourceType: "flush_name",
		Conds:        conds,
	})
	var ret string
	res := []yundiapi.FlushName{}
	err = proxy.HandleTask(task, &res, &ret)
	if ret != "" {
		panic(ret)
	}
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(res)
	<-time.After(2 * time.Second)
}
