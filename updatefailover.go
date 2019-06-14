package main

import (
	"fmt"
	"time"

	"quark"
	"quark/haproxy"
	"quark/registry"
	"quark/rest"
	"zcloud-go/failover"
)

func main() {
	register, err := registry.NewEtcdRegistry("127.0.0.1", []string{"http://127.0.0.1:2379"})
	if err != nil {
		panic("create registry failed:" + err.Error())
	}

	rrmanagerProxy := haproxy.GetRestProxy(register, "failover", failover.SupportedResources())

	failoverTask := quark.NewTask()
	failoverTask.User = "admin"
	fo := &failover.Failover{
		Id:                    "eef9cac140424cca804eefb6752bc57a",
		Zone:                  "zhanglei.test.",
		Zdnsuser:              "admin",
		MinConsensusNodeCount: 1,
		Name:              "www.zhanglei.test.",
		Type:              "aw",
		View:              "others",
		Active:            true,
		ProbeMethod:       "HTTP",
		Nodes:             []string{"beijing", "shanghai", "tianjin"},
		Frequency:         60,
		BearableThreshold: 5000,
		RetryCount:        0,
		SwitchDelay:       0,
		RetryInterval:     0,
		Port:              80,
		Uri:               "",
		Action:            "switch_to_backup",
		Backups:           []string{"10 192.168.1.1"},
		Host:              "www.zhanglei.com",
	}
	failoverTask.AddCmd(&rest.PutCmd{NewResource: fo})

	var ret string
	err = rrmanagerProxy.HandleTask(failoverTask, nil, &ret)
	if err != nil {
		fmt.Println(err.Error())
	}
	if ret != "" {
		fmt.Println("add failover failed: " + ret)
	}
	ret = ""
	<-time.After(2 * time.Second)
}
