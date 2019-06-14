package main

import (
	"flag"
	"fmt"
	"quark"
	"quark/haproxy"
	"quark/registry"
	"quark/rest"
	"strings"
	"time"
	"zcloud-go/failover"
)

var (
	etcd  string
	local string
	user  string
)

func init() {
	flag.StringVar(&etcd, "e", "", "etcd")
	flag.StringVar(&local, "l", "", "local ip")
	flag.StringVar(&user, "u", "", "user")
}

func main() {
	flag.Parse()
	if local == "" || etcd == "" {
		panic("local or etcd is nil")
	}
	etcds := strings.Split(etcd, ",")
	if len(etcds) <= 0 {
		panic("etcd is nil")
	}

	register, err := registry.NewEtcdRegistry(local, etcds)
	if err != nil {
		panic("create registry failed:" + err.Error())
	}

	rrmanagerProxy := haproxy.GetRestProxy(register, "failover", failover.SupportedResources())

	failoverTask := quark.NewTask()
	failoverTask.User = user
	fo1 := &failover.Failover{
		Zone:                  "ctripgslb.com.",
		Zdnsuser:              user,
		MinConsensusNodeCount: 4,
		Name:                  "ctslb-pcisite.ctrip.com.jq.tel.ctripgslb.com.",
		Type:                  "a",
		View:                  "others",
		Active:                false,
		ProbeMethod:           "HTTPS",
		Nodes:                 []string{"beijing", "qingdao", "shenzhen", "hangzhou"},
		Frequency:             60,
		BearableThreshold:     5000,
		RetryCount:            0,
		SwitchDelay:           0,
		RetryInterval:         0,
		Port:                  0,
		Uri:                   "/do_not_delete/noc.gif",
		Action:                "switch_to_backup",
		Backups:               []string{"140.207.228.43"},
		Host:                  "ctslb-pcisite.ctrip.com",
	}
	failoverTask.AddCmd(&rest.PostCmd{NewResource: fo1})

	fo2 := &failover.Failover{
		Zone:                  "ctripgslb.com.",
		Zdnsuser:              user,
		MinConsensusNodeCount: 4,
		Name:                  "ctslb-pcisite.ctrip.com.oy.cmc.ctripgslb.com.",
		Type:                  "a",
		View:                  "others",
		Active:                false,
		ProbeMethod:           "HTTPS",
		Nodes:                 []string{"beijing", "qingdao", "shenzhen", "hangzhou"},
		Frequency:             60,
		BearableThreshold:     5000,
		RetryCount:            0,
		SwitchDelay:           0,
		RetryInterval:         0,
		Port:                  0,
		Uri:                   "/do_not_delete/noc.gif",
		Action:                "switch_to_backup",
		Backups:               []string{"140.206.211.82"},
		Host:                  "ctslb-pcisite.ctrip.com",
	}
	failoverTask.AddCmd(&rest.PostCmd{NewResource: fo2})

	fo3 := &failover.Failover{
		Zone:                  "ctripgslb.com.",
		Zdnsuser:              user,
		MinConsensusNodeCount: 4,
		Name:                  "ctslb-pcisite.ctrip.com.cnc.jq.ctripgslb.com.",
		Type:                  "a",
		View:                  "others",
		Active:                false,
		ProbeMethod:           "HTTPS",
		Nodes:                 []string{"beijing", "qingdao", "shenzhen", "hangzhou"},
		Frequency:             60,
		BearableThreshold:     5000,
		RetryCount:            0,
		SwitchDelay:           0,
		RetryInterval:         0,
		Port:                  0,
		Uri:                   "/do_not_delete/noc.gif",
		Action:                "switch_to_backup",
		Backups:               []string{"140.206.211.82", "101.226.248.39"},
		Host:                  "ctslb-pcisite.ctrip.com",
	}
	failoverTask.AddCmd(&rest.PostCmd{NewResource: fo3})

	fo4 := &failover.Failover{
		Zone:                  "ctripgslb.com.",
		Zdnsuser:              user,
		MinConsensusNodeCount: 4,
		Name:                  "ctslb-pcisite.ctrip.com.tel.oy.ctripgslb.com.",
		Type:                  "a",
		View:                  "others",
		Active:                false,
		ProbeMethod:           "HTTPS",
		Nodes:                 []string{"beijing", "qingdao", "shenzhen", "hangzhou"},
		Frequency:             60,
		BearableThreshold:     5000,
		RetryCount:            0,
		SwitchDelay:           0,
		RetryInterval:         0,
		Port:                  0,
		Uri:                   "/do_not_delete/noc.gif",
		Action:                "switch_to_backup",
		Backups:               []string{"101.226.248.39", "140.206.211.82"},
		Host:                  "ctslb-pcisite.ctrip.com",
	}
	failoverTask.AddCmd(&rest.PostCmd{NewResource: fo4})

	fo5 := &failover.Failover{
		Zone:                  "ctripgslb.com.",
		Zdnsuser:              user,
		MinConsensusNodeCount: 4,
		Name:                  "ctslb-pcisite.ctrip.com.cmc.jq.ctripgslb.com.",
		Type:                  "a",
		View:                  "others",
		Active:                false,
		ProbeMethod:           "HTTPS",
		Nodes:                 []string{"beijing", "qingdao", "shenzhen", "hangzhou"},
		Frequency:             60,
		BearableThreshold:     5000,
		RetryCount:            0,
		SwitchDelay:           0,
		RetryInterval:         0,
		Port:                  0,
		Uri:                   "/do_not_delete/noc.gif",
		Action:                "switch_to_backup",
		Backups:               []string{"221.130.198.232", "140.207.228.43"},
		Host:                  "ctslb-pcisite.ctrip.com",
	}
	failoverTask.AddCmd(&rest.PostCmd{NewResource: fo5})

	fo6 := &failover.Failover{
		Zone:                  "ctripgslb.com.",
		Zdnsuser:              user,
		MinConsensusNodeCount: 4,
		Name:                  "ctslb-pciws.ctrip.com.cmc.oy.ctripgslb.com.",
		Type:                  "a",
		View:                  "others",
		Active:                false,
		ProbeMethod:           "HTTPS",
		Nodes:                 []string{"beijing", "qingdao", "shenzhen", "hangzhou"},
		Frequency:             60,
		BearableThreshold:     5000,
		RetryCount:            0,
		SwitchDelay:           0,
		RetryInterval:         0,
		Port:                  0,
		Uri:                   "/do_not_delete/noc.gif",
		Action:                "switch_to_backup",
		Backups:               []string{"117.184.207.223", "140.206.211.83"},
		Host:                  "ctslb-pciws.ctrip.com",
	}
	failoverTask.AddCmd(&rest.PostCmd{NewResource: fo6})

	var ret string
	err = rrmanagerProxy.HandleTask(failoverTask, nil, &ret)
	if err != nil {
		fmt.Println(err.Error())
	}
	if ret != "" {
		fmt.Println("add failover failed: " + ret)
	}
	ret = ""
	<-time.After(5 * time.Second)
}
