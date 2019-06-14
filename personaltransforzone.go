package main

import (
	"fmt"

	"flag"
	"quark"
	"quark/haproxy"
	"quark/registry"
	"strings"
	"time"
	"zcloud-go/personalservice/transferzone"
)

var (
	etcd    string
	localIp string
)

func init() {
	flag.StringVar(&etcd, "e", "http://127.0.0.1:2379", "etcd")
	flag.StringVar(&localIp, "i", "127.0.0.1", "local ip")
}

func main() {
	flag.Parse()
	register, err := registry.NewEtcdRegistry(localIp, strings.Split(etcd, ","))
	if err != nil {
		panic("create registry failed:" + err.Error())
	}

	transferProxy := haproxy.GetHttpCmdProxy(register, "personalservice_cmd", transferzone.SupportedCmds())

	t := quark.NewTask()
	t.User = "admin"
	t.AddCmd(&transferzone.TransferZone{
		User: "admin",
		Zone: "zhanglei.test.",
	})
	var ret string
	err = transferProxy.HandleTask(t, nil, &ret)
	if ret != "" {
		fmt.Println("get users failed: " + ret + " usesr: admin")
		return
	}
	if err != nil {
		fmt.Println("user : admin " + " transfer zone err:" + err.Error())
		return
	}

	<-time.After(2 * time.Second)
}
