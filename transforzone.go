package main

import (
	"flag"
	"quark"
	"quark/haproxy"
	"quark/registry"
	"strings"
	"time"
	"zcloud-go/transferzone"
)

var (
	etcd    string
	localIp string
	user    string
	zone    string
	view    string
)

func init() {
	flag.StringVar(&etcd, "e", "http://127.0.0.1:2379", "etcd")
	flag.StringVar(&localIp, "i", "127.0.0.1", "local ip")
	flag.StringVar(&user, "u", "", "user")
	flag.StringVar(&view, "v", "", "view")
	flag.StringVar(&zone, "z", "", "zone")
}

func main() {
	flag.Parse()
	register, err := registry.NewEtcdRegistry(localIp, strings.Split(etcd, ","))
	if err != nil {
		panic("create registry failed:" + err.Error())
	}

	transferProxy := haproxy.GetHttpCmdProxy(register, "transferzone_cmd", transferzone.SupportedCmds())

	t := quark.NewTask()
	t.User = "admin"
	t.AddCmd(&transferzone.TransferZone{
		User: user,
		Zone: zone,
		View: view,
	})
	var ret string
	err = transferProxy.HandleTask(t, nil, &ret)
	if ret != "" {
		panic(ret)
	}
	if err != nil {
		panic(err.Error())
	}

	<-time.After(5 * time.Second)
}
