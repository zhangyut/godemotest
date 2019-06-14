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
	"zcloud-go/failover"
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

	err = store.Select("failover-rrset")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if probeId == "" && diff != nil {
		for _, probeId = range diff {
			fmt.Println("unknown task received,id %s", probeId)
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

	if stype == "probe" {
		rr := &failover.Rr{}
		if found, err := store.Get(probeId, &rr); found == false {
			fmt.Println("unknown task received,probe id %s", probeId)
			if err != nil {
				fmt.Println(err.Error())
			}
			return
		}

		fmt.Println("probe rr:%v", *rr)
	} else if stype == "failover" {
		rrset := &failover.RrSet{}
		if found, err := store.Get(probeId, &rrset); found == false {
			fmt.Println("unknown task received,failover id %s", probeId)
			if err != nil {
				fmt.Println(err.Error())
			}
			return
		}
		fmt.Println("failover rrset:%v", *rrset)
	}
}

func getIds() []string {
	diff := SortedArray{Ids: []string{}}
	db52 := SortedArray{Ids: []string{}}
	db53 := SortedArray{Ids: []string{}}
	db59 := SortedArray{Ids: []string{}}

	dbwrapper, err := bolt.NewBoltKVStore("bolt52.db")
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	db, err := dbwrapper.Select("failover-rrset")
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	myPrint := func(key string, value []byte) error {
		db52.Ids = append(db52.Ids, key)
		boltvaluemap[key] = value
		return nil
	}
	sort.Sort(&db52)
	err = db.ForEachKeyValue(myPrint)
	if err != nil {
		fmt.Println("read db52 err:", err.Error())
		return nil
	}

	///////////////
	dbwrapper, err = bolt.NewBoltKVStore("bolt53.db")
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	db, err = dbwrapper.Select("failover-rrset")
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	myPrint = func(key string, value []byte) error {
		db53.Ids = append(db53.Ids, key)
		boltvaluemap[key] = value
		return nil
	}
	sort.Sort(&db53)
	err = db.ForEachKeyValue(myPrint)
	if err != nil {
		fmt.Println("read db53 err:", err.Error())
		return nil
	}
	///////////////
	dbwrapper, err = bolt.NewBoltKVStore("bolt59.db")
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	db, err = dbwrapper.Select("failover-rrset")
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	myPrint = func(key string, value []byte) error {
		db59.Ids = append(db59.Ids, key)
		boltvaluemap[key] = value
		return nil
	}
	sort.Sort(&db59)
	err = db.ForEachKeyValue(myPrint)
	if err != nil {
		fmt.Println("read db59 err:", err.Error())
		return nil
	}

	for i, v := range db52.Ids {
		index := sort.SearchStrings(db53.Ids, v)
		if index == len(db53.Ids) || db53.Ids[index] != v {
			diff.Ids = append(diff.Ids, db52.Ids[i])
		}
	}

	for i, v := range db52.Ids {
		index := sort.SearchStrings(db59.Ids, v)
		if index == len(db59.Ids) || db59.Ids[index] != v {
			diff.Ids = append(diff.Ids, db52.Ids[i])
		}
	}

	for i, v := range db53.Ids {
		index := sort.SearchStrings(db52.Ids, v)
		if index == len(db52.Ids) || db52.Ids[index] != v {
			diff.Ids = append(diff.Ids, db53.Ids[i])
		}
	}

	for i, v := range db53.Ids {
		index := sort.SearchStrings(db59.Ids, v)
		if index == len(db59.Ids) || db59.Ids[index] != v {
			diff.Ids = append(diff.Ids, db53.Ids[i])
		}
	}

	for i, v := range db59.Ids {
		index := sort.SearchStrings(db52.Ids, v)
		if index == len(db52.Ids) || db52.Ids[index] != v {
			diff.Ids = append(diff.Ids, db59.Ids[i])
		}
	}

	for i, v := range db59.Ids {
		index := sort.SearchStrings(db53.Ids, v)
		if index == len(db53.Ids) || db53.Ids[index] != v {
			diff.Ids = append(diff.Ids, db59.Ids[i])
		}
	}

	sort.Sort(&diff)
	diffArray := []string{}
	for i, v := range diff.Ids {
		if i != 0 {
			if diffArray[len(diffArray)-1] == v {
				continue
			} else {
				diffArray = append(diffArray, diff.Ids[i])
			}
		} else {
			diffArray = append(diffArray, diff.Ids[i])
		}
	}

	return diffArray
}
