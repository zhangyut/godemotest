package main

import (
	"fmt"
	"quark"
	"quark/haproxy"
	"quark/registry"
	"time"
	"zcloud-go/prune"
)

func main() {
	register, err := registry.NewEtcdRegistry("127.0.0.1", []string{"http://127.0.0.1:2379"})
	if err != nil {
		panic("create registry failed:" + err.Error())
	}

	proxy := haproxy.GetHttpCmdProxy(register, "prune_cmd", prune.SupportedCmds())

	task := quark.NewTask()
	task.User = "admin"
	am := &prune.GetRrs{
		Zone: "zhanglei.test.",
		Name: []string{"www.zhanglei.test.", "aaa.zhanglei.test."},
		View: []string{"others"},
	}
	task.AddCmd(am)

	var ret string
	rr := []prune.Rr{}
	err = proxy.HandleTask(task, &rr, &ret)
	if ret != "" {
		fmt.Println("get rr failed: " + ret)
	}
	if err != nil {
		fmt.Println("get rr error: " + err.Error())
	}
	fmt.Println(rr)
	<-time.After(2 * time.Second)
}
