package main

import (
	"fmt"
	"quark"
	"quark/haproxy"
	"quark/registry"
	"quark/rest"
	"time"
	"zcloud-go/rrmanager"
)

func main() {
	register, err := registry.NewEtcdRegistry("127.0.0.1", []string{"http://127.0.0.1:2379"})
	if err != nil {
		panic("create registry failed:" + err.Error())
	}

	cmdProxy := haproxy.GetHttpCmdProxy(register, "rrmanager_cmd", nil)
	fmt.Println(cmdProxy)
	rrmanagerProxy := haproxy.GetRestProxy(register, "rrmanager", nil)
	rrTask := quark.NewTask()
	rrTask.User = "admin"
	rrTask.AddCmd(&rest.GetCmd{
		ResourceType: "rr",
		Conds:        map[string]interface{}{"type": "aw", "zone": "zhanglei.test."},
	})
	var ret string
	rrs := []rrmanager.Rr{}
	err = rrmanagerProxy.HandleTask(rrTask, &rrs, &ret)
	if ret != "" {
		fmt.Println("get rr failed: " + ret)
	}
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(rrs)
	fmt.Println("=================")
	<-time.After(2 * time.Second)

	rrs[0].Type = "cnamew"
	rrs[0].Rdata = "10 www.baidu.com."
	rrs[1].Type = "cnamew"
	rrs[1].Rdata = "10 www.google.com."
	rrs[2].Type = "cnamew"
	rrs[2].Rdata = "20 www.163.com."
	rrs[2].Flags = 1
	batchupdaterr := rrmanager.BatchUpdateRr{rrs}
	cmdTask := quark.NewTask()
	cmdTask.User = "admin"
	cmdTask.AddCmd(&batchupdaterr)

	retrrs := []rrmanager.Rr{}
	err = cmdProxy.HandleTask(cmdTask, &retrrs, &ret)
	if ret != "" {
		fmt.Println("update rr failed: " + ret)
	}
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(retrrs)
	<-time.After(2 * time.Second)
}
