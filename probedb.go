package main

import (
	center "crystalball/probecenter"
	"flag"
	"fmt"
)

var (
	probeId string
)

func init() {
	flag.StringVar(&probeId, "t", "", "probe id")
}

func main() {
	flag.Parse()
	store, err := center.NewProbeStore("./probenode_qingdao_probe.db")
	if err != nil {
		panic("create probe store failed:" + err.Error())
	}

	probes := []center.AddProbe{}
	store.ForEachProbe(func(probe *center.AddProbe) {
		probes = append(probes, *probe)
	})

	for _, v := range probes {
		if probeId != "" {
			if v.Params.GetId() == probeId {
				fmt.Println(v, v.Params.GetTimeout())
			}
		} else {
			fmt.Println(v, v.Params.GetTimeout())
		}

	}
}
