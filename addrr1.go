package main

import (
	"time"

	"quark"
	"quark/haproxy"
	"quark/registry"
	"quark/rest"
	"zcloud-go/rrmanager"
)

func main() {
	register, err := registry.NewEtcdRegistry("192.168.142.132", []string{"http://192.168.249.12:2379"})
	if err != nil {
		panic("create registry failed:" + err.Error())
	}

	rrmanagerProxy := haproxy.GetRestProxy(register, "rrmanager", rrmanager.SupportedResources())

	rrTask := quark.NewTask()
	rrTask.User = "zhanglei"
	rr0 := &rrmanager.Rr{
		Id:    "0.zhanglei.test.id",
		Zone:  "zhanglei.test.",
		Name:  "www.zhanglei.test.",
		View:  "others",
		Type:  "a",
		Rdata: "192.168.1.0",
		Ttl:   3600,
	}
	rr1 := &rrmanager.Rr{
		Id:    "1.zhanglei.test.id",
		Zone:  "zhanglei.test.",
		Name:  "www.zhanglei.test.",
		View:  "c_ct_heilongjiang",
		Type:  "a",
		Rdata: "192.168.1.1",
		Ttl:   3600,
	}
	rr2 := &rrmanager.Rr{
		Id:    "2.zhanglei.test.id",
		Zone:  "zhanglei.test.",
		Name:  "www.zhanglei.test.",
		View:  "c_ct_jilin",
		Type:  "a",
		Rdata: "192.168.1.2",
		Ttl:   3600,
	}
	rr3 := &rrmanager.Rr{
		Id:    "3.zhanglei.test.id",
		Zone:  "zhanglei.test.",
		Name:  "www.zhanglei.test.",
		View:  "c_ct_liaoning",
		Type:  "a",
		Rdata: "192.168.1.3",
		Ttl:   3600,
	}
	rr4 := &rrmanager.Rr{
		Id:    "4.zhanglei.test.id",
		Zone:  "zhanglei.test.",
		Name:  "www.zhanglei.test.",
		View:  "c_ct_beijing",
		Type:  "a",
		Rdata: "192.168.1.4",
		Ttl:   3600,
	}
	rr5 := &rrmanager.Rr{
		Id:    "5.zhanglei.test.id",
		Zone:  "zhanglei.test.",
		Name:  "www.zhanglei.test.",
		View:  "c_ct_hebei",
		Type:  "a",
		Rdata: "192.168.1.5",
		Ttl:   3600,
	}
	rr6 := &rrmanager.Rr{
		Id:    "6.zhanglei.test.id",
		Zone:  "zhanglei.test.",
		Name:  "www.zhanglei.test.",
		View:  "c_ct_neimeng",
		Type:  "a",
		Rdata: "192.168.1.6",
		Ttl:   3600,
	}
	rr7 := &rrmanager.Rr{
		Id:    "7.zhanglei.test.id",
		Zone:  "zhanglei.test.",
		Name:  "www.zhanglei.test.",
		View:  "c_ct_shanxi",
		Type:  "a",
		Rdata: "192.168.1.7",
		Ttl:   3600,
	}
	rr8 := &rrmanager.Rr{
		Id:    "8.zhanglei.test.id",
		Zone:  "zhanglei.test.",
		Name:  "www.zhanglei.test.",
		View:  "c_ct_tianjign",
		Type:  "a",
		Rdata: "192.168.1.8",
		Ttl:   3600,
	}
	rr9 := &rrmanager.Rr{
		Id:    "9.zhanglei.test.id",
		Zone:  "zhanglei.test.",
		Name:  "www.zhanglei.test.",
		View:  "c_ct_anhui",
		Type:  "a",
		Rdata: "192.168.1.9",
		Ttl:   3600,
	}
	rr10 := &rrmanager.Rr{
		Id:    "10.zhanglei.test.id",
		Zone:  "zhanglei.test.",
		Name:  "www.zhanglei.test.",
		View:  "c_ct_fujian",
		Type:  "a",
		Rdata: "192.168.1.10",
		Ttl:   3600,
	}
	rrTask.AddCmd(&rest.PostCmd{NewResource: rr0})
	rrTask.AddCmd(&rest.PostCmd{NewResource: rr1})
	rrTask.AddCmd(&rest.PostCmd{NewResource: rr2})
	rrTask.AddCmd(&rest.PostCmd{NewResource: rr3})
	rrTask.AddCmd(&rest.PostCmd{NewResource: rr4})
	rrTask.AddCmd(&rest.PostCmd{NewResource: rr5})
	rrTask.AddCmd(&rest.PostCmd{NewResource: rr6})
	rrTask.AddCmd(&rest.PostCmd{NewResource: rr7})
	rrTask.AddCmd(&rest.PostCmd{NewResource: rr8})
	rrTask.AddCmd(&rest.PostCmd{NewResource: rr9})
	rrTask.AddCmd(&rest.PostCmd{NewResource: rr10})
	var ret string
	rrmanagerProxy.HandleTask(rrTask, nil, &ret)
	if ret != "" {
		panic(ret)
	}
	<-time.After(2 * time.Second)
}
