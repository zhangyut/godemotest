package main

import (
	"fmt"
	"time"

	"quark"
	"quark/httpcmd"
	"quark/registry"

	"flag"
	"zcloud-go/yundiapi"
	"zcloud-go/yundiapi/util"
)

var (
	token      string
	imageId    string
	imageValue string
)

func init() {
	flag.StringVar(&imageId, "imageid", "", "imageid")
	flag.StringVar(&imageValue, "imagevalue", "", "imagevalue")
}

func main() {
	flag.Parse()

	registry, _ := registry.NewEtcdRegistry("127.0.0.1", []string{"http://127.0.0.1:2379"})
	authproxy, err := httpcmd.GetProxy(registry, "yundiapi_cmd", yundiapi.SupportedCmds())
	if err != nil {
		fmt.Println(err.Error())
	}
	//set login notify
	task := quark.NewTask()
	task.AddCmd(&yundiapi.GetFlushNameResult{})
	errMsg := ""

	ret := util.ResultGetFlushName{}
	err = authproxy.HandleTask(task, &ret, &errMsg)
	if errMsg != "" {
		fmt.Println(errMsg)
		return
	}
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("success.")
	fmt.Println(ret)

	<-time.After(1 * time.Second)
}
