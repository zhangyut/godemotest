package main

import (
	"flag"
	"fmt"
	"sqd"
	"strings"

	"quark"
	"quark/haproxy"
	"quark/registry"
)

var (
	etcdcluster string
	local       string
	topic       string
	addrs       string
	opt         string
	user        string
)

func init() {
	flag.StringVar(&etcdcluster, "e", "http://127.0.0.1:2379", "etcd cluster address.")
	flag.StringVar(&local, "i", "127.0.0.1", "local ip addr")
	flag.StringVar(&topic, "topic", "", "topic name")
	flag.StringVar(&addrs, "addrs", "0.0.0.0/0", "white list split with ','")
	flag.StringVar(&opt, "opt", "add", "add/delete/update opt")
	flag.StringVar(&user, "user", "", "user id")
}

func main() {
	flag.Parse()
	if topic == "" || user == "" {
		panic("topic or user is nil.")
	}
	whitelists := strings.Split(addrs, ",")

	registry, err := registry.NewEtcdRegistry(local, strings.Split(etcdcluster, ","))
	if err != nil {
		panic("create registry failed:" + err.Error())
	}

	//init nodeservice group and view
	proxy := haproxy.GetHttpCmdProxy(registry, "sqd_cmd", sqd.SupportedCmds())
	task := quark.NewTask()
	task.User = user
	if opt == "add" {
		task.AddCmd(&sqd.SetWhitelist{
			Topic: topic,
			Addrs: whitelists,
		})
	} else if opt == "delete" {
		task.AddCmd(&sqd.DeleteWhitelist{
			Topic: topic,
		})
	} else if opt == "update" {
		task.AddCmd(&sqd.UpdateWhitelist{
			Topic: topic,
			Addrs: whitelists,
		})
	} else {
		panic("not support")
	}
	ret := ""
	err = proxy.HandleTask(task, nil, &ret)
	if err != nil {
		fmt.Println(err.Error())
	}
	if ret != "" {
		fmt.Println(ret)
	}
}
