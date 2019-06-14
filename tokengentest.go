package main

import (
	"fmt"
	"time"

	"quark"
	"quark/httpcmd"
	"quark/registry"

	"flag"
	"zcloud-go/auth"
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

	registry, _ := registry.NewEtcdRegistry("202.173.9.22", []string{"http://202.173.9.22:2379"})
	authproxy, err := httpcmd.GetProxy(registry, "auth_cmd", auth.SupportedCmds())
	if err != nil {
		fmt.Println(err.Error())
	}
	//set login notify
	task := quark.NewTask()
	task.AddCmd(&auth.TokenGenCmd{
		GrantType:  "password",
		Username:   "admin",
		Password:   "admin123456",
		ExpireIn:   2,
		Address:    "127.0.0.1",
		ImageId:    imageId,
		ImageValue: imageValue,
	})
	errMsg := ""

	ret := auth.GenTokenResult{}
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
