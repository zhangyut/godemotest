package main

import (
	"time"

	"flag"
	"quark"
	"quark/haproxy"
	"quark/registry"
	"quark/rest"
	"zcloud-go/rrmanager"
)

var (
	localip string
	etcd    string
	zone    string
	rrtype  string
	rrview  string
	rrdata  string
	userid  string
	rrname  string
	rrttl   int
)

func init() {
	flag.StringVar(&localip, "i", "127.0.0.1", "local ip")
	flag.StringVar(&etcd, "e", "http://127.0.0.1:2379", "etcd address")
	flag.StringVar(&zone, "zone", "", "zone")
	flag.StringVar(&rrtype, "type", "", "type")
	flag.StringVar(&rrview, "view", "", "view")
	flag.StringVar(&rrdata, "data", "", "data")
	flag.StringVar(&userid, "user", "", "user id")
	flag.StringVar(&rrname, "name", "", "name")
	flag.IntVar(&rrttl, "ttl", 0, "ttl")
}

func main() {
	flag.Parse()
	if zone == "" || userid == "" || rrview == "" || rrtype == "" || rrdata == "" || rrname == "" || rrttl < 60 {
		panic("param error.")
	}
	register, err := registry.NewEtcdRegistry(localip, []string{etcd})
	if err != nil {
		panic("create registry failed:" + err.Error())
	}

	rrmanagerProxy := haproxy.GetRestProxy(register, "rrmanager", rrmanager.SupportedResources())

	rrTask := quark.NewTask()
	rrTask.User = userid
	rr := &rrmanager.Rr{
		Zone:  zone,
		Name:  rrname,
		View:  rrview,
		Type:  rrtype,
		Rdata: rrdata,
		Ttl:   rrttl,
	}
	rrTask.AddCmd(&rest.PostCmd{NewResource: rr})
	var ret string
	err = rrmanagerProxy.HandleTask(rrTask, nil, &ret)
	if err != nil {
		panic(err.Error())
	}
	if ret != "" {
		panic(ret)
	}
	<-time.After(2 * time.Second)
}
