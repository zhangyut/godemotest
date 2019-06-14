package main

import (
	"fmt"

	"flag"
	"quark"
	"quark/haproxy"
	"quark/registry"
	"quark/rest"
	"strconv"
	"zcloud-go/yundiapi"
)

var (
	blockIpId string
	status    string
)

func init() {
	flag.StringVar(&blockIpId, "id", "", "flush name id")
	flag.StringVar(&status, "s", "", "flush name status")
}

func main() {
	flag.Parse()
	if blockIpId == "" {
		panic("block ip id is null.")
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
		ResourceType: "block_ip",
		Conds:        map[string]interface{}{"id": blockIpId},
	})
	var ret string
	blockIps := []yundiapi.BlockIp{}
	err = proxy.HandleTask(get, &blockIps, &ret)
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
	blockIps[0].Status, _ = strconv.Atoi(status)
	t.AddCmd(&rest.PutCmd{
		NewResource: &blockIps[0],
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
