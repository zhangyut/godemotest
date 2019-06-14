package main

import (
	"fmt"

	"cement/uuid"
	"flag"
	"quark"
	"quark/haproxy"
	"quark/registry"
	"quark/rest"
	"strings"
	"zcloud-go/agentweb/resource"
)

var (
	etcd  string
	local string
)

func init() {
	flag.StringVar(&etcd, "e", "", "etcd address")
	flag.StringVar(&local, "i", "", "local address")
}

func main() {
	flag.Parse()
	if etcd == "" || local == "" {
		fmt.Println("parameter error.")
		return
	}
	register, err := registry.NewEtcdRegistry(local, strings.Split(etcd, ","))
	if err != nil {
		panic("create registry failed:" + err.Error())
	}

	proxy := haproxy.GetRestProxy(register, "agentweb", resource.SupportedResources())

	id, _ := uuid.Gen()
	task := quark.NewTask()
	task.User = "admin"
	agent := &resource.Agent{
		Id:       id,
		Name:     "agent1",
		Password: "agentp1",
		Describe: "代理商1",
	}
	task.AddCmd(&rest.PostCmd{NewResource: agent})

	var ret string
	err = proxy.HandleTask(task, nil, &ret)
	if ret != "" {
		fmt.Println("add agent failed: " + ret)
	}
	if err != nil {
		fmt.Println("add agent error: " + err.Error())
	}

	task.ClearCmd()
	agent.Name = "agent2"
	agent.Password = "agentp2"
	agent.Describe = "代理商2"
	task.AddCmd(&rest.PutCmd{NewResource: agent})
	ret = ""
	err = proxy.HandleTask(task, nil, &ret)
	if ret != "" {
		fmt.Println("update agent failed: " + ret)
	}
	if err != nil {
		fmt.Println("update agent error: " + err.Error())
	}

	task.ClearCmd()
	tid, _ := uuid.Gen()
	template := &resource.Template{
		Id:    tid,
		Agent: agent.Id,
	}
	task.AddCmd(&rest.PostCmd{NewResource: template})
	ret = ""
	err = proxy.HandleTask(task, nil, &ret)
	if ret != "" {
		fmt.Println("add template failed: " + ret)
	}
	if err != nil {
		fmt.Println("add template error: " + err.Error())
	}
}
