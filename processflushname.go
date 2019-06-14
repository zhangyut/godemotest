package main

import (
	"flag"
	"fmt"
	"quark"
	"quark/haproxy"
	"quark/registry"
	"quark/rest"
	"strconv"
	"zcloud-go/yundiapi"
)

var (
	flushNameId string
	status      string
)

func init() {
	flag.StringVar(&flushNameId, "id", "", "flush name id")
	flag.StringVar(&status, "s", "", "flush name status")
}

func main() {
	flag.Parse()
	if flushNameId == "" {
		panic("flush name id is null.")
	}
	if status == "" {
		status = "0"
	}
	register, err := registry.NewEtcdRegistry("127.0.0.1", []string{"http://127.0.0.1:2379"})
	if err != nil {
		panic("create registry failed:" + err.Error())
	}

	proxy := haproxy.GetRestProxy(register, "yundiapi", yundiapi.SupportedResources())

	get := quark.NewTask()
	get.User = "admin"
	get.AddCmd(&rest.GetCmd{
		ResourceType: "flush_name",
		Conds:        map[string]interface{}{"id": flushNameId},
	})
	var ret string
	flushNames := []yundiapi.FlushName{}
	err = proxy.HandleTask(get, &flushNames, &ret)
	if ret != "" {
		fmt.Println(ret)
		return
	}
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	ret = ""

	t := quark.NewTask()
	t.User = "admin"
	flushNames[0].Status, _ = strconv.Atoi(status)
	t.AddCmd(&rest.PutCmd{
		NewResource: &flushNames[0],
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
