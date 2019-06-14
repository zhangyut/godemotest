package main

import (
	"fmt"
	"strings"
	"time"

	"quark"
	"quark/registry"

	"dragondance/backend/bolt"
	"dragondance/client"
	"dragondance/dd"
	"flag"
	"sort"
)

var (
	failoverId string
	etcd       string
	localIp    string
	stype      string
)

func init() {
	flag.StringVar(&etcd, "e", "", "etcd address")
	flag.StringVar(&failoverId, "id", "", "failoverId")
	flag.StringVar(&localIp, "ip", "127.0.0.1", "local ip")
}

type SortedArray struct {
	Ids []string
}

func (p *SortedArray) Len() int {
	return len(p.Ids)
}
func (p *SortedArray) Less(i, j int) bool {
	return p.Ids[i] < p.Ids[j]
}
func (p *SortedArray) Swap(i, j int) {
	tmp := p.Ids[i]
	p.Ids[i] = p.Ids[j]
	p.Ids[j] = tmp
}

var boltvaluemap map[string][]byte

func main() {
	flag.Parse()

	boltvaluemap = make(map[string][]byte)
	diff := getIds()
	fmt.Println(diff)

	return

	if etcd == "" {
		fmt.Println("etcd can't nil.")
		return
	}

	if probeId == "" && diff == nil {
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

	if probeId == "" && diff != nil {
		for _, probeId = range diff {
			value, ok := boltvaluemap[probeId]
			if ok && value != nil {
				err = store.PutRaw(probeId, value)
				if err != nil {
					fmt.Println("put probeid:", probeId, " failed")
				}
			}
		}
		return
	}
}
