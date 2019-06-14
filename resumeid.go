package main

import (
	"fmt"
	"time"

	"quark"
	"quark/haproxy"
	"quark/registry"
	"zcloud-go/httpdns"
)

type Ret struct {
	Id int
}

func main() {
	//register, err := registry.NewEtcdRegistry("192.168.142.132", []string{"http://192.168.249.12:2379"})
	register, err := registry.NewEtcdRegistry("192.168.142.132", []string{"http://202.173.9.11:2379"})
	if err != nil {
		panic("create registry failed:" + err.Error())
	}

	httpdnsProxy := haproxy.GetHttpCmdProxy(register, "httpdns_cmd", httpdns.SupportedCmds())

	userTask := quark.NewTask()
	userTask.User = "admin"
	userTask.AddCmd(&httpdns.Resume{
		User: "zhanglei",
	})

	var ret string
	id := &Ret{}
	httpdnsProxy.HandleTask(userTask, id, &ret)
	if ret != "" {
		panic("create httpdns failed:" + ret)
	} else {
		fmt.Printf("succeed.")
		fmt.Println(*id)
	}

	<-time.After(2 * time.Second)
}
