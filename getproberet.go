package main

import (
	"fmt"
	"quark"
	"quark/haproxy"
	"quark/registry"
	"quark/rest"
	"time"
	"zcloud-go/proberetshow"
)

func main() {
	//register, err := registry.NewEtcdRegistry("127.0.0.1", []string{"http://127.0.0.1:2379"})
	register, err := registry.NewEtcdRegistry("192.168.79.31", []string{"http://202.173.9.22:2379"})
	if err != nil {
		panic("create registry failed:" + err.Error())
	}

	proxy := haproxy.GetRestProxy(register, "proberet", nil)
	task := quark.NewTask()
	task.User = "admin"
	task.AddCmd(&rest.GetCmd{
		ResourceType: "probe_ret",
		Conds: map[string]interface{}{
			"name":  "www",
			"zone":  "zhanglei.test.",
			"begin": "2018-03-05T11:18:31+08:00",
			"end":   "2018-03-07T15:18:31+08:00",
			"node":  "beijing",
		},
	})
	var msg string
	ret := []proberetshow.ProbeRet{}
	err = proxy.HandleTask(task, &ret, &msg)
	if msg != "" {
		fmt.Println("get probe ret failed: " + msg)
	}
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(ret)
	fmt.Println("=================")
	<-time.After(2 * time.Second)
}
