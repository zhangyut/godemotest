package main

import (
	"fmt"

	"flag"
	"quark"
	"quark/haproxy"
	"quark/registry"
	"strings"
	"time"
	"zcloud-go/authorize"
	"zcloud-go/authorize/resource"
)

var (
	etcd      string
	local     string
	onlyCount int
	limit     int
	offset    int
	domain    string
	scope     string
	btime     string
	etime     string
)

func init() {
	flag.StringVar(&etcd, "e", "", "etcd address")
	flag.StringVar(&local, "i", "", "local address")
	flag.IntVar(&onlyCount, "o", 0, "only count")
	flag.IntVar(&limit, "l", 100, "limit")
	flag.IntVar(&offset, "s", 0, "offset")
	flag.StringVar(&domain, "d", "", "domain")
	flag.StringVar(&scope, "c", "", "scope")
	flag.StringVar(&btime, "bt", "", "begin time")
	flag.StringVar(&etime, "et", "", "end time")
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

	proxy := haproxy.GetHttpCmdProxy(register, "authorize_cmd", authorize.SupportedCmds())

	beginTime := time.Time{}
	endTime := time.Time{}
	if btime != "" {
		var err error
		beginTime, err = time.Parse("2006-01-02T15:04:05+08:00", btime)
		if err != nil {
			panic(err.Error())
		}
	}
	if etime != "" {
		var err error
		endTime, err = time.Parse("2006-01-02T15:04:05+08:00", etime)
		if err != nil {
			panic(err.Error())
		}
	}
	task := quark.NewTask()
	task.User = "admin"
	am := &authorize.GetAuthMessage{
		OnlyCount:       onlyCount != 0,
		Limit:           10,
		Offset:          offset,
		Zone:            domain,
		CreateTimeBegin: beginTime,
		CreateTimeEnd:   endTime,
	}
	task.AddCmd(am)

	var ret string
	if onlyCount == 0 {
		out := []resource.AuthMessage{}
		err = proxy.HandleTask(task, &out, &ret)
		if ret != "" {
			fmt.Println("add authorize failed: " + ret)
		}
		if err != nil {
			fmt.Println("add authorize error: " + err.Error())
		}
		fmt.Println(out)
		<-time.After(2 * time.Second)
	} else {
		var out int
		err = proxy.HandleTask(task, &out, &ret)
		if ret != "" {
			fmt.Println("add authorize failed: " + ret)
		}
		if err != nil {
			fmt.Println("add authorize error: " + err.Error())
		}
		fmt.Println(out)
		<-time.After(2 * time.Second)

	}
}
