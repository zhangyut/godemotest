package main

import (
	"fmt"
	"strings"
	"time"

	"quark"
	"quark/registry"

	"dragondance/client"
	"dragondance/dd"
	"flag"
)

var (
	probeId string
	etcd    string
	localIp string
	stype   string
)

func init() {
	flag.StringVar(&etcd, "e", "", "etcd address")
	flag.StringVar(&probeId, "id", "", "probe id")
	flag.StringVar(&localIp, "ip", "127.0.0.1", "local ip")
	flag.StringVar(&stype, "type", "probe", "select type")
}

func main() {
	flag.Parse()

	if etcd == "" {
		fmt.Println("etcd can't nil.")
		return
	}

	if probeId == "" {
		fmt.Println("probe id can't nil.")
		return
	}

	var kvconf dd.DdConf

	registry, err := registry.NewEtcdRegistry(localIp, strings.Split(etcd, ","))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = registry.GetStruct(registry.NodeKey(quark.ThirdPartyKey, "kvstore"), &kvconf)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	store, err := proxy.NewKVClient(strings.Split(kvconf.Addrs, ","), 10*time.Second)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = store.Select("probe")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var probe interface{}
	if found, err := store.Get(probeId, &probe); found == false {
		fmt.Println("unknown task received,probe id %s", probeId)
		if err != nil {
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println("probe :%v", probe)
}
