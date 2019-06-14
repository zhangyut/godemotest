package main

import (
	"fmt"

	"cement/reflector"
	"crystalball/prober"
	"crystalball/prober/cdn"
	"crystalball/prober/cdnhttp"
	"crystalball/prober/dig"
	"crystalball/prober/http"
	"crystalball/prober/ping"
	"crystalball/prober/tcp"
	"dragondance/client"
	"dragondance/dd"
	"quark"
	"quark/haproxy"
	"quark/registry"
	"quark/rest"
	"strings"
	"time"
	"zcloud-go/failover"
	"zcloud-go/usermanager"

	"flag"
)

var (
	user        string
	etcdcluster string
	local       string
	knownTask   map[uint8]prober.Task
)

func init() {
	flag.StringVar(&user, "u", "", "user id")
	flag.StringVar(&etcdcluster, "e", "", "etcd")
	flag.StringVar(&local, "l", "127.0.0.1", "local ip")
	knownTask = map[uint8]prober.Task{
		uint8(48): &ping.PingTask{},
		uint8(49): &dig.DigTask{},
		uint8(50): &http.HttpTask{},
		uint8(51): &cdn.CDNTask{},
		uint8(52): &tcp.TcpConnTestTask{},
		uint8(53): &cdnhttp.CDNHttpTask{},
	}
}

type Policy struct {
	MinConsensusNodeCount int
	RetryCount            int
	Interval              int
	Nodes                 []string
}

type Probe struct {
	Policy  Policy
	Handler string
	Task    prober.Task
}

func getUsers(register *registry.EtcdRegistry) []usermanager.Zdnsuser {
	userProxy := haproxy.GetRestProxy(register, "usermanager", usermanager.SupportedResources())
	userTask := quark.NewTask()
	userTask.User = "admin"
	userTask.AddCmd(&rest.GetCmd{
		ResourceType: "zdnsuser",
		Conds:        map[string]interface{}{},
	})

	var ret string
	users := []usermanager.Zdnsuser{}
	err := userProxy.HandleTask(userTask, &users, &ret)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	if ret != "" {
		fmt.Println("get zdnsuser failed: " + ret)
		return nil
	}

	return users
}

func getFailovers(user string, register *registry.EtcdRegistry) []failover.Failover {
	failoverProxy := haproxy.GetRestProxy(register, "failover", failover.SupportedResources())

	failoverTask := quark.NewTask()
	failoverTask.User = user
	failoverTask.AddCmd(&rest.GetCmd{
		ResourceType: "failover",
		Conds:        map[string]interface{}{},
	})

	var ret string
	failovers := []failover.Failover{}
	err := failoverProxy.HandleTask(failoverTask, &failovers, &ret)
	if err != nil {
		fmt.Println("get failover failed1:" + err.Error())
		return nil
	}
	if ret != "" {
		fmt.Println("get failover failed2: " + ret)
		return nil
	}
	return failovers
}

func getProbeIds(kv *proxy.Client, fail failover.Failover) []string {
	//fmt.Println(failoverId)
	kv.Select("failover-rrset")
	ids := failover.RrSet{}
	kv.Get(fail.Id, &ids)
	for _, v := range ids.Ids {
		probeRr := failover.Rr{}
		found, err := kv.Get(v, &probeRr)
		if !found {
			fmt.Println(fmt.Sprintf("can't find prober rr :%s", v))
		}
		if err != nil {
			fmt.Println(err.Error())
		}
		/*
			fmt.Print("rr:")
			fmt.Print("\t")
			fmt.Println(probeRr.InnerRr)
		*/

	}
	if len(ids.Ids) <= 0 {
		fmt.Println(fmt.Sprintf("can't find prober id for failover:%v", fail))
	}
	return ids.Ids
}

func checkProbeId(kv *proxy.Client, probeId string) error {
	err := kv.Select("probe")
	task, ok := knownTask[probeId[0]]
	if ok == false {
		return nil
	}
	clone, _ := reflector.CloneStruct(task)
	p := Probe{
		Task: clone.(prober.Task),
	}
	found, err := kv.Get(probeId, &p)
	if !found {
		fmt.Println(fmt.Sprintf("can't find prober:%s", probeId))
		return nil
	}
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	/*
		fmt.Print("probe:")
		fmt.Print("\t")
		fmt.Println(p.Task)
	*/
	return nil
}

func restartFailover(user string, failovers []failover.Failover, register *registry.EtcdRegistry) {
	failoverProxy := haproxy.GetRestProxy(register, "failover", failover.SupportedResources())
	if len(failovers) <= 0 {
		fmt.Println("can't find failover")
		return
	}
	interval := time.Duration(60*1000/len(failovers)) * time.Millisecond
	for i, _ := range failovers {
		fmt.Print("restart failover:")
		fmt.Println(failovers[i])
		upfailoverTask := quark.NewTask()
		upfailoverTask.User = user
		if failovers[i].Active == true {
			upfailoverTask.AddCmd(&rest.PutCmd{
				NewResource: &failovers[i],
			})
		} else {
			fmt.Println("restart failover: skip")
			continue
		}
		var ret string
		err := failoverProxy.HandleTask(upfailoverTask, nil, &ret)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if ret != "" {
			fmt.Println("restart failover failed: " + ret)
			return
		}

		fmt.Println("restart failover succeed")
		<-time.After(interval)
	}
}

func main() {
	flag.Parse()
	if etcdcluster == "" {
		fmt.Println("param error.")
		return
	}
	register, err := registry.NewEtcdRegistry(local, strings.Split(etcdcluster, ","))
	if err != nil {
		fmt.Println("create registry failed:" + err.Error())
		return
	}
	var kvconf dd.DdConf
	err = register.GetStruct(register.NodeKey(quark.ThirdPartyKey, "kvstore"), &kvconf)
	if err != nil {
		panic("get kv conf failed:" + err.Error())
	}
	store, err := proxy.NewKVClient(strings.Split(kvconf.Addrs, ","), 10*time.Second)
	if err != nil {
		fmt.Println("create kv failed:" + err.Error())
		return
	}

	if user == "" {
		users := getUsers(register)
		for _, v := range users {
			failovers := getFailovers(v.Id, register)
			//restartFailover(v.Id, failovers, register)
			for _, v := range failovers {
				if !v.Active {
					continue
				}
				probeIds := getProbeIds(store, v)
				for i, _ := range probeIds {
					checkProbeId(store, probeIds[i])
				}
			}
		}
	} else {
		failovers := getFailovers(user, register)
		//restartFailover(user, failovers, register)
		for _, v := range failovers {
			if !v.Active {
				continue
			}
			probeIds := getProbeIds(store, v)
			for i, _ := range probeIds {
				checkProbeId(store, probeIds[i])
			}
		}
	}
}
