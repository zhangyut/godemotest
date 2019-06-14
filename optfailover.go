package main

import (
	"fmt"

	"quark"
	"quark/haproxy"
	"quark/registry"
	"quark/rest"
	"strings"
	"zcloud-go/failover"

	"flag"
)

var (
	user        string
	etcdcluster string
	local       string
	domainName  string
	viewId      string
	onlyLook    string
	opt         bool
)

func init() {
	flag.StringVar(&user, "u", "", "user id")
	flag.StringVar(&etcdcluster, "e", "", "etcd")
	flag.StringVar(&local, "l", "127.0.0.1", "local ip")
	flag.BoolVar(&opt, "opt", false, "active/disable failover")
}

func main() {
	flag.Parse()
	fmt.Println(opt)
	if etcdcluster == "" {
		fmt.Println("param error.")
		return
	}
	register, err := registry.NewEtcdRegistry(local, strings.Split(etcdcluster, ","))
	if err != nil {
		fmt.Println("create registry failed:" + err.Error())
		return
	}

	rrmanagerProxy := haproxy.GetRestProxy(register, "failover", failover.SupportedResources())

	failoverTask := quark.NewTask()
	failoverTask.User = user
	failoverTask.AddCmd(&rest.GetCmd{
		ResourceType: "failover",
		Conds:        map[string]interface{}{},
	})

	var ret string
	failovers := []failover.Failover{}
	err = rrmanagerProxy.HandleTask(failoverTask, &failovers, &ret)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if ret != "" {
		fmt.Println("get failover failed: " + ret)
		return
	}

	if len(failovers) == 0 {
		return
	}

	fmt.Println(failovers)

	upfailoverTask := quark.NewTask()
	upfailoverTask.User = user
	for i, _ := range failovers {
		failovers[i].Active = opt
		upfailoverTask.AddCmd(&rest.PutCmd{
			NewResource: &failovers[i],
		})
	}

	err = rrmanagerProxy.HandleTask(upfailoverTask, nil, &ret)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if ret != "" {
		fmt.Println("restart failover failed: " + ret)
		return
	}

	fmt.Println("restart succeed.")
}
